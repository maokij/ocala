#!/usr/bin/env ruby
# frozen_string_literal: true

require "set"
require "strscan"
require File.expand_path("llpg", __dir__)

class GenerateArch
  module_eval LLPg.generate_code(<<-GRAMMAR)
    prog: ^{ v = [:prog] }! (sexp { v << val[0] }!)* { v }!;
    sexp: INT { val[0].value }!
        | STR { val[0].value }!
        | SYM { val[0].value }!
        | '(' { v = [] }! (sexp { v << val[0] }!)* ')' { v }!
        | '[' { v = [] }! (sexp { v << val[0] }!)* ']' { v }!
        | '`' sexp { [:quote, val[1]] }! ;
  GRAMMAR
  include Parser

  BOPS = %i[* / % + - << >> >>> < <= > >= == != & | ^ && ||].freeze
  INTERNS = {
    "#.REP": "KwREP",
    "#.INVALID": "KwINVALID",
    "#.jump": "KwJump",
    "#.call": "KwCall"
  }.freeze
  Operand = Struct.new(:name, :go, :oc, :asm, :alt, :temp)
  Token = Struct.new(:kind, :value, :pos)
  Arg = Struct.new(:v, :x)
  BMap = Struct.new(:x, :name, :mask, :min, :max, :map) # rubocop:disable Lint/StructNewOverride

  BCode = Struct.new(:t, :v) do # byte(n) low(x) high(x) imp(b x) temp
    def inspect
      "#{t.upcase}(#{v.map { |i| format('%02x', i) }.join(' ')})"
    end

    def pretty_print(pp)
      pp.text(inspect)
    end
  end

  OpArg = Struct.new(:t, :v) do # k(x) s(x t) v(x)
    def inspect
      t == :k ? v[0].to_s : "#{v[1]}(#{v[0]})"
    end

    def pretty_print(pp)
      pp.text(inspect)
    end
  end

  def token
    @token ||= scan
  end

  def consume_token
    @token.tap { @token = nil }
  end

  def parse_error(token, expected)
    n = @s.string[...token.pos].lines.size
    line = @s.string.lines[n]
    raise "#{token} / #{expected}\n#{n + 1}|#{line}"
  end

  def scan
    while @s.scan(/[ \t\n]+|;[^\n]*/)
      nil # NOP
    end

    pos = @s.pos
    if @s.eos?
      Token.new(:EOF, :nil, pos)
    elsif @s.scan(/[\[\](){}`]/)
      Token.new(@s[0], :nil, pos)
    elsif @s.scan(/-?(0x[_0-9a-f]+|0b[_01]+|0|[1-9][0-9]*)/i)
      Token.new(:INT, Integer(@s[0]), pos)
    elsif @s.scan(/"([^"]*)"/)
      Token.new(:STR, @s[1], pos)
    elsif @s.scan(%r{[_A-Za-z$?&|+\-*/.@=<>!^:#][0-9_A-Za-z$?&|+\-*/.@=<>!^:#]*})
      Token.new(:SYM, @s[0].intern, pos)
    else
      parse_error(Token.new(:ERR, nil, pos), :EOF)
    end
  end

  class Arch
    ATTRIBUTES = Set.new(%i[
                           opcodes operands operators bops uops bmaps
                           registers conditions examples inst_aliases token_aliases
                         ]).freeze
    attr_reader :name, :base, *ATTRIBUTES

    def initialize(name, base = nil)
      @name = name
      @base = base
      @opcodes = {}
      @operands = {}
      @operators = {}
      @bops = {}
      @uops = {}
      @bmaps = base ? base.bmaps : {}
      @examples = {}
      @inst_aliases = {}
      @token_aliases = {}

      @registers = []
      @conditions = []

      @operands[:_] = Operand.new(:_, "KwAny") unless base
    end

    def deep_merge(*rest)
      merge = lambda do |a, b|
        return b unless b.is_a?(Hash)
        return b unless a.is_a?(Hash)

        b.each { |k, v| a[k] = merge.call(a[k], v) }.then { a }
      end
      rest.reduce({}) { |r, i| merge.call(r, i) }
    end

    def all(attr, acc = [])
      raise unless ATTRIBUTES.include?(attr)

      (base ? base.all(attr, acc) : acc).push(send(attr))
    end

    def merged(attr)
      deep_merge(*all(attr))
    end

    def fetch(attr, *ks)
      raise unless ATTRIBUTES.include?(attr)

      send(attr)&.dig(*ks) || base&.fetch(attr, *ks)
    end
  end
  attr_reader :arch

  class Env
    attr_reader :outer, :items

    def initialize(outer = nil)
      @outer = outer
      @items = {}
    end

    def [](k)
      @items[k] || @outer&.[](k)
    end

    def []=(k, v)
      @items[k] = v
    end

    def enter
      Env.new(self)
    end
  end

  def read_arch(src)
    @form = {
      "=": lambda { |env, name, t = nil|
        OpArg.new(:s, [evaluate(env, name).x, t])
      },
      "=l": lambda { |env, name|
        BCode.new(:low, [evaluate(env, name).x])
      },
      "=h": lambda { |env, name|
        BCode.new(:high, [evaluate(env, name).x])
      },
      "=rl": lambda { |env, name, d, t = 1|
        BCode.new(:rlow, [evaluate(env, name).x, d, t])
      },
      "=rh": lambda { |env, name, d, t = 1|
        BCode.new(:rhigh, [evaluate(env, name).x, d, t])
      },
      "=m": lambda { |env, name, map|
        BCode.new(:map, [evaluate(env, name).x, arch.bmaps[map].x])
      },
      "=t": lambda { |_env|
        BCode.new(:temp, [0])
      },
      "=i": lambda { |env, name, base, mask, shift|
        BCode.new(:imp, [evaluate(env, name).x, evaluate(env, base), mask, shift])
      },
      "=U": lambda { |env, name|
        BCode.new(:unsupported, [evaluate(env, name).x])
      },
      "+": lambda { |env, *args|
        args.reduce(0) { |r, i| r | evaluate(env, i) }
      },
      quote: lambda { |_env, v|
        v
      },
      prog: lambda { |env, *rest|
        rest.each { |i| evaluate(env, i) }
      },
      arch: lambda { |env, name, *rest|
        name, base = Array(name)
        base = env[base] or raise "unknown base" if base

        env[name] = @arch = Arch.new(name, base)
        rest.each { |i| evaluate(env, i) }
        adjust_temp
        @archs << arch
      },
      map: lambda { |env, name, *kvs|
        env[name] = kvs.each_slice(2).to_h
        @form[name] = lambda { |ienv, k, shift = 0|
          k = evaluate(ienv, k).v
          v = evaluate(ienv, name)[k] or
            raise "!invalid key #{k} in #{name}"
          v << shift
        }
      },
      aliases: lambda { |_env, name, *names|
        arch.inst_aliases[name] = names
      },
      registers: lambda { |_env, *items|
        arch.registers.replace(items)
      },
      conditions: lambda { |_env, *items|
        items.each do |name, *aliases|
          arch.conditions << name
          aliases.each { arch.token_aliases[_1] = name }
        end
      },
      bytemap: lambda { |_env, name, mask, min, max, *items|
        x = arch.bmaps.size
        arch.bmaps[name] = BMap.new(x, name, mask, min, max, items)
      },
      operand: lambda { |env, name, go, oc, asm, alt = nil, temp = nil|
        env[name] = arch.operands[name] = Operand.new(name, "kw#{go}", oc, asm, alt, temp)
      },
      example: lambda { |env, name, *patterns|
        name = name.intern
        arch.examples[name] ||= {}

        if patterns[0] == [:*]
          oc, asm = patterns[1..].each_slice(2).to_a.transpose.map { |i| i.join(";") }
          arch.examples[name][[:*]] = [oc, asm]
        else
          patterns.each_slice(3) do |pat, oc, asm|
            puts "--> test: #{name} #{pat}" if ENV["DEBUG"]

            pat = pat.map { |i| env[i].is_a?(Hash) ? env[i].keys : [i] }
            pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

            pat.each do |operands|
              arch.examples[name][operands] = [oc, asm]
            end
          end
        end
      },
      opcode: lambda { |env, name, args, *patterns|
        name = name.intern
        arch.opcodes[name] ||= {}
        patterns.each_slice(2) do |pat, body|
          puts "--> #{name} #{pat}" if ENV["DEBUG"]
          raise "?size mismatch #{args} / #{pat}" unless args.size == pat.size

          pat = pat.map { |i| env[i].is_a?(Hash) ? env[i].keys : [i] }
          pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

          pat.each do |operands|
            tenv = env.enter
            args.each_with_index { |i, x| tenv[i] = Arg.new(operands[x], x) }

            found = arch.fetch(:opcodes, name, operands)
            bcodes = body.map do |i|
              i = evaluate(tenv, i)
              if i.is_a?(Integer)
                BCode.new(:byte, [i])
              else
                i.is_a?(BCode) ? i : (raise "invalid bcode")
              end
            end

            if found
              warn "already defined (#{name} #{operands.join(', ')}) as #{found}, skip #{bcodes}"
            else
              arch.opcodes[name][operands] = bcodes
            end
          end
        end
      },
      operator: lambda { |env, name, args, *patterns|
        name = name.intern
        if args.size == 1
          arch.uops[name] ||= 0
        elsif args.size == 2
          arch.bops[name] ||= 0
        end

        mapcode = lambda do |ienv, body|
          body.map do |e|
            [OpArg.new(:k, [e[0].intern])] + e[1..].map do |i|
              i = evaluate(ienv, i)
              case i
              when Integer then OpArg.new(:v, [i])
              when Operand then OpArg.new(:k, [i.name])
              when OpArg then i
              when Array then mapcode.call(ienv, i)
              else raise "invalid oparg #{i.inspect}"
              end
            end
          end
        end

        arch.operators[name] ||= {}
        patterns.each_slice(2) do |pat, body|
          puts "--> #{name} #{pat}" if ENV["DEBUG"]
          raise "?size mismatch #{args} / #{pat}" unless args.size == pat.size

          pat = pat.map { |i| env[i].is_a?(Hash) ? env[i].keys : [i] }
          pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

          pat.each do |operands|
            tenv = env.enter
            args.each_with_index { |i, x| tenv[i] = Arg.new(operands[x], x) }

            raise "already defined (#{name} #{operands.join(', ')})" if arch.operators[name][operands]

            arch.operators[name][operands] = mapcode.call(tenv, body)
          end
        end
      }
    }

    @archs = []
    @s = StringScanner.new(src)
    evaluate(Env.new, _parse([]))
  end

  def nest_map(m)
    r = {}
    m.each do |k, v|
      q = r
      k.each do |i|
        q = (q[i] ||= {})
      end
      q[nil] = v
    end
    r
  end

  def nest_map2(m)
    r = {}
    m.each do |k, v|
      a, b = k
      r[a] ||= {}
      r[a][b] = v
    end
    r
  end

  def adjust_temp
    alts = arch.merged(:operands).filter { |_, v| v.alt }
    arch.opcodes.each_value do |pats|
      as = {}
      pats.each do |pat, body|
        temp = false
        apat = pat.map do |i|
          j = alts[i]
          if j
            temp |= j.temp
            j.alt
          else
            i
          end
        end
        if apat != pat && !pats[apat]
          as[apat] = temp ? [BCode.new(:temp, [0])] * body.size : body
        end
      end
      pats.merge!(as)
    end
  end

  def evaluate(env, e)
    case e
    in Integer | String
      e
    in Symbol
      env[e] or raise "?unknown symbol #{e}"
    in [Symbol => name, *args]
      raise "?unknown func #{name}" unless @form[name]

      @form[name].call(env, *args)
    else
      raise "?invalid form #{e}"
    end
  end

  def intern(name)
    return INTERNS[name] if INTERNS[name]

    name = name.to_s
    @interns[name] ||= "kw#{name.gsub(/\W/, '')}"
  end

  def generate_table_code(_package, arch)
    operand_code = lambda do |a|
      arch.fetch(:operands, a)&.go || "nil"
    end

    bcodes_code = lambda do |ls|
      s = ls.map do |i|
        v = i.v.map.with_index do |j, x|
          format("A%d: 0x%02x", x, [j].pack("C")[0].ord)
        end
        "{Kind: Bc#{i.t.capitalize}, #{v.join(', ')}}"
      end
      "{\n#{s.join(",\n")},\n}"
    end

    scodes_code = lambda do |es, nest = false|
      prefix = nest ? "&Vec" : ""
      s = es.map do |ls|
        t = ls.map do |i|
          next scodes_code.call(i, true) if i.is_a?(Array)

          case i.t
          when :k
            k = i.v[0]
            a = arch.fetch(:operands, k)
            if a
              "&Operand{Kind: #{a.go}}"
            elsif !k.start_with?("#") && !arch.fetch(:opcodes, k)
              raise "unknown opcode #{k}"
            else
              intern(k)
            end
          when :v
            "Int(#{i.v[0]})"
          when :s
            "&Vec{Int(#{i.v[0]}), #{operand_code.call(i.v[1])}}"
          end
        end
        "#{prefix}{#{t.join(', ')}}"
      end
      "#{prefix}{\n#{s.join(",\n")},\n}"
    end

    map_code_list = lambda do |pats|
      s = []
      pats.each do |k, v|
        if v.is_a? Array
          s << %(#{operand_code.call(k)}: InstDat#{bcodes_code.call(v)},)
        else
          s << %(#{operand_code.call(k)}: InstPat{)
          s.concat(map_code_list.call(v))
          s << %(},)
        end
      end
      s
    end

    code = []
    suffix = arch.base ? arch.name.to_s.gsub(/\W/, "").capitalize : ""

    arch.operands.each_value do |v|
      code << %!var #{v.go} = Intern("#{v.oc}")! unless v.go.start_with?("Kw")
    end
    code << ""

    code << "var asmOperands#{suffix} = map[*Keyword]AsmOperand{"
    arch.operands.each_value do |v|
      next unless v.asm

      s = v.asm.gsub(/%[BW]/, "%")
      t = s == v.asm ? "false" : "true"
      code << %(  #{v.go}: { "#{s}", #{t} },)
    end
    code << "}" << ""

    code << "var tokenWords#{suffix} = [][]string{"
    [arch.registers, arch.conditions].each do |v|
      s = v.map do |i|
        a = arch.fetch(:operands, i)&.oc or raise "invalid operand #{i}"
        %("#{a}")
      end
      code << "    { #{s.join(', ')} },"
    end
    uops = arch.uops.keys.map { %("#{_1}") }
    base_bops = (arch.base&.merged(:bops) || {}).keys | BOPS
    bops = (arch.bops.keys - base_bops).map { %("#{_1}") }
    code << "    { #{uops.join(', ')} },"
    code << "    { #{bops.join(', ')} },"
    code << "}" << ""

    unless arch.token_aliases.empty?
      code << "var tokenAliases#{suffix} = map[string]string{"
      arch.token_aliases.each do |k, v|
        code << %(  "#{k}": "#{v}",)
      end
      code << "}" << ""
    end

    unless arch.inst_aliases.empty?
      code << "var instAliases#{suffix} = map[string][]string{"
      arch.inst_aliases.each do |k, v|
        v = v.map { %("#{_1}") }.join(", ")
        code << %(  "#{k}": { #{v} },)
      end
      code << "}" << ""
    end

    unless arch.opcodes.empty?
      code << ""
      code << "var instMap#{suffix} = InstPat{"
      arch.opcodes.each do |name, pats|
        code << %(  #{intern(name)}: InstPat{)
        code.concat(map_code_list.call(nest_map(pats)))
        code << %(  },)
      end
      code << "}" << ""
    end

    return if arch.operators.empty?

    code << "var ctxOpMap#{suffix} = CtxOpMap{"
    arch.operators.each do |name, pats|
      code << %!        Intern("#{name}"): {!
      nest_map2(pats).each do |a, bs|
        code << %(            #{operand_code.call(a)}: {)
        bs.map do |b, template|
          code << %(                #{operand_code.call(b)}: #{scodes_code.call(template)},)
        end
        code << %(            },)
      end
      code << %(        },)
    end
    code << "}" << ""
  end

  def generate_arch(path)
    package = File.basename(File.dirname(path))
    @interns = {}
    s = @archs.map do |arch|
      generate_table_code(package, arch).join("\n")
    end

    code = ["package #{package}", ""]
    code << %(import . "ocala/internal/core" //lint:ignore ST1001 core) << ""

    code << "var bmaps = [][]byte{"
    @archs[0].bmaps.sort_by { |_, v| v.x }.each do |k, v|
      map = v.map.join(", ")
      code << %(  {#{v.mask}, #{v.min}, #{v.max}, #{map}}, // #{v.x}: #{k})
    end
    code << "}" << ""

    @interns.sort_by { |_, v| v }.each do |k, v|
      code << %!var #{v} = Intern("#{k}")!
    end

    s = gofmt(s.unshift(code.join("\n")).join("\n"))
    File.write(path.sub(/\.lisp\z/, ".g.go"), s)
  end

  def testdata_path(path, name)
    dir = File.expand_path("../testdata", path)
    File.join(dir, name)
  end

  def pfill(s)
    s.sub("%B", "5").sub("%W", "1234")
  end

  def use_example(name, pat, a, b)
    c = @all_examples.dig(name, pat) or return
    c[0] == :_ and return

    s, t = c.map { _1.split(/(?<! );/).map(&:strip).join("\n    ") }
    a << "    #{s}" unless s.empty?
    b << "    #{t}" unless t.empty?
    true
  end

  def find_asm_example(name, pat)
    c = @all_examples.dig(name, pat) or return
    c[1] if c[0] == :_
  end

  def generate_opcodes_oc(path, arch)
    nlabels = 0
    cocl = []
    casm = []

    add_label = lambda do
      (nlabels += 1).tap do |n|
        cocl << "L#{n}:"
        casm << "L#{n}:"
      end
    end

    suffix = arch.base ? arch.name.to_s.gsub(/\A\W/, "_") : ""
    use_example(:$prologue, [:*], cocl, casm)
    arch.merged(:opcodes).each do |name, pats|
      n = add_label.call
      next if use_example(name, [:*], cocl, casm)

      pats.each do |operands, body|
        next if body[0].t == :temp
        next if use_example(name, operands, cocl, casm)

        alt = find_asm_example(name, operands)
        operands = operands.map { |i| arch.fetch(:operands, i) or raise "unknown operand #{i.inspect}" }
        ocs = operands.map { |i| pfill(i.oc) }
        asms = operands.map { |i| pfill(i.asm) }
        rel = body.find { |i| i.t == :rlow }
        if rel
          x = rel.v[0]
          asms[x] = ocs[x] = "L#{n - 1}"
          cocl << "    #{name} #{ocs.join(' ')}"
          casm << "    #{name} #{asms.join(', ')}"

          asms[x] = ocs[x] = "L#{n + 2}"
        end

        cocl << "    #{name} #{ocs.join(' ')}"
        casm <<
          if alt
            "    #{alt}"
          else
            "    #{name} #{asms.join(', ')}"
          end
      end
    end
    use_example(:$epilogue, [:*], cocl, casm)

    File.write(testdata_path(path, "opcodes#{suffix}.oc"), cocl.join("\n") << "\n")
    File.write(testdata_path(path, "opcodes#{suffix}.asm"), casm.join("\n") << "\n")
  end

  def generate_operators_oc(path, arch)
    suffix = arch.base ? arch.name.to_s.gsub(/\A\W/, "_") : ""

    alts = arch.merged(:operands).filter { |_, v| v.alt }
    arch.merged(:operators).each_value do |pats|
      as = pats.filter_map do |pat, body|
        apat = pat.map { |i| alts[i]&.alt || i }
        apat != pat && !pats[apat] && [apat, body]
      end
      pats.merge!(as.to_h)
    end

    allopr = arch.merged(:operands).keys - [:_]
    cocl = []
    casm = []

    expand = lambda do |pat, template|
      op = template[0].v[0]
      next if op.start_with?("#")

      inst = arch.fetch(:opcodes, op) or raise "unknown opcode #{template}"
      oprs = []
      asms = []
      template[1..].map do |i|
        raise "!#{i}" if i.is_a?(Array)

        case i.t
        when :k
          oprs << i.v[0]
          asms << pfill(arch.fetch(:operands, i.v[0]).asm)
        when :v
          oprs << ((-128..255).cover?(i.v[0]) ? :N : :NN)
          asms << "##{i.v[0]}"
        when :s
          oprs << (i.v[1] || pat[i.v[0]])
          asms << pfill(arch.fetch(:operands, oprs.last).asm)
        end
      end
      m = inst[oprs]
      m and m[0].t != :temp and [op.to_s, *asms]
    end

    find_all = lambda do |pat, body|
      pat = pat.map { |i| i == :_ ? allopr : [i] }
      pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

      pat.map do |args|
        expanded = body.map { |i| expand.call(args, i) }
        [args, expanded]
      end
    end

    use_example(:$prologue, [:*], cocl, casm)
    arch.merged(:operators).each do |name, pats|
      next if use_example(name, [:*], cocl, casm)

      m = {}
      pats.each do |operands, body|
        next if use_example(name, operands, cocl, casm)

        find_all.call(operands, body).each do |args, expanded|
          next if m[args]

          m[args] = true
          next if use_example(name, args, cocl, casm)

          unless expanded.all?
            # warn "?? #{name} #{operands} --> #{args} #{body.inspect} #{expanded.inspect}"
            next
          end

          alt = find_asm_example(name, args)
          args = args.map { |i| pfill(arch.fetch(:operands, i).oc) }
          args[0] = "$(#{args[0]})" if args[0].match?(/\A\d/)
          cocl << "    #{args[0]} #{name} #{args[1]}".rstrip
          if alt
            casm << "    #{alt}"
          else
            expanded.each do |i|
              casm << "    #{i[0]} #{i[1..].join(', ')}"
            end
          end
        end
      end
    end
    use_example(:$operators, [:*], cocl, casm)
    use_example(:$epilogue, [:*], cocl, casm)

    File.write(testdata_path(path, "operators#{suffix}.oc"), cocl.join("\n") << "\n")
    File.write(testdata_path(path, "operators#{suffix}.asm"), casm.join("\n") << "\n")
  end

  def generate_testdata(path)
    @archs.each do |arch|
      @all_examples = arch.merged(:examples)
      generate_opcodes_oc(path, arch)
      generate_operators_oc(path, arch)
    end
  end

  def gofmt(s)
    IO.popen(["gofmt", "-s"], "r+") do |io|
      io.puts(s)
      io.close_write
      io.read
    end
  end

  def run(args)
    raise if args.empty?

    command = args.shift
    case command
    when "arch"
      read_arch(File.read(args[0]))
      generate_arch(args[0])
    when "testdata"
      read_arch(File.read(args[0]))
      generate_testdata(args[0])
    else
      raise "unknown subcommand: #{command}"
    end
  end
end

GenerateArch.new.run(ARGV.dup)

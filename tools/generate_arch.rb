#!/usr/bin/env ruby
# frozen_string_literal: true

require "strscan"
require File.expand_path("llpg", __dir__)

class GenerateArch
  module_eval LLPg.generate_code(<<-GRAMMAR)
    sexp: INT { val[0].value }!
        | STR { val[0].value }!
        | SYM { val[0].value }!
        | '(' { v = [] }! (sexp { v << val[0] }!)* ')' { v }!
        | '[' { v = [] }! (sexp { v << val[0] }!)* ']' { v }!
        | '`' sexp { [:quote, val[1]] }! ;
  GRAMMAR
  include Parser

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
        BCode.new(:map, [evaluate(env, name).x, @bmaps[map].x])
      },
      "=t": lambda { |_env|
        BCode.new(:temp, [0])
      },
      "=i": lambda { |env, name, base, mask, shift|
        BCode.new(:imp, [evaluate(env, name).x, evaluate(env, base), mask, shift])
      },
      "+": lambda { |env, *args|
        args.reduce(0) { |r, i| r | evaluate(env, i) }
      },
      quote: lambda { |_env, v|
        v
      },
      arch: lambda { |env, _name, *rest|
        rest.each { |i| evaluate(env, i) }
      },
      map: lambda { |env, name, *kvs|
        env[name] = kvs.each_slice(2).to_h
        @form[name] = lambda { |ienv, k, shift = 0|
          k = evaluate(ienv, k).v
          v = evaluate(ienv, name)[k] or
            raise "!invalid key #{k} in #{map}"
          v << shift
        }
      },
      aliases: lambda { |_env, name, *names|
        @inst_aliases[name] = names
      },
      registers: lambda { |_env, *items|
        @scanmap[:REG] = items
      },
      conditions: lambda { |_env, *items|
        items.each_slice(2) do |as, bs|
          a = as[0]
          as = as[1..]
          b = bs[0]
          bs = bs[1..]
          @scanmap[:COND] << a << b
          as.each { @token_aliases[_1] = a }
          bs.each { @token_aliases[_1] = b }
          @condmap[a] = b
          @condmap[b] = a
        end
      },
      bytemap: lambda { |_env, name, mask, min, max, *items|
        x = @bmaps.size
        @bmaps[name] = BMap.new(x, name, mask, min, max, items)
      },
      operand: lambda { |env, name, go, oc, asm, alt = nil, temp = nil|
        env[name] = @operands[name] = Operand.new(name, "kw#{go}", oc, asm, alt, temp)
      },
      example: lambda { |env, name, *patterns|
        name = name.intern
        @examples[name] ||= {}

        if patterns[0] == [:*]
          oc, asm = patterns[1..].each_slice(2).to_a.transpose.map { |i| i.join(";") }
          @examples[name][[:*]] = [oc, asm]
        else
          patterns.each_slice(3) do |pat, oc, asm|
            puts "--> test: #{name} #{pat}" if ENV["DEBUG"]

            pat = pat.map { |i| env[i].is_a?(Hash) ? env[i].keys : [i] }
            pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

            pat.each do |operands|
              @examples[name][operands] = [oc, asm]
            end
          end
        end
      },
      opcode: lambda { |env, name, args, *patterns|
        name = name.intern
        @opcodes[name] ||= {}
        patterns.each_slice(2) do |pat, body|
          puts "--> #{name} #{pat}" if ENV["DEBUG"]
          raise "?size mismatch #{args} / #{pat}" unless args.size == pat.size

          pat = pat.map { |i| env[i].is_a?(Hash) ? env[i].keys : [i] }
          pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

          pat.each do |operands|
            tenv = env.enter
            args.each_with_index { |i, x| tenv[i] = Arg.new(operands[x], x) }

            found = @opcodes[name][operands]
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
              @opcodes[name][operands] = bcodes
            end
          end
        end
      },
      operator: lambda { |env, name, args, *patterns|
        name = name.intern
        if args.size == 1
          @uops[name] ||= 0
        elsif args.size == 2
          @bops[name] ||= 0
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

        @operators[name] ||= {}
        patterns.each_slice(2) do |pat, body|
          puts "--> #{name} #{pat}" if ENV["DEBUG"]
          raise "?size mismatch #{args} / #{pat}" unless args.size == pat.size

          pat = pat.map { |i| env[i].is_a?(Hash) ? env[i].keys : [i] }
          pat = pat.empty? ? [[]] : pat[0].product(*pat[1..])

          pat.each do |operands|
            tenv = env.enter
            args.each_with_index { |i, x| tenv[i] = Arg.new(operands[x], x) }

            raise "already defined (#{name} #{operands.join(', ')})" if @operators[name][operands]

            @operators[name][operands] = mapcode.call(tenv, body)
          end
        end
      }
    }
    @s = StringScanner.new(src)
    @opcodes = {}
    @operands = {}
    @operators = {}
    @bops = {}
    @uops = {}
    @bmaps = {}
    @scanmap = { REG: [], COND: [] }
    @condmap = {}
    @examples = {}
    @inst_aliases = {}
    @token_aliases = {}
    e = _parse([])
    env = Env.new
    @operands[:_] = env[:_] = Operand.new(:_, "KwAny")
    evaluate(Env.new, e)
    adjust_temp
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
    alts = @operands.filter_map { |k, v| v.alt && [k, v] }.to_h
    @opcodes.each_value do |pats|
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

  def generate_table_code(package)
    operand_code = lambda do |a|
      @operands[a]&.go || "nil"
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
            a = @operands[k]
            if a
              "&Operand{Kind: #{a.go}}"
            elsif !k.start_with?("#") && !@opcodes[k]
              raise "unknown opcode #{k}"
            else
              %!Intern("#{k}")!
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
    code << "package #{package}" << ""
    code << %(import . "ocala/internal/core" //lint:ignore ST1001 core) << ""

    @operands.each_value do |v|
      code << %!var #{v.go} = Intern("#{v.oc}")! unless v.go.start_with?("Kw")
    end
    code << ""

    code << "var operandToAsmMap = map[*Keyword](struct{ s string; t bool }){"
    @operands.each_value do |v|
      next unless v.asm

      s = v.asm.gsub(/%[BW]/, "%")
      t = s == v.asm ? "false" : "true"
      code << %(  #{v.go}: { s: "#{s}", t: #{t} },)
    end
    code << "}" << ""

    code << "var tokenWords = [][]string{"
    @scanmap.slice(:REG, :COND).each do |k, v|
      s = v.map do |i|
        a = @operands[i]&.oc or raise "invalid operand #{i} for #{k}"
        %("#{a}")
      end
      code << "    { #{s.join(', ')} },"
    end
    uops = @uops.keys.map { %("#{_1}") }
    bops = %i[* / % + - << >> >>> < <= > >= == != & | ^ && ||]
    bops = (@bops.keys - bops).map { %("#{_1}") }
    code << "    { #{uops.join(', ')} },"
    code << "    { #{bops.join(', ')} },"
    code << "}" << ""

    code << "var oppositeConds = map[*Keyword]*Keyword{"
    @condmap.each do |k, v|
      code << %(    #{@operands[k].go}: #{@operands[v].go},)
    end
    code << "}" << ""

    code << "var bmaps = [][]byte{"
    @bmaps.sort_by { |_k, v| v.x }.each do |k, v|
      map = v.map.join(", ")
      code << %(  {#{v.mask}, #{v.min}, #{v.max}, #{map}}, // #{v.x}: #{k})
    end
    code << "}" << ""

    code << "var tokenAliases = map[string]string{"
    @token_aliases.each do |k, v|
      code << %(  "#{k}": "#{v}",)
    end
    code << "}" << ""

    code << "var instAliases = map[string][]string{"
    @inst_aliases.each do |k, v|
      v = v.map { %("#{_1}") }.join(", ")
      code << %(  "#{k}": { #{v} },)
    end
    code << "}" << ""

    code << ""
    code << "var instMap = InstPat{"
    @opcodes.each do |name, pats|
      code << %!  Intern("#{name}"): InstPat{!
      code.concat(map_code_list.call(nest_map(pats)))
      code << %(  },)
    end
    code << "}" << ""

    code << "var ctxOpMap = map[*Keyword]map[*Keyword]map[*Keyword][][]Value{"
    @operators.each do |name, pats|
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
    s = generate_table_code(package).join("\n")
    s = gofmt(s)
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
    c = @examples.dig(name, pat) or return
    s, t = c.map { _1.split(/(?<! );/).map(&:strip).join("\n    ") }
    a << "    #{s}" unless s.empty?
    b << "    #{t}" unless t.empty?
    true
  end

  def generate_opcodes_oc(path)
    nlabels = 0
    cocl = []
    casm = []

    add_label = lambda do
      (nlabels += 1).tap do |n|
        cocl << "L#{n}:"
        casm << "L#{n}:"
      end
    end

    use_example(:$prologue, [:*], cocl, casm)
    @opcodes.each do |name, pats|
      n = add_label.call
      next if use_example(name, [:*], cocl, casm)

      pats.each do |operands, body|
        next if body[0].t == :temp
        next if use_example(name, operands, cocl, casm)

        operands = operands.map { |i| @operands[i] }
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
        casm << "    #{name} #{asms.join(', ')}"
      end
    end
    use_example(:$epilogue, [:*], cocl, casm)

    File.write(testdata_path(path, "opcodes.oc"), cocl.join("\n") << "\n")
    File.write(testdata_path(path, "opcodes.asm"), casm.join("\n") << "\n")
  end

  def generate_operators_oc(path)
    alts = @operands.filter_map { |k, v| v.alt && [k, v] }.to_h
    @operators.each_value do |pats|
      as = pats.filter_map do |pat, body|
        apat = pat.map { |i| alts[i]&.alt || i }
        apat != pat && !pats[apat] && [apat, body]
      end
      pats.merge!(as.to_h)
    end

    allopr = @operands.keys - [:_]
    cocl = []
    casm = []

    expand = lambda do |pat, template|
      op = template[0].v[0]
      next if op.start_with?("#")

      inst = @opcodes[op] or raise "unknown opcode #{template}"
      oprs = []
      asms = []
      template[1..].map do |i|
        raise "!#{i}" if i.is_a?(Array)

        case i.t
        when :k
          oprs << i.v[0]
          asms << pfill(@operands[i.v[0]].asm)
        when :v
          oprs << ((-128..255).cover?(i.v[0]) ? :N : :NN)
          asms << "##{i.v[0]}"
        when :s
          oprs << (i.v[1] || pat[i.v[0]])
          asms << pfill(@operands[oprs.last].asm)
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
    @operators.each do |name, pats|
      next if use_example(name, [:*], cocl, casm)

      m = {}
      pats.each do |operands, body|
        next if use_example(name, operands, cocl, casm)

        find_all.call(operands, body).each do |args, expanded|
          next if m[args]

          m[args] = true

          unless expanded.all?
            # warn "?? #{name} #{operands} --> #{args} #{body.inspect} #{expanded.inspect}"
            next
          end

          args = args.map { |i| pfill(@operands[i].oc) }
          args[0] = "$(#{args[0]})" if args[0].match?(/\A\d/)
          cocl << "    #{args[0]} #{name} #{args[1]}".rstrip

          expanded.each do |i|
            casm << "    #{i[0]} #{i[1..].join(', ')}"
          end
        end
      end
    end
    use_example(:$operators, [:*], cocl, casm)
    use_example(:$epilogue, [:*], cocl, casm)

    File.write(testdata_path(path, "operators.oc"), cocl.join("\n") << "\n")
    File.write(testdata_path(path, "operators.asm"), casm.join("\n") << "\n")
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
      generate_opcodes_oc(args[0])
      generate_operators_oc(args[0])
    else
      raise "unknown subcommand: #{command}"
    end
  end
end

GenerateArch.new.run(ARGV.dup)

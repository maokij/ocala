#!/usr/bin/env ruby
# frozen_string_literal: true

require "pp"
require "yaml"
require "strscan"
require "forwardable"
require "optparse"

module LLPg
  VERSION = "0.0.1"

  class << self
    attr_accessor :debug_mode
  end

  GRAMMAR_RULES = <<~GRAMMAR
    # grammar: rule*;
    grammar: opt_rule_list;
    opt_rule_list: rule opt_rule_list;
    opt_rule_list:;

    # rule: NONTERM ':' rhs_list ';';
    rule: lhs ':' opt_init rhs_list ';' { add_rules(val[0], val[2], val[3], nil) }!;
    lhs: NONTERM                        { @lhs = val[0] }!;

    # rhs_list: rhs ('|' rhs)*;
    rhs_list: rhs opt_rhs_list_tail              { [val[0], *val[1]] }!;
    opt_rhs_list_tail: '|' rhs opt_rhs_list_tail { [val[1], *val[2]] }!;
    opt_rhs_list_tail:                           { [] }!;

    # rhs: ACTION? | elem elem*
    rhs: opt_action                 { [[nil, val[0]]] }!;
    rhs: elem opt_rhs_tail          { [val[0], *val[1]] }!;
    opt_rhs_tail: elem opt_rhs_tail { [val[0], *val[1]] }!;
    opt_rhs_tail:                   { [] }!;

    # elem: prim ('?' | '*')? ACTION?
    elem: prim opt_qualifier opt_action { [adjust_token(@lhs, val[0], *val[1]), val[2]] }!;
    opt_qualifier: '?'                  { [:opt, val[0].value] }!;
    opt_qualifier: '*'                  { :rep }!;
    opt_qualifier:                      { false }!;

    opt_action: ACTION                  { val[0].value }!;
    opt_action:                         { nil }!;

    opt_init: '^' ACTION                { val[1].value}!;
    opt_init:                           { nil}!;

    # prim: NONTERM | TERM | '(' rhs_list ')';
    prim: TERM                      { val[0] }!;
    prim: NONTERM                   { val[0] }!;
    prim: '(' opt_init rhs_list ')' { add_rules(@lhs, val[1], val[2], :group) }!;
  GRAMMAR

  SYMBOL_TOKENS = %i[TERM NONTERM].freeze

  Token = Struct.new(:kind, :value, :var, :path, :pos, :bolpos, :lineno, :colno, :s) do
    def format_line
      line = s.string.lines[lineno].rstrip
      "#{path}:#{lineno + 1}:#{colno}\n| #{line}\n| #{' ' * colno}^"
    end
  end

  Rule = Struct.new(:lhs, :rhs, :vars, :actions, :inline_option) do
    def to_s
      "#{lhs} -> #{rhs.join(' ')}"
    end

    def inspect
      "<Rule: #{self}>"
    end
  end

  module Utils
    def debug(s, obj = nil)
      return unless LLPg.debug_mode

      warn "--- [DEBUG] [#{caller(1..1).first}] #{s}"
      yield if block_given?
      warn obj.pretty_inspect if obj
    end

    def indent(s, n = 1)
      s = s.join("\n") unless s.is_a?(String)
      m = self.class.const_defined?(:INDENT_SIZE) ? self.class::INDENT_SIZE : 2
      s.lines.map { |i| (" " * n * m) << i }.join
    end
  end

  extend Utils

  class ParseError < RuntimeError
    attr_accessor :token

    def initialize(message, token)
      super(message)
      @token = token
    end
  end

  class GrammarReader
    def initialize(reader_options)
      @re_bof = reader_options[:re_bof] || /#/
      @re_eof = reader_options[:re_eof] || /%%/
      @additional_action = reader_options[:additional_action]
      setup_internal_parser
    end

    def make_token(kind, value, var, pos)
      Token.new(kind, value, var, @path, pos, @bolpos, @lineno, pos - @bolpos, @s)
    end

    def parse_error(token, expected)
      kind = token.kind.is_a?(Symbol) ? token.kind.to_s : "'#{token.kind}'"
      raise ParseError.new("Invalid token #{kind}, expected #{expected}", token)
    end

    def seek_to_bof
      return unless @s.skip_until(@re_bof)

      pos = @s.pos
      @s.pos = 0
      while @s.skip_until(/\n/) && @s.pos < pos
        @lineno += 1
        @bolpos = @s.pos
      end
      @s.pos = pos
    end

    def scan
      loop do
        if @s.bol? && @s.scan(/[ \t]*%options +([^\n]+)/)
          begin
            @options.update(YAML.safe_load("{ #{@s[1].strip} }"))
          rescue Psych::SyntaxError
            raise ParseError.new("Invalid %options", make_token(:OPTIONS, "", nil, @s.pos))
          end
        elsif @s.scan(/[ \t]+|#[^\n]*/)
          # SKIP
        elsif @s.scan(/\n/)
          @lineno += 1
          @bolpos = @s.pos
        else
          break
        end
      end

      pos = @s.pos
      if @s.eos? || @s.scan(@re_eof)
        make_token(:EOF, "", nil, pos)
      elsif @s.scan(/([A-Z][A-Z_0-9]*)(-([a-z_][a-z_0-9]+))?|('[^']+')/)
        if @s[1]
          make_token(:TERM, @s[1], @s[3], pos)
        else
          make_token(:TERM, @s[4], nil, pos)
        end
      elsif @s.scan(/([a-z_][a-z_0-9]*)(-([a-z_][a-z_0-9]+))?/)
        make_token(:NONTERM, @s[1], @s[3], pos)
      elsif @s.scan(/\{(.+?)\}!/m)
        make_token(:ACTION, @s[1], nil, pos)
      elsif @s.scan(/\?(:([\w.]+))?/m)
        make_token("?", @s[2], nil, pos)
      elsif @s.scan(/[:;|*?()^]/)
        make_token(@s[0], @s[0], nil, pos)
      else
        raise ParseError.new("invalid charactor", make_token(nil, nil, nil, pos))
      end
    end

    def token
      @token ||= scan # .tap { |it| pp it}
    end

    def consume_token
      @token.tap { @token = nil }
    end

    def initialize_state(src)
      @table = Table.new
      @token = nil
      @options = {}
      @unique_id = 0
      @s = StringScanner.new(src)
      @s.pos = @lineno = @bolpos = 0
    end

    def setup_internal_parser
      debug_mode = LLPg.debug_mode
      LLPg.debug_mode = nil
      initialize_state(GRAMMAR_RULES)
      @path = "<internal>"

      while token.kind != :EOF
        token.kind == :NONTERM or parse_error(token, "NONTERM")
        lhs = consume_token

        token.kind == ":" or parse_error(token, ":")
        consume_token

        rhs = []
        rhs << [consume_token] while SYMBOL_TOKENS.include?(token.kind)
        rhs << [nil, consume_token.value] if token.kind == :ACTION

        add_rules(lhs, nil, [rhs], false)
        token.kind == ";" or parse_error(token, ";")
        consume_token
      end

      code = RubyCodeGenerator.new.generate(@table.build, @options)
      instance_eval(code << " ; extend Parser")
      LLPg.debug_mode = debug_mode
    end

    def make_inline_lhs(token, inline_option)
      @unique_id += 1
      suffix = "#{@unique_id}#{inline_option[0]}"
      token.dup.tap { |it| it.value = "_#{it.value}_#{suffix}" }
    end

    def add_rules(lhs, init, rhs_list, inline_option)
      lhs = make_inline_lhs(lhs, inline_option) if inline_option

      @table.inits[lhs.value] = init
      rhs_list.each do |rhs|
        syms = []
        vars = []
        actions = []

        rhs.each do |token, action|
          if token
            syms << token.value
            vars << token.var
          else
            actions.pop
          end
          actions << action
        end
        @table.rules << Rule.new(lhs.value, syms, vars, actions, inline_option)
      end

      lhs
    end

    def adjust_token(lhs, token, inline_option, default = nil)
      return token unless inline_option

      lhs = make_inline_lhs(lhs, inline_option)
      if inline_option == :opt
        @table.rules << Rule.new(lhs.value, [token.value], [nil], [nil], :opt)
        @table.rules << Rule.new(lhs.value, [], [], default ? ["RET #{default}"] : [], :skip)
      elsif inline_option == :rep
        @table.rules << Rule.new(lhs.value, [token.value, lhs.value], [nil, nil], [nil, nil], :rep)
        @table.rules << Rule.new(lhs.value, [], [], [], :skip)
      end

      lhs
    end

    def read(path)
      @path = path
      parse(path == "-" ? $stdin.read : File.read(path))
    end

    def parse(src)
      initialize_state(src)
      seek_to_bof
      _parse([])
      @additional_action&.call(@options, @s)
      [@table.build, @options, @s.rest]
    end
  end

  class Table
    include Enumerable
    extend Forwardable
    def_delegators :@table, :[], :keys, :values, :each, :each_key, :each_value

    EMPTY = :@EMPTY
    attr_reader :rules, :start, :nonterms, :inits

    def initialize
      @rules = []
      @inits = {}
    end

    def term?(sym)
      !nonterm?(sym)
    end

    def nonterm?(sym)
      !!@nonterms[sym]
    end

    def inline_option(sym)
      @nonterms.dig(sym, 0)&.inline_option
    end

    def build(start = nil)
      @start = start || @rules.find { |i| !i.inline_option }.lhs

      @nonterms = {}
      @first = {}
      @follow = {}

      LLPg.debug "rules" do
        @rules.each { |i| warn i }
      end

      @rules.each do |i|
        @nonterms[i.lhs] ||= []
        @nonterms[i.lhs] << i
        @first[i.lhs] ||= []
        @first[i.lhs] << EMPTY if i.rhs.empty?
        @follow[i.lhs] ||= []
      end

      make_first
      make_follow
      make_table
      LLPg.debug "summary" do
        @nonterms.each_key { |i| warn make_summary(i) unless inline_option(i) }
      end

      self
    end

    def first(syms, items = [], empty_case = [EMPTY])
      syms = syms.dup
      loop do
        sym = syms.shift
        if sym.nil? # empty
          items.replace(items | empty_case)
          break
        elsif term? sym
          items.replace(items | [sym])
          break
        else
          items.replace(items | (@first[sym] - [EMPTY]))
          break unless @first[sym].include?(EMPTY)
        end
      end

      items
    end

    def make_first
      loop do
        changed = false

        @rules.each do |rule|
          set = @first[rule.lhs]
          size = set.size
          changed |= first(rule.rhs, set).size > size
        end

        break unless changed
      end

      LLPg.debug "first set" do
        @first.each { |k, v| warn "#{k} : #{v.join(', ')}" }
      end
    end

    def make_follow
      @follow[@start] << "EOF"

      loop do
        changed = false

        @rules.each do |rule|
          rule.rhs.each_with_index do |sym, x|
            next unless nonterm?(sym)

            set = @follow[sym]
            size = set.size
            set.replace(set | first(rule.rhs[(x + 1)..], [], @follow[rule.lhs]))
            changed |= set.size > size
          end
        end

        break unless changed
      end

      LLPg.debug "follow set" do
        @follow.each { |k, v| warn "#{k} : #{v.join(', ')}" }
      end
    end

    def make_table
      @table = {}
      @rules.each do |rule|
        set = (@table[rule.lhs] ||= {})
        first(rule.rhs, [], @follow[rule.lhs]).each do |i|
          raise "NOT LL(1) input: #{i}, conflict: #{set[i]} / #{rule}" if set[i]

          set[i] = rule
        end
      end

      LLPg.debug "table" do
        @table.each do |sym, set|
          warn "#{sym}:"
          set.group_by(&:last).each do |rule, pairs|
            warn "  #{pairs.map(&:first).join(', ')} =>\n    #{rule}"
          end
        end
      end
    end

    def make_summary(sym, memo = {})
      return sym if memo[sym]

      found = @nonterms.dig(sym, 0)
      return sym unless found
      return sym if !found.inline_option && !memo.empty?

      memo[sym] = true

      case found.inline_option
      when :opt
        "#{make_summary(found.rhs[0], memo)}?"
      when :rep
        "#{make_summary(found.rhs[0], memo)}*"
      when :group, nil
        s = @nonterms[sym].map do |rule|
          rule.rhs.map do |i|
            make_summary(i, memo)
          end.join(" ")
        end.join("\n| ")

        if found.inline_option
          "(\n  #{LLPg.indent(s)}\n  )"
        else
          "#{found.lhs} ->\n  #{LLPg.indent(s)}"
        end
      else
        raise
      end
    end
  end

  class RubyCodeGenerator
    include Utils

    INDENT_SIZE = 2
    INPUT_SUFFIX_RE = /\.rb\z/.freeze
    OUTPUT_SUFFIX = ".g.rb"
    READER_OPTIONS = {
      re_bof: /^=begin[ \t]+#[^\n]*/,
      re_eof: /=end[ \t]+#/
    }.freeze

    def generate(table, options)
      @table = table
      @options = options

      code = []
      code << "# GENERATED BY llpg"
      code << ""

      klass = @options["class"]
      if klass
        prefix = "#{klass}::"
        code << "class #{klass} ; end"
        code << ""
      end

      code << "module #{prefix}Parser"
      code << "  def _inline_rule(); yield([]); end"
      code << "  def _parse(st)"
      code << "    result = #{method_name(@table.start)}(st)"
      code << "    token.kind == :EOF or parse_error(token, 'EOF')"
      code << "    result"
      code << "  end"

      @table.each_key do |sym|
        next if @table.inline_option(sym)

        code << ""
        code << "  def #{method_name(sym)}(st)"
        code << "    val = []"
        code << indent(generate_rules_code(sym), 2)
        code << "  end"
      end

      code << "end\n"
      code.join("\n")
    end

    def generate_rules_code(sym, wrapper: nil, rep: nil)
      set = @table[sym]
      skip = @table[wrapper]&.filter_map { |k, v| v.rhs.empty? && k } || []
      expected = set.keys | skip

      annotations = []
      code = []
      code << adjust_action(@table.inits[sym]) if @table.inits[sym]
      code << "case token.kind"
      set.group_by(&:last).each do |rule, pairs|
        annotations << "# #{rule}"
        pat = pairs.map { |i| sym_in_code(i[0]) }.join(", ")
        code << "when #{pat}"
        if rule.rhs.empty?
          action = rule.actions[0]
          code << indent(adjust_action(action)) if action
          code << (rep ? "  break" : "  # NOP")
        else
          rule.rhs.size.times do |x|
            code << indent(generate_matcher_code(rule, x, rep))
          end
          code << "  redo" if wrapper && rep
        end
      end

      if wrapper
        annotations.unshift "# #{wrapper} -> #{sym}#{rep ? '*' : '?'}"
        unless skip.empty?
          pat = skip.map { |i| sym_in_code(i) }.join(", ")
          code << "when #{pat}"
          code << "  break"
        end
      end

      code << "else"
      code << "  parse_error(token, #{expected.join(', ').inspect})"
      code << "end"
      annotations.unshift("## Rules").concat(code)
    end

    def generate_matcher_code(rule, x, rep)
      code = []
      sym = rule.rhs[x]
      var = rule.vars[x] || "_"
      action = rule.actions[x]
      assign = "#{var} = val[#{x}] ="

      if rep && sym == rule.lhs
        code << "redo"
      elsif @table.inline_option(sym)
        code << assign
        code << indent(generate_inline_code(sym), 1)
      elsif @table.nonterm?(sym)
        code << "#{assign} #{method_name(sym)}(st)"
      else
        code << "token.kind == #{sym_in_code(sym)} or parse_error(token, #{sym.inspect})" if x > 0
        code << "#{assign} consume_token"
      end
      code << adjust_action(action) if action
      code
    end

    def generate_inline_code(sym)
      rule = @table.nonterms[sym][0]
      inner = rule.rhs[0]

      body =
        case rule.inline_option
        when :group
          generate_rules_code(sym)
        when :opt
          if @table.inline_option(inner) == :group
            generate_rules_code(inner, wrapper: sym)
          else
            generate_rules_code(sym)
          end
        when :rep
          if @table.inline_option(inner) == :group
            generate_rules_code(inner, wrapper: sym, rep: true)
          else
            generate_rules_code(sym, rep: true)
          end
        else
          raise "invalid inline option"
        end

      code = []
      code << "_inline_rule do |val|"
      code << indent(body)
      code << "end"
    end

    def adjust_action(action)
      action.strip
    end

    def method_name(sym)
      "_parse_#{sym}"
    end

    def sym_in_code(sym)
      sym[0] == "'" ? sym : ":#{sym}"
    end
  end

  class GoCodeGenerator
    include Utils

    INDENT_SIZE = 4
    INPUT_SUFFIX_RE = /(\.llpg)?\.go\z/.freeze
    OUTPUT_SUFFIX = ".g.go"
    READER_OPTIONS = {
      re_bof: %r{^[ \t]*/\*\*:rules[^\n]*},
      re_eof: %r{\*\*/},
      additional_action: ->(options, s) { reader_action(options, s) }
    }.freeze

    ASCII_SYMBOL_MAP = {
      "!" => "EX",
      "#" => "HA",
      "$" => "DL",
      "%" => "PE",
      "&" => "AM",
      "(" => "LP",
      ")" => "RP",
      "*" => "AS",
      "+" => "PL",
      "," => "CM",
      "-" => "MI",
      "." => "DT",
      "/" => "SL",
      ":" => "CL",
      ";" => "SC",
      "<" => "LT",
      "=" => "EQ",
      ">" => "GT",
      "?" => "QU",
      "@" => "AT",
      "[" => "LS",
      "]" => "RS",
      "^" => "CA",
      "_" => "UN",
      "`" => "GA",
      "{" => "LC",
      "|" => "VE",
      "}" => "RC",
      "~" => "TI"
    }.freeze

    attr_accessor :parser_type, :value_type, :token_type, :error_type

    def self.reader_action(options, s)
      options["reserved_words"] = []
      _, words = *s.rest.match(%r{/\*\*:reservedWords:(.+?)\*\*/}m)
      return unless words

      words.strip.lines.each do |line|
        _, a, b, = *line.strip.match(/(\w+):(.+)/)
        next unless a

        options["reserved_words"] << [a, b.strip.split]
      end
    end

    def generate(table, options)
      @table = table
      @options = {
        "parser_type" => "*Parser",
        "value_type" => "Value",
        "token_type" => "*Token"
      }.merge(options)

      code = []
      code << "// GENERATED BY llpg"
      code << "//lint:file-ignore U1000 generated code"
      code << "//lint:file-ignore SA4004 generated code"
      code << ""

      package = @options["package"] or raise
      code << "package #{package}"
      code << ""

      prologue = @options["prologue"]
      if prologue
        code << prologue
        code << ""
      end

      @parser_type = @options["parser_type"] or raise
      @value_type = @options["value_type"] or raise
      @token_type = @options["token_type"] or raise

      @table.rules.each do |rule|
        next unless rule.inline_option == :opt

        rule.actions[0] = "RET _1 // DEFAULT OPT ACTION"
      end

      terms = @table.rules.map(&:rhs).flatten.reject { |i| @table.nonterm?(i) }.uniq
      code << "const ("
      code << "    tkEOF = iota"
      code << "    tkSCANERROR"
      terms.each do |i|
        code << "    #{sym_in_code(i)}"
      end
      code << ")"
      code << ""

      code << "var tokenKinds = map[string]int32{"
      terms.each do |i|
        next unless i[0] == "'"

        code << "    \"#{i[1...-1]}\": #{sym_in_code(i)},"
      end
      code << "}"
      code << ""

      code << "var tokenLabels = [...]string{"
      code << "    \"EOF\", // tkEOF"
      code << "    \"SCANERROR\", // tkSCANERROR"
      terms.each do |i|
        code << "    \"#{i}\", // #{sym_in_code(i)}"
      end
      code << "}"
      code << ""

      reserved_words = @options["reserved_words"]
      unless reserved_words.empty?
        code << "var reservedWords = map[string]int32{"
        reserved_words.each do |token, names|
          names.each do |i|
            code << "  \"#{i}\": tk#{token},"
          end
        end
        code << "}"
        code << ""
      end

      code << "func (p #{parser_type}) _parse() (res #{value_type}, ok bool) {"
      code << "    res = p.#{method_name(@table.start)}()"
      code << "    if (p.PeekToken().Kind != tkEOF) {"
      code << "         p.ErrorUnexpected(p.PeekToken(), \"EOF\")"
      code << "    }"
      code << "    ok = true"
      code << "    return"
      code << "}"

      @table.each_key do |sym|
        next if @table.inline_option(sym)

        code << ""
        code << "/** ** **"
        code << @table.make_summary(sym).lines.map { |i| "  #{i}" }.join
        code << " ** ** **/"
        code << "func (p #{parser_type}) #{method_name(sym)}() (res #{value_type}) {"
        code << indent(generate_rules_code(sym))
        code << "    return"
        code << "}"
      end

      code.join("\n") << "\n"
    end

    def generate_rules_code(sym, inline: false, wrapper: nil, rep: nil)
      set = @table[sym]
      annotations = []

      code = []
      code << adjust_action(@table.inits[sym], inline) if @table.inits[sym]
      code << "switch p.PeekToken().Kind {"
      set.group_by(&:last).each do |rule, pairs|
        annotations << "// RULE: #{rule}"
        pat = pairs.map { |i| sym_in_code(i[0]) }.join(", ")
        code << "case #{pat}:"
        if rule.rhs.empty?
          action = rule.actions[0]
          action = action ? adjust_action(action, inline) : "// NOP"
          code << indent(action, 2)
        else
          rule.rhs.size.times do |x|
            code << indent(generate_matcher_code(rule, x, rep, inline), 2)
          end
          code << "" << indent("continue", 2) if wrapper && rep
        end
      end

      skip = []
      if wrapper
        annotations.unshift "// RULE: #{wrapper} -> #{sym}#{rep ? '*' : '?'}"
        skip = @table[wrapper].filter_map { |k, v| v.rhs.empty? && k } || []
        unless skip.empty?
          pat = skip.map { |i| sym_in_code(i) }.join(", ")
          code << "case #{pat}:"
          code << "    // NOP"
        end
      end

      expected = set.keys | skip
      code << "default:"
      code << "    p.ErrorUnexpected(p.PeekToken(), #{syms_as_string(expected)})"
      code << "}"

      annotations + code
    end

    def generate_matcher_code(rule, x, rep, inline)
      sym = rule.rhs[x]
      var = rule.vars[x] || "_#{x + 1}"
      action = rule.actions[x]
      sym_inline = @table.inline_option(sym)

      rhs_comment = rule.rhs.dup.tap { |it| it.insert(x, "@") }.join(" ")
      code = ["// #{rhs_comment}"]
      if rep && sym == rule.lhs
        code << "continue"
      elsif sym_inline
        opt_term = sym_inline == :opt && @table.term?(@table.nonterms[sym][0].rhs[0])
        var_type = opt_term ? token_type : value_type
        code << "var #{var} #{var_type}; _ = #{var}"
        code << generate_inline_code(sym, x)
      elsif @table.nonterm?(sym)
        code << "#{var} := p.#{method_name(sym)}(); _ = #{var}"
      else
        if x > 0
          code << "if p.PeekToken().Kind != #{sym_in_code(sym)} {"
          code << "    p.ErrorUnexpected(p.PeekToken(), #{syms_as_string([sym])})"
          code << "}"
        end
        code << "#{var} := p.ConsumeToken(); _ = #{var}"
      end
      code << adjust_action(action, inline) if action
      code << "" if x < rule.rhs.size - 1
      code
    end

    def generate_inline_code(sym, x)
      rule = @table.nonterms[sym][0]
      inner = rule.rhs[0]

      body =
        case rule.inline_option
        when :group
          generate_rules_code(sym, inline: true)
        when :opt
          if @table.inline_option(inner) == :group
            generate_rules_code(inner, inline: true, wrapper: sym)
          else
            generate_rules_code(sym, inline: true)
          end
        when :rep
          if @table.inline_option(inner) == :group
            generate_rules_code(inner, inline: true, wrapper: sym, rep: true)
          else
            generate_rules_code(sym, inline: true, rep: true)
          end
        else
          raise "invalid inline option"
        end

      code = []
      code << "for {"
      code << "    res := &_#{x + 1}; _ = res"
      code << indent(body)
      code << "    break"
      code << "}"
    end

    def adjust_action(action, inline)
      var = inline ? "*res" : "res"
      action.strip.gsub(/(^| )RET +/, "\\1#{var} = ")
    end

    def method_name(sym)
      parts = sym.to_s.split("_")
      "_parse#{parts.map(&:capitalize).join}"
    end

    def sym_in_code(sym)
      sym = sym[1..-2].chars.map { |i| ASCII_SYMBOL_MAP[i] }.join if sym[0] == "'"
      "tk#{sym}"
    end

    def syms_as_string(syms)
      syms.join(", ").inspect
    end
  end

  GENERATORS = [RubyCodeGenerator, GoCodeGenerator].freeze

  class CLI
    def make_output_path(path, options)
      base = path.sub(@generator::INPUT_SUFFIX_RE, "")
      suffix = options["output_suffix"] || @generator::OUTPUT_SUFFIX
      raise "invalid output suffix" if suffix.match?(%r{[/\\]|\.\.})

      @output_path = base + suffix
    end

    def find_generator(path)
      @generator = GENERATORS.find { |i| i::INPUT_SUFFIX_RE.match?(path) } or
        raise "unknown generator type"
    end

    def run(args)
      opt = OptionParser.new
      opt.banner += " path"
      opt.version = LLPg::VERSION
      opt.on("-o arg", "--output arg") { |arg| @output_path = arg }
      opt.on("-t arg", "--type arg") { |arg| find_generator ".#{arg}" }
      opt.on("--debug") { LLPg.debug_mode = true }
      opt.parse!(args)

      raise "input file required" unless args.size == 1

      path = args.shift
      find_generator path unless @generator

      reader = GrammarReader.new(@generator::READER_OPTIONS)
      table, options = reader.read(path)
      code = @generator.new.generate(table, options)

      make_output_path path, options unless @output_path
      File.write(@output_path, code)
    rescue ParseError => e
      warn "parse error: #{e.message} at #{e.token.format_line}"
      LLPg.debug_mode ? raise : exit(1)
    rescue RuntimeError => e
      warn "runtime error: #{e.message}"
      LLPg.debug_mode ? raise : exit(1)
    end
  end

  def self.generate_code(src)
    reader = GrammarReader.new(RubyCodeGenerator::READER_OPTIONS)
    table, options = reader.parse(src)
    RubyCodeGenerator.new.generate(table, options)
  end
end

LLPg::CLI.new.run(ARGV.dup) if $PROGRAM_NAME == __FILE__

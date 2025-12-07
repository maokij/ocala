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

    opt_init: '^' ACTION                { val[1].value }!;
    opt_init:                           { nil }!;

    # prim: NONTERM | TERM | '(' rhs_list ')';
    prim: TERM                      { val[0] }!;
    prim: NONTERM                   { val[0] }!;
    prim: '(' opt_init rhs_list ')' { add_rules(@lhs, val[1], val[2], :group) }!;
  GRAMMAR

  SYMBOL_TOKENS = %i[TERM NONTERM].freeze

  Token = Struct.new(:kind, :value, :var, :pos, :lineno, :owner) do
    def format_line
      colno = pos - owner.line_indexes[lineno]
      line = owner.lines[lineno]
      "#{owner.path}:#{lineno + 1}:#{colno}\n| #{line}| #{' ' * colno}^"
    end
  end

  Var = Struct.new(:name, :desc)

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
      ws = " " * n * m
      s.lines.map { |i| /\S/.match?(i) ? "#{ws}#{i}" : i }.join
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
    attr_reader :path, :lines, :line_indexes

    def initialize(reader_options)
      @re_bof = reader_options[:re_bof] || /#/
      @re_eof = reader_options[:re_eof] || /%%/
      @additional_action = reader_options[:additional_action]
      setup_internal_parser
    end

    def make_token(kind, value, var, pos)
      Token.new(kind, value, var, pos, @lineno, self) # .tap { warn _1.format_line }
    end

    def parse_error(token, expected)
      kind = token.kind.is_a?(Symbol) ? token.kind.to_s : "'#{token.kind}'"
      raise ParseError.new("Invalid token #{kind}, expected #{expected}", token)
    end

    def adjust_lineno
      @lineno = @line_indexes.bsearch_index { |i| i > @s.pos } - 1
    end

    def seek_to_bof
      return unless @s.skip_until(@re_bof)

      adjust_lineno
    end

    def scan
      loop do
        if @s.bol? && @s.scan(/[ \t]*%options +([^\n]*)/)
          yaml = @s[1].strip
          @s.scan(/\n([^\n]*)/) and yaml << @s[1].strip while yaml.chomp!("\\")
          begin
            @options.update(YAML.safe_load("{ #{yaml} }"))
          rescue Psych::SyntaxError
            raise ParseError.new("Invalid %options", make_token(:OPTIONS, "", nil, @s.pos))
          end
        elsif @s.scan(/[ \t]+|#[^\n]*/)
          # SKIP
        elsif @s.scan(/\n/)
          @lineno += 1
        else
          break
        end
      end

      pos = @s.pos
      if @s.eos? || @s.scan(@re_eof)
        make_token(:EOF, "", nil, pos)
      elsif @s.scan(/('[^']+')/)
        make_token(:TERM, @s[1], nil, pos)
      elsif @s.scan(/([A-Z][A-Z_0-9]*)(-([a-z_][a-z_0-9]+))?/)
        make_token(:TERM, @s[1], Var.new(@s[3]), pos)
      elsif @s.scan(/(<([^>]+|\\.)+>\s*)?([a-z_][a-z_0-9]*)(-([a-z_][a-z_0-9]+))?/)
        make_token(:NONTERM, @s[3], Var.new(@s[5], @s[2]), pos)
          .tap { adjust_lineno if @s[1] }
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
      @token ||= scan
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
      @s.pos = @lineno = 0
      @lines = (src.chomp << "\n\n").lines.to_a # for EOF
      @line_indexes = @lines.reduce([0]) { |r, i| r << (r.last + i.bytesize) }
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
      token.dup.tap do |it|
        it.var = nil
        it.value = "_#{it.value}_#{suffix}"
      end
    end

    def add_rules(lhs, init, rhs_list, inline_option)
      lhs = make_inline_lhs(lhs, inline_option) if inline_option

      @table.descs[lhs.value] = lhs.var&.desc
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
    attr_reader :rules, :start, :nonterms, :terms, :inits, :descs

    def initialize
      @rules = []
      @inits = {}
      @descs = {}
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

      verify
      make_first
      make_follow
      make_table
      LLPg.debug "summary" do
        @nonterms.each_key { |i| warn make_summary(i) unless inline_option(i) }
      end

      self
    end

    def verify
      @terms = @rules.map(&:rhs).flatten.reject { |i| nonterm?(i) }.uniq
      unused = @terms.filter { |i| i.match(/[a-z]/) && !i.start_with?("'") }
                     .uniq.sort
      raise "unused nonterm: #{unused.join(', ')}" unless unused.empty?
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

    def follow(sym)
      @follow[sym]
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
        code << indent(rules_code(sym), 2)
        code << "  end"
      end

      code << "end\n"
      code.join("\n")
    end

    def rules_code(sym, wrapper: nil, rep: nil)
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
            code << indent(matcher_code(rule, x, rep))
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

    def matcher_code(rule, x, rep)
      code = []
      sym = rule.rhs[x]
      var = rule.vars[x]&.name || "_"
      action = rule.actions[x]
      assign = "#{var} = val[#{x}] ="

      if rep && sym == rule.lhs
        code << "redo"
      elsif @table.inline_option(sym)
        code << assign
        code << indent(inline_code(sym), 1)
      elsif @table.nonterm?(sym)
        code << "#{assign} #{method_name(sym)}(st)"
      else
        code << "token.kind == #{sym_in_code(sym)} or parse_error(token, #{sym.inspect})" if x > 0
        code << "#{assign} consume_token"
      end
      code << adjust_action(action) if action
      code
    end

    def inline_code(sym)
      rule = @table.nonterms[sym][0]
      inner = rule.rhs[0]

      body =
        case rule.inline_option
        when :group
          rules_code(sym)
        when :opt
          if @table.inline_option(inner) == :group
            rules_code(inner, wrapper: sym)
          else
            rules_code(sym)
          end
        when :rep
          if @table.inline_option(inner) == :group
            rules_code(inner, wrapper: sym, rep: true)
          else
            rules_code(sym, rep: true)
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
    module Strict
      def setup_code
        code = []
        code << "func (p #{parser_type}) _parse() #{value_type} {"
        code << "    res := p.#{method_name(@table.start)}()"
        code << "    if (p.PeekToken().Kind != tkEOF) {"
        code << "         p.ErrorUnexpected(\"EOF\")"
        code << "    }"
        code << "    return res"
        code << "}"
      end

      def error_handled
        "\x09"
      end

      def handle_error(*)
        block_given? ? yield : "\x09"
      end

      def seq(x = nil)
        if x
          [
            "for {",
            "    res := &_#{x + 1}; _ = res",
            indent(yield),
            "    break",
            "}"
          ]
        else
          [
            "res := &ret; _ = res",
            yield
          ]
        end
      end
    end

    module Recoverable
      def setup_code
        code = []
        code << "const ("
        code << "    tpEOF = 1 << iota"
        code << "    tpSCANERROR"
        @table.terms.each do |i|
          code << "    #{sym_in_code(i, 'tp')}"
        end
        code << ")"
        code << ""

        code << "func (p #{parser_type}) _parse() #{value_type} {"
        code << "    p.Follow(tpEOF)"
        code << "    res := p.#{method_name(@table.start)}()"
        code << "    if res == SYNTAX_ERROR {"
        code << "        return SYNTAX_ERROR"
        code << "    }"
        code << "    if (p.PeekToken().Kind != tkEOF) {"
        code << "        p.Follow(tpEOF)"
        code << "        p.HandleError(SYNTAX_ERROR)"
        code << "        return SYNTAX_ERROR"
        code << "    }"
        code << "    return res"
        code << "}"
      end

      def error_handled
        "*res = SYNTAX_ERROR"
      end

      def handle_error(rule, x, handle, rep: false, continue: false)
        @label_stack.last[1] = true

        follow = rep ? @table.first(rule.rhs, [], []) : []
        follow = @table.first(rule.rhs[(x + 1)..], [], follow)
        follow &= @recover
        follow = follow.map { |i| sym_in_code(i, "tp") }.join("|")
        body = block_given? ? yield : "\x09"
        label_id = @label_stack.last[0]

        if follow.empty?
          if handle == "SYNTAX_ERROR"
            [
              body,
              error_handled.to_s,
              "break loop#{label_id}#{' // cannot continue' if continue}"
            ]
          else
            [
              body,
              "if #{handle} == SYNTAX_ERROR {",
              "    #{error_handled}",
              "    break loop#{label_id}",
              "}",
              continue ? "continue" : "\x09"
            ]
          end
        else
          s = rule.rhs.dup.tap { |it| it.insert(x + 1, "@") }.join(" ")
          [
            "p.Follow(#{follow}) // #{s}",
            body,
            "if !p.HandleError(#{handle}) {",
            "    #{error_handled}",
            "    break loop#{label_id}",
            "}",
            continue ? "continue" : "\x09"
          ]
        end
      end

      def seq(x = nil)
        if x
          @label_id += 1
          @label_stack << [@label_id, false]
          s = "&_#{x + 1}"
        else
          @label_id = 0
          @label_stack = [[nil, false]]
          s = "&ret"
        end
        body = yield # maybe modify @label_stack
        id, handled = @label_stack.pop

        [
          "#{'// ' unless handled}loop#{id}:",
          "for {",
          "    res := #{s}; _ = res",
          indent(body),
          "    break",
          "}"
        ]
      end
    end

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
      _, section = *s.rest.match(%r{/\*\*:reservedWords:(.+?)\*\*/}m)
      (section || "").strip.lines.each do |line|
        _, a, b, = *line.strip.match(/(\w+):(.+)/)
        next unless a

        options["reserved_words"] << [a, b.strip.split]
      end

      options["symbol_descriptions"] = {}
      _, section = *s.rest.match(%r{/\*\*:symbolDescriptions:(.+?)\*\*/}m)
      (section || "").strip.lines.each do |line|
        _, a, b, = *line.strip.match(/((?:'[^']+'|[^'\s]+)+):(.+)/)
        next unless a

        b.strip.split.each do |i|
          options["symbol_descriptions"][i] = a
        end
      end
    end

    def generate(table, options)
      @table = table
      @options = {
        "parser_type" => "*Parser",
        "value_type" => "Value",
        "token_type" => "*Token"
      }.merge(options)

      @recover = @options["recover"]
      extend(@recover ? Recoverable : Strict)

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
      @expected_syms = []

      @table.rules.each do |rule|
        next unless rule.inline_option == :opt

        rule.actions[0] = "RET _1 // DEFAULT OPT ACTION"
      end

      code << "const ("
      code << "    tkEOF = iota"
      code << "    tkSCANERROR"
      @table.terms.each do |i|
        code << "    #{sym_in_code(i)}"
      end
      code << ")"
      code << ""

      code << "var tokenKinds = map[string]int32{"
      @table.terms.each do |i|
        next unless i[0] == "'"

        code << "    \"#{i[1...-1]}\": #{sym_in_code(i)},"
      end
      code << "}"
      code << ""

      (@table.nonterms.keys + @table.terms).each do |i|
        @table.descs[i] ||= i.tr("_", "-") unless /\A['_]/.match?(i)
      end
      @table.descs.update(@options["symbol_descriptions"])
      code << "var tokenLabels = [...]string{"
      code << "    \"EOF\", // tkEOF"
      code << "    \"SCANERROR\", // tkSCANERROR"
      @table.terms.each do |i|
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

      code << setup_code
      @table.each_key do |sym|
        next if @table.inline_option(sym)

        @expected_stack = []
        collect_expected_syms(Rule.new(sym, [sym], [nil], [nil]), 0)

        code << ""
        code << "/** ** **"
        code << @table.make_summary(sym).lines.map { |i| "  #{i}" }.join
        code << " ** ** **/"
        code << "func (p #{parser_type}) #{method_name(sym)}() (ret #{value_type}) {"
        code << indent(seq { rules_code(sym) })
        code << "    return"
        code << "}"
        raise unless @expected_stack.empty?
      end

      debug "expected symboles" do
        @expected_syms.each_with_index do |(rule, x, syms), i|
          rhs = rule.rhs.dup.tap { |it| it.insert(x, "@") }.join(" ")
          if rule.lhs == rule.rhs[0]
            first = @table.first([rule.lhs], [], [Table::EMPTY])
            mark = first.include?(Table::EMPTY) ? "*(allow empty)" : "*"
            warn format("%03d%s: %s", i, mark, rule.lhs)
          else
            warn format("%03d: %s -> %s", i, rule.lhs, rhs)
          end
          warn "     [#{syms.size}] #{syms.join(' ')}"
        end
      end

      code.join("\n").gsub(/^\s*\x09\n/, "") << "\n"
    end

    def rules_code(sym, inline: nil, wrapper: nil, rep: nil)
      set = @table[sym] # expected => rule
      annotations = []

      code = []
      code << adjust_action(@table.inits[sym]) if @table.inits[sym]
      code << "switch p.PeekToken().Kind {"
      set.group_by(&:last).each do |rule, pairs|
        annotations << "// RULE: #{rule}"
        pat = pairs.map { |i| sym_in_code(i[0]) }.join(", ")
        code << "case #{pat}:"
        if rule.rhs.empty?
          action = rule.actions[0]
          action = action ? adjust_action(action) : "// NOP"
          code << indent(action)
        else
          rule.rhs.size.times do |x|
            code << indent(matcher_code(rule, x, rep, inline))
          end
          code << "" << indent("continue") if wrapper && rep
        end
      end

      skip = []
      if wrapper
        annotations.unshift "// RULE: #{wrapper} -> #{sym}#{rep ? '*' : '?'}"
        skip = @table[wrapper].filter_map { |k, v| v.rhs.empty? && k }
        unless skip.empty?
          pat = skip.map { |i| sym_in_code(i) }.join(", ")
          code << "case #{pat}:"
          code << "    // NOP"
        end
      end

      code << "default:"
      code << indent(detect_error)
      if rep
        rule, x = inline
        code << indent(handle_error(rule, x - 1, "SYNTAX_ERROR", continue: true))
      else
        code << indent(error_handled)
      end
      code << "}" # end switch

      annotations + code
    end

    def matcher_code(rule, x, rep, _inline)
      sym = rule.rhs[x]
      var = rule.vars[x]&.name || "_#{x + 1}"
      action = rule.actions[x]
      sym_inline = @table.inline_option(sym)

      rhs_comment = rule.rhs.dup.tap { |it| it.insert(x, "@") }.join(" ")
      code = ["// #{rhs_comment}"]
      if rep && sym == rule.lhs
        code << "continue"
      elsif sym_inline
        opt_term = sym_inline == :opt && @table.term?(@table.nonterms[sym][0].rhs[0])
        var_type = opt_term ? token_type : value_type

        collect_expected_syms(rule, x)
        code << "var #{var} #{var_type}; _ = #{var}"
        code << handle_error(rule, x, var, rep: rep) do
          inline_code(rule, x)
        end
      elsif @table.nonterm?(sym)
        code << handle_error(rule, x, var, rep: rep) do
          "#{var} := p.#{method_name(sym)}(); _ = #{var}"
        end
      elsif x == 0 # term?(sym)
        code << "#{var} := p.ConsumeToken(); _ = #{var}"
      else # term(sym) && x > 0
        collect_expected_syms(rule, x)
        code << "#{var} := p.PeekToken(); _ = #{var}"
        code << "if #{var}.Kind == #{sym_in_code(sym)} {"
        code << "    p.ConsumeToken()"
        code << "} else {"
        code << indent(detect_error)
        code << indent(handle_error(rule, x, "SYNTAX_ERROR", rep: rep))
        code << "}"
      end
      code << adjust_action(action) if action
      code << "" if x < rule.rhs.size - 1
      code
    end

    def inline_code(rule, x)
      sym = rule.rhs[x]
      expanded = @table.nonterms[sym][0]
      inner = expanded.rhs[0]
      wrapper = rep = nil

      case expanded.inline_option
      when :group
        # OK
      when :opt
        if @table.inline_option(inner) == :group
          wrapper = sym
          sym = inner
        end
      when :rep
        rep = true
        if @table.inline_option(inner) == :group
          wrapper = sym
          sym = inner
        end
      else
        raise "invalid inline option"
      end

      seq(x) { rules_code(sym, inline: [rule, x], wrapper: wrapper, rep: rep) }
    end

    def collect_expected_syms(rule, x)
      collect = lambda do |rule, x, acc|
        rule.rhs[x..].each_with_index do |sym, offset|
          if (option = @table.inline_option(sym))
            case option
            when :group
              @table.nonterms[sym].each { |rule| collect.call(rule, 0, acc) }
              break
            when :opt, :rep
              inner_rule = @table.nonterms[sym][0]
              inner = inner_rule.rhs[0]
              if @table.inline_option(inner) == :group
                collect.call(inner_rule, 0, acc)
              else
                acc << sym_desc(rule, x + offset, inner)
              end
            else
              raise
            end
          else
            acc << sym_desc(rule, x + offset, sym)
            first = @table.first([rule.rhs[x + offset]], [], [Table::EMPTY])
            break unless first.include?(Table::EMPTY)
          end
        end
      end

      syms = [].tap { |it| collect.call(rule, x, it) }.uniq
      first = @table.first(rule.rhs[x..], [], [Table::EMPTY])
      syms << "etc." if first.include?(Table::EMPTY)
      @expected_stack << syms
      @expected_syms << [rule, x, syms]
    end

    def detect_error
      expected = @expected_stack.pop.join(", ").inspect
      "p.ErrorUnexpected(#{expected})"
    end

    def adjust_action(action)
      action.strip.gsub(/(^| )RET +/, "\\1*res = ")
    end

    def method_name(sym)
      parts = sym.to_s.split("_")
      "_parse#{parts.map(&:capitalize).join}"
    end

    def sym_desc(rule, x, sym)
      desc = rule.vars[x]&.desc if rule
      desc || @table.descs[sym] || sym
    end

    def sym_in_code(sym, prefix = "tk")
      sym = sym[1..-2].chars.map { |i| ASCII_SYMBOL_MAP[i] }.join if sym[0] == "'"
      "#{prefix}#{sym}"
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

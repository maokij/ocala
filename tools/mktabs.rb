#!/usr/bin/env ruby
# frozen_string_literal: true

require "strscan"

class Mktab
  def generate_tabs(path)
    output_path = File.expand_path("../tabs.g.go", path)

    tabs = {}.tap do |it|
      s = StringScanner.new(File.read(path))
      until s.eos?
        line = s.scan_until(/\n|\z/).strip
        _, kind, name = *line.match(%r{// +(SYNTAX|FUN|SPECIAL): +\(([^\s)]+)})
        next unless kind

        line = s.scan_until(/\n|\z/).strip
        _, go_name = *line.match(/.+?Compiler\w*\) +([sf]\w+)/)
        it[kind] ||= {}
        it[kind][name] = go_name if go_name
      end
    end

    code = ["package core", ""]
    tabs.each do |kind, m|
      code << "var #{kind.capitalize}Map = map[*Keyword]SyntaxFn{"
      m.each do |k, v|
        code << %!    Intern("#{k}"): SyntaxFn((*Compiler).#{v}),!
      end
      code << "}" << ""
    end
    File.write(output_path, gofmt(code.join("\n")))
  end

  def generate_ocalajson(paths)
    require "json"
    outpath = paths.pop
    include, examples = paths.partition { |i| i.include?("/ocala/include/") }
    include = include.map do |i|
      [File.basename(i), File.read(i)]
    end
    examples = examples.map do |i|
      [i.gsub(%r{\A(examples|internal)/}, ""), File.read(i)]
    end
    data = {
      incMap: include.sort_by(&:first).flatten,
      examples: examples.sort_by(&:first).to_h
    }
    File.write(outpath, JSON.generate(data))
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
    when "tabs"
      generate_tabs(args[0])
    when "ocalajson"
      generate_ocalajson(args)
    else
      raise "unknown subcommand: #{command}"
    end
  end
end

Mktab.new.run(ARGV.dup)

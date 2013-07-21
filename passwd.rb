#!/usr/bin/ruby
#
# This script will parse the /etc/passwd file and produce the output like
#ENTRY:
#    user: foo
#    password: x
#    uid: 1172
#    gid: 100
#    gecos: foo
#    home: /home/foo
#    shell: /bin/bash

class PasswdFile
  class PasswdEntry
    def initialize
      @attributes = {}
    end

    def method_missing(name, *args)
      attribute = name.to_s
      if attribute =~ /=$/
        @attributes[attribute.chop] = args[0]
      else
        @attributes[attribute]
      end
    end

    def import(line)
      field_names = [:user, :password, :uid, :gid, :gecos, :home, :shell]
      fields = line.split(":")
      puts "Error: Not a passwd line, not 7 fields" unless fields.length == 7
      field_names.each do |field_name|
        self.send("#{field_name}=", fields[field_names.index(field_name)])
      end
    end
  end

  attr_reader :filename
  def initialize(filename="/etc/passwd")
    @entries = []
    import(filename)
  end

  def import(filename)
    File.open(filename, "r").each do |line|
      entry = PasswdEntry.new
      entry.import line
      @entries << entry
    end
  end

  def each
    @entries.each do |entry|
      yield entry
    end
  end

  def print
    field_names = [:user, :password, :uid, :gid, :gecos, :home, :shell]
    
    @entries.each do |entry|
      puts "ENTRY:"
      field_names.each do |field_name|
        puts "    #{field_name}: " + entry.send("#{field_name}")
      end
    end
  end
end
entries = PasswdFile.new
entries.print

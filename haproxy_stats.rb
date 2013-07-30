#!/usr/bin/env ruby

# Add the below lines to the types.db
#hproxy_status          status:GAUGE:0:U
#haproxy_traffic         stot:COUNTER:0:U, eresp:COUNTER:0:U, chkfail:COUNTER:0:U
#haproxy_sessions        qcur:GAUGE:0:U, scur:GAUGE:0:U


require 'optparse'
require 'fileutils'
require 'socket'
require 'pp'
require 'socket'


# default options
options = {
  :socket => "/var/run/haproxy/haproxy.sock",
  :wait_time => 60
}

STATSCOMMAND = "show stat\n"
BACKENDSTRING = "BACKEND"

class HaproxyStats
  attr_accessor :haproxy_vars
  
  def initialize(section,name)
    @columnTypes = nil
    @section = section
    @name = name
    
    # typei = type-instance
    
    @haproxy_vars = {
      "haproxy_status" => {
        :typei => [:server],
        :data_desc => ["status"]
      },
      "haproxy_traffic" => {
        :typei => [:server, :total],
        :data_desc => ["stot","eresp", "chkfail"]
      },
      "haproxy_session" => {
        :typei => [:server, :total],
        :data_desc => ["qcur","scur"]
      }
    }
    @translate = { "status" => [
      [/^UP$/, "2"],
      [/^UP.$/, "-1"], # going down
      [/^DOWN$/, "2"], 
      [/^DOWN.*$/, "1"], # going up
      [/^no check$/, "0"] 
      ]
    }
  end
  
  def parse(input,time)
    output = ""
    backend_line = nil
    accumulator = {}
    
    if @columnTypes == nil
      parseColumnTypes(input)
    end
    
    # init accumulator
    
    @haproxy_vars.each do |type, data|
      if data[:typei].include?(:total)
        accumulator[type] = [].fill(0, 0, data[:data_desc].length)
      end
    end
    
    input.each_with_index do |line, index|
      if index == 0
        next
      end
      
      if line =~ /^#{@section}/
        values = line.split(',')
        if values[@columnTypes["svname"]] == BACKENDSTRING
          backend_line = line.clone
        else
          @haproxy_vars.each do |type, data|
            if data[:typei].include?(:server)
              output << "PUTVAL #{Socket.gethostname}/haproxy/haproxy-#{@name}/" + type + "-" + values[@columnTypes["svname"]].downcase.gsub(/-/,'_') + " interval=#{time} N"
              data[:data_desc].each_with_index do |column,index|
                if @translate[column] != nil
                  output << ":"
                  @translate[column].each do |pattern,val|
                    if values[@columnTypes[column]] =~ pattern
                      output << val
                      break
                    end
                  end
                else
                  if values[@columnTypes[column]] == ""
                    output << ":0"
                  else
                    if data[:typei].include?(:total)
                      accumulator[type][index] += values[@columnTypes[column]].to_i
                    end
                    output << ":" << values[@columnTypes[column]]
                  end
                end
              end
              output << "\n"
            end
          end
        end
      end
    end
    
    values = backend_line.split(",")
    if values[@columnTypes["svname"]] != BACKENDSTRING
      raise "Unable to find the backend string"
    end
    
    #handle backend string
    
    @haproxy_vars.each do |type,data|
      if data[:typei].include?(:total)
        output << "PUTVAL #{Socket.gethostname}/haproxy/haproxy-#{@name}" + type + "-total interval=#{time} N"
        data[:data_desc].each_with_index do |column, index|
          if values[@columnTypes[column]] == ""
            output << ":" << accumulator[type][index].to_s
          else
            output << ":" < values[@columnTypes[column]]
          end
        end
        output << "\n"
      end
    end
    
    output
  end
  
  private 
  def parseColumnTypes(input)
    if input.length == 0
      raise "Empty input not allowed"
    end
    
    match = input.index("\n") 
    if match.nil? || match == 0
      raise "Invalid input: No line breaks found"
    end
    
    line = input[0..match]
    
    if line[0].chr != "#"
      raise "Invalid input: No column types found"
      end
      
      @columnTypes = {}
      line.slice!(0..1)
      line.chomp!
      columns = line.split(',')
      
      # Build hash
      columns.each_with_index do |column, index|
        if !column.nil? && column != ""
          @columnTypes[column] = index 
        end
      end
      
      @columnTypes
    end
  end
  
  opts = OptionParser.new
  
  opts.banner = "Usage: haproxy-stats.rb [options]"
  
  opts.separator ""
  opts.separator "Specific option:"
  opts.on("-sSOCKET","--socket = SOCKET", "Location of the haproxy socket", "Default: /var/run/haproxy/haproxy.sock") do |str|
    options[:socket] = str
  end
  opts.on("-wWAITING", "--wait = WAITING", "Time to wait between samples", "Default: 60") do |str|
    options[:wait_time] = str.to_i
  end
  opts.on("-eSECTION", "--section = SECTION", "HAProxy section to search for stats") do |str|
    options[:section] = str
  end
  opts.on("-nNAME", "--name = NAME", "Config name for haproxy stats (ie www)") do |str|
    options[:name] = str
  end
  opts.separator ""
  opts.separator "Common options:"
  opts.on_tail("-h", "--help", "Show this message") {
    exit
  }
  
  begin
    opts.parse(ARGV)
    if ARGV.length == 0
      exit
    end
    
    # Check for required args
    raise "Section is requried " unless options[:section]
    raise "Name is required" unless options[:name]
    
  rescue SystemExit
    puts opts
    exit
  rescue Exception => e
    puts "Error: #{e}"
    puts opts
    exit
  end
  
  # Option Done
  
  begin
    socket_data = ""
    stats = HaproxyStats.new(options[:section],options[:name])
    
      socket = UNIXSocket.open(options[:socket])
      socket.write(STATSCOMMAND)
      socket_data = socket.read
      socket.close
      
      puts stats.parse(socket_data,options[:wait_time])
      
  rescue Exception => e
    puts e
    puts e.backtrace
  ensure
    socket.close if !socket.nil? && !socket.closed?
  end

#!/usr/bin/ruby

require 'rubygems'
require 'optparse'
require 'yaml'
require 'json'
require 'net/http'
require 'uri'
require 'pp'

class CheckGenerationError < StandardError
end

class CheckGenerationTime
  def self.check_generation_time(host,genTime,maxTime)
    msg = "CRITICAL: Generation Time for RR YAML file #{host}"
    status = 2
    begin
    checkTime = Time.now.utc.to_i - genTime
    if checkTime > maxTime.to_i
      msg = "CRITICAL:#{host}  generation Time " +  Time.at(genTime).to_s.split(" ")[3] + " > max Time #{maxTime}"
     else
       msg = "OK : #{host}  generation Time " + Time.at(genTime).to_s.split(" ")[3] + " < max Time #{maxTime}"
       status = 0
    end
  rescue Exception => e
    msg = "CRITICAL : #{host} check error #{e.message}"
    status = 2
    end
  [ status, msg ]
  end

  def self.parse(args)
    options = {}
    opts = OptionParser.new do |opts|
      opts.banner = "Usage: #{__FILE__} [options]"
      opts.separator ""
      opts.on("-t", "--maxTime TIME", "Specify Maximum time to check for") do |t|
        options[:mTime] = t
      end
      opts.on_tail("-h","--help","Show this messsage") do
        puts opts
        exit
      end
    end
      begin
        opts.parse!(args)
        raise OptionParser::MissingArgument, "-t , no max time specified." if options[:mTime].empty?
        rescue SystemExit
          exit
        rescue Exception => e
          puts "Error " + e, "#{__FILE__} -h for options"
          exit
        end
        options
  end
      def self.run(args)
        opts = parse(args)
        nag_msg = []
        status = 0
        begin
         config_file = "/opt/inmobi/logs-manifest/conf/rr-puller.yaml"
         # load the stats
         stats = YAML.load_file(config_file)
         primary_end_point = stats["master"]["primary"]
         secondary_end_point = stats["master"]["secondary"]
         dir_tocheck = stats["destdir"]
         if primary_end_point
           @host_list = Net::HTTP.get URI.parse(primary_end_point)
          else
            @host_list = Net::HTTP.get URI.parse(secondary_end_point)
        end
         hosts = JSON.load(@host_list)
        hosts.each do |h|
          grepString = open("#{dir_tocheck}#{h}-manifest") { |f| f.lines.find { |line| line.include?("generate_time") } }
          genTime = (grepString.chomp!).split(": ").last.to_i
          st, msg = check_generation_time(h,genTime,opts[:mTime])
          status = 2 if st != 0
          nag_msg << msg
        end
        puts nag_msg.join("\n")
        exit status
      end
    end
end
CheckGenerationTime.run(ARGV)



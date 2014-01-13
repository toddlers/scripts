#!/usr/bin/env ruby

require 'pp'
require 'yaml'
require 'optparse'


# Example to take config params from both 
# commandline as well from a config file

class TEST
  def self.parse(args)
    options = {}
    opts = OptionParser.new do |opts|
      opts.banner = "Usage: #{__FILE__} [options]"
      opts.separator ""
      opts.on("-p","--provider PROVIDER", "Provider name") do |p|
        options[:provider] = p
      end
      opts.on("-r", "--region REGION", "provide region name") do |r|
        options[:region] = r
      end
      opts.on("-c","--config FILE", "Read options from file") do |file|
        fp = YAML::load(File.open('config1.yaml'))
        options[:provider] =  fp["provider"]
        options[:region] = fp["aws_region"]
      end

      opts.on_tail("-h","--help","Show this message") do
        puts opts
        exit
      end
    end
    begin
      opts.parse!(args)
      raise OptionParser::MissingArgument, "-p, no provider specified" if not options[:provider]
      raise OptionParser::MissingArgument, "-r, no region specified" if not options[:region]
    rescue SystemExit
      exit
    rescue Exception => e
      puts "Error " + e , "#{__FILE__} -h for options"
      exit
    end
    options
  end

  def self.run(args)
    opts = parse(args)
    pp opts
  end
end

TEST.run(ARGV)


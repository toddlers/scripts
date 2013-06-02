#!/usr/bin/ruby

# Used to check the status of Puppet

require 'rubygems'
require 'yaml'
require 'optparse'
require 'pp'
require 'yaml'

class CheckPuppetError < StandardError
end

class CheckPuppet
  def self.check_last_run_time(lastrun_time,maxTime)
    msg = "CRITICAL: Puppet was last run #{lastrun_time}"
    status = 2
    begin
    checkTime = Time.now.utc.to_i - lastrun_time
    maxMin = maxTime.to_i/60
    if checkTime > maxTime.to_i
      msg = "CRITICAL:  Puppet was last run " +  Time.at(lastrun_time).to_s.split(" ")[3] + " > max Time #{maxMin} min"
     else
       msg = "OK : Puppet was last run " + Time.at(lastrun_time).to_s.split(" ")[3] + " < max Time #{maxMin} min"
       status = 0
    end
  rescue Exception => e
    msg = "CRITICAL : Puppet check error #{e.message}"
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
     puppet_summary = "/var/lib/puppet/state/last_run_summary.yaml"
     if File.exist?(puppet_summary)
     # get the last run time
     stats = YAML.load_file(puppet_summary)
     last_run = stats["time"]["last_run"]
   else
     puts "#{puppet_summary} file does exists"
   end
      st, msg = check_last_run_time(last_run,opts[:mTime])
      status = 2 if st != 0
      nag_msg << msg
    puts nag_msg.join("\n")
    exit status
  end
end 
end
CheckPuppet.run(ARGV)

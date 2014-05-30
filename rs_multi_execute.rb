#!/usr/bin/env ruby
require 'optparse'
require 'right_api_client'
require 'pp'

class Deploy
  
  def self.parse(args)
    options = {}
    opts = OptionParser.new do |opts|
      opts.banner = "Usage: #{__FILE__} [options]"
      opts.separator ""
      opts.on("-a","--array SERVER_ARRAY","Server Array Name") do |a|
        options[:sarray] = a
      end
      opts.on("-s","--script SCRIPT_ID","RS Script ID") do |s|
        options[:script] = s
      end
      opts.on("-i","--account_id ACCOUNT_ID","RS Environment Account id") do |i|
        options[:account_id] = i
      end
      opts.on("-u","--user USERNAME","RS Environment Account user name") do |u|
        options[:user] = u
      end
      opts.on("-p","--passwd PASSWORD","RS Environment Account password") do |p|
        options[:passwd] = p
      end
      opts.on_tail("-h","--help","Show this message") do
        puts opts
        exit
      end
    end
    begin
      opts.parse!(args)
      raise OptionParser::MissingArgument " -a no server array specified" if not options[:sarray]
      raise OptionParser::MissingArgument " -s no script specified" if not options[:script]
      raise OptionParser::MissingArgument " -i no account id  specified" if not options[:account_id]
      raise OptionParser::MissingArgument " -u no username  specified" if not options[:user]
      raise OptionParser::MissingArgument " -p no password specified " if not options[:passwd]
    rescue SystemExit
      exit
    rescue Exception => e
      puts e
      exit
    end
    options
  end


  def self.run(args)
    opts = parse(args)
    pp opts
    server_array = opts[:sarray]
    account_id = opts[:account_id]
    script = "right_script_href=/api/right_scripts/" + opts[:script]
    pp script
    username = opts[:user]
    password = opts[:passwd]
    pp username
    pp password
    begin
      @client = RightApi::Client.new(:email => username , :password => password, :account_id => account_id)
      task = @client.server_arrays.index(:filter => ["name==#{server_array}"])[0].multi_run_executable(script)
      puts task.show.summary
    rescue Exception => e
      puts e
      exit
    end
  end
end
Deploy.run(ARGV)

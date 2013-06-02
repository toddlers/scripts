#!/usr/bin/env ruby
require 'json'
require 'rest-client' 
require 'pp'
require 'optparse'
options = {}
opt_parser = OptionParser.new do |opt|
  opt.banner = "Usage json_parser [OPTIONS]"
  opt.separator ""
  opt.separator "Options"
  opt.on("-e","--endpt ENDPOINT","Provide the json endpoint or URL") do |l|
    options[:endpoint] = l
  end
  opt.on("-h","--help","help") do
    puts opt_parser
  end
end
opt_parser.parse!
json_input_url = options[:endpoint]
# The killer one line ruby json parsing! 
pp JSON.parse(RestClient.get(json_input_url)) if json_input_url

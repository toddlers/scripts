#!/usr/bin/env ruby
#
# Downloading xkcd images 
#
require 'nokogiri'
require 'pp'
require 'open-uri'
require 'optparse'
require 'fileutils'


class XkcdDownload

    def self.parse(args)
        options = {}
        opts = OptionParser.new do |opts|
            opts.banner = "Usage: #{__FILE__} [options]"
            opts.separator ""
            opts.on("-s", "--start NUM", "enter the starting comic page no. ") do |s|
                options[:start] = s
            end
            opts.on("-e","--end NUM","enter the last comic page no. :") do |e|
                options[:pend] = e
            end

            opts.on_tail("-h", "--help", "Show this message") do
                puts opts
                exit
            end
        end
        begin
            opts.parse!(args)
            raise OptionsParser::MissingArgument "-s , no starting comic page no." if options[:start].empty?
            raise OptionsParser::MissingArgument "-e , no end comic page no." if options[:pend].empty?
        rescue SystemExit
            exit
        rescue Exception => e
            puts "Error" +e, "#{__FILE__} -h for options"
            exit
        end
        options
    end

    # This method will return the hash of urls between the mentioned page
    # numbers
    def self.get_urls(start,pend)
        list_url = []
        puts "fetching urls ...."
        (start .. pend+1).each do |comic_num|
            page = Nokogiri::HTML(open("http://xkcd.com/#{comic_num}"))
            list_url << page.css('div img')[1].attributes["src"].value
        end
            return list_url
    end

    def self.download_comic(urls,start)
        count = start
        urls.each do |url|
            puts "Downaloading xkcd " + count.to_s
            filename = "xkcd" + count.to_s
            File.open(filename, 'wb') do |save_file|
                open(url,'rb') do |read_file|
                    save_file.write(read_file.read)
                    count += 1
                end
            end
        end
    end

    def self.run(args)
        opts = parse(args)
        start = opts[:start].to_i
        pend = opts[:pend].to_i
        path = "#{Dir.pwd}/xkcd"
        FileUtils.mkdir_p(path) unless File.exists?(path)
        Dir.chdir(path)
        puts "Saving files to #{Dir.pwd}"
        begin
        urls = get_urls(start,pend)
        download_comic(urls,start)
        rescue Exception => e
            puts e 
            exit
        end
    end

end

XkcdDownload.run(ARGV)

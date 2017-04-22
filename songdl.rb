#!/usr/bin/env ruby

require 'pp'
require 'json'
require 'net/http'
require 'open-uri'
require 'nokogiri'
require 'optparse'
require 'fileutils'

class Songdl

  def self.parse(args)
    options = {}
    opts = OptionParser.new do |opts|
      opts.banner = "Usage: #{__FILE__} [options]"
      opts.separator ""
      opts.on("-s", "--search MOVIE_NAME", "Enter the movie name with year") do |s|
        options[:search] = s
      end
      opts.on_tail("-h", "--help", "Show this message") do
        puts opts
        exit
      end
    end
    begin
      opts.parse!(args)
      raise OptionParser::MissingArgument "-s , no search string provided " if options[:search].empty?
    rescue SystemExit
      exit
    rescue Exception => e
      puts "Error" +e, "#{__FILE__} -h for options"
      exit
    end
    options
  end


  def self.googleSearch(params)
    puts "searching for #{params[:q]}"
    uri = URI("https://www.googleapis.com/customsearch/v1?")
    uri.query = URI.encode_www_form(params)
    http = Net::HTTP.new(uri.host,uri.port)
    http.use_ssl = true
    response = http.get(uri.request_uri)
    result = JSON.parse(response.body)
    links = []
    result['items'].each do |i|
      links << i['link']
    end
    getUrls(links)
  end

  def self.getUrls(links)
    flinks = []
    links.each do |l|
      page = Nokogiri::HTML(open(l))
      extractLinks = page.css('a')
      extractLinks.each do |e|
        if e["href"] =~ /^http.*\.mp3$/
          flinks << e["href"]
         end
      end
    end
    return flinks
  end

  def self.downloadSongs(urls)
    urls.each do |url|
      filename = url.split("/").last
      File.open(filename,'wb') do |save_file|
        puts "saving file #{filename}"
        open(url,'rb') do |read_file|
          save_file.write(read_file.read)
        end
      end
    end
  end

  def self.run(args)
    opts = parse(args)
    query = opts[:search] + " download mp3 songs"
    params = {
      :key => 'AIzaSyAxoxKFaWHdn3otGrzggUr1Oa2gvhLUOpc',
      :cx => '006551824314330792304:xompbxiw3xa',
      :q => "#{query}",
      :page => '1'
    }
    name, year = "#{query}".split(" ")
    dirname = "#{name}" + "_" + "#{year}"
    path = "#{Dir.pwd}/#{dirname}"
    FileUtils.mkdir_p(path) unless File.exists?(path)
    Dir.chdir(path)
    puts "Saving songs to #{Dir.pwd}"
    begin
      slinks = googleSearch(params)
      downloadSongs(slinks)
    rescue Exception => e
      puts e
      exit
    end
  end
end

Songdl.run(ARGV)

#!/usr/bin/env ruby
require 'json'
require 'net/http'
require 'pp'
require 'open-uri'

def get_profile_pic
    #return a dict of url of all the profile pics of your friend on facebookk
    params = { 
            :q => "SELECT pic from user WHERE uid IN (SELECT uid2 FROM friend WHERE uid1 = me()) LIMIT 1000", 
            :access_token => 'CAACEdEose0cBAKZCjkZCIVcIO6K2y9TSSWeoA4GhjvPk198ZBAxvYXYzxtXmeTrW1woJZAs0R87WZBC93sI7D46dxIpLTGKyw9QcG2pZC2DObLL33PLyuD41HZBTSgaMNp0bAxEE7GzKvKk8mljQYwOq0MrOmfK27nlDo48Y8Qm6IUJ4IJM27I3JE3x3xolBz0ZD' 
             }
    uri = URI('https://graph.facebook.com/fql')
    uri.query = URI.encode_www_form(params)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    response = http.get(uri.request_uri)
    result = JSON.parse(response.body)
    return result['data']
end

#download all the pics

def save_photos(photos)
    count = 0
    photos.each do |p|
        puts "successfully downloaded cover pic" + count.to_s
        filename = "pic" + count.to_s
        File.open(filename,'wb') do |save_file|
            open(p['pic'],'rb') do |read_file|
                save_file.write(read_file.read)
                count += 1
            end
        end
    end
end

save_photos(get_profile_pic)


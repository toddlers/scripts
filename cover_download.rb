#!/usr/bin/env ruby
require 'json'
require 'net/http'
require 'pp'
require 'open-uri'

def get_profile_pic
    #return a dict of url of all the profile pics of your friend on facebookk
    params = { 
            :q => "SELECT src FROM photo_src WHERE photo_id IN (SELECT object_id FROM photo WHERE aid IN (SELECT aid FROM album WHERE owner IN (SELECT uid2 FROM friend WHERE uid1= me()))) AND width > 1000  LIMIT 1000", 
            :access_token => 'CAACEdEose0cBAOxogR9F9oT1CBZAIdGxYKxnj4p5m4Pjtmugpo70Fs0dn3CcXZCNoJuKPVoJs23fVEfcKrdpLSmdmRpwbOrZBLpzNjqydZBEkXtnhaXjJ4M1ryIQyLEAIYq5PdOZAlWdZCB5r2foUfykYSB0znOBjBhuU3UMHbxEAYuzYDc6RLtUCHFJJeJ4EZD'
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
            open(p['src'],'rb') do |read_file|
                save_file.write(read_file.read)
                count += 1
            end
        end
    end
end

save_photos(get_profile_pic)


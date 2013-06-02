#!/usr/bin/env ruby
require 'net/https'
require 'net/smtp'
require 'optparse'
require 'socket'

options={}
optparse=OptionParser.new do |opts|
  opts.banner="Usage:nagios_report.rb [options]"
  options[:username]="nagios"
  opts.on('-u','--username UNAME','Specify the nagios username') do |uname|
          options[:username]=uname
          end
  options[:passwd]="passwd"
  opts.on('-p','--passwd PASSWD','Specify the password name') do |p|
    options[:passwd]=p
  end
  opts.on('-h','--help','Usage Info') do
    puts optparse
    exit
  end
end
optparse.parse!

 def colo(h)
  c = nil
   if h =~ /\.([[:alnum:]]+)\.inmobi\.com$/
    c=$1
   end
   c
 end

  def mnode(h)
    colo_name = colo(h)
    case colo_name
    when "colo_name"
      "hostnames"
    when "colo_name"
      "hostname"
    when "colo_name"
      "hostname"
    when "colo_name"
      "hostname"
    end
end

h = Socket.gethostname
monserver = mnode(h)
puts monserver
# If you're accessing through SSL:
#http = Net::HTTP.new("nagios1.domain.com", 443)
#http.use_ssl = true
# Otherwise:
http = Net::HTTP.new(monserver,9999)
http.use_ssl = false

http.start do |http|
#  req = Net::HTTP::Get.new("/nagios/cgi-bin/status.cgi?host=all",{"User-Agent" => "nagios_checker"})
  req = Net::HTTP::Get.new('/nagios/cgi-bin/status.cgi?host=all&limit=200&limit=0',{"User-Agent" => "nagios_checker"})
  req.basic_auth("#{options[:username]}", "#{options[:passwd]}")
  response = http.request(req)
  puts response.code
  resp = response.body
  host_problems_flag = service_problems_flag = false

  msg = ""
  msg << "From: Nagios <nagios@#{monserver}>\n"
  msg << "To: Adserve Ops  <adserve-ops@inmobi.com>\n"

  status = resp.scan(/Notifications are disabled/)
  if status.size > 0
    msg << "Subject: Nagios Status: WARNING!\n"
    msg << "\n"
    msg << "WARNING: ALL notifications are DISABLED!\n"
  else
    msg << "Subject: Nagios Status: #{h} \n"
    msg << "\n"
  end

  host = service = nil
  hosts = resp.scan(/A HREF='extinfo.cgi\?type=1&(.*?)'><IMG SRC=.*ALT='Notifications/)
  msg << "\n*** Notifications are disabled for the following hosts ***\n"
  hosts.each do |h|
    h.to_s.split("&").each do |e|
      key, value = e.split("=")
      msg << "Host: #{value}\n"
    end
    host_problems_flag = true
  end

  services = resp.scan(/A HREF='extinfo.cgi\?type=2&(.*?)'><IMG SRC=.*ALT='Notifications/)
  msg << "\n*** Notifications are disabled for the following services ***\n"
  services.each do |s|
    s.to_s.split("&").each do |e|
      key, value = e.split("=")
      if key == 'host'
        host = value
      elsif key == 'service'
        service = value.gsub('+', ' ')
      end 
    end
    msg << "Host: #{host}, Service: #{service}\n"
    service_problems_flag = true
  end

  if host_problems_flag || service_problems_flag
    smtp = Net::SMTP.start('mail.boo.com', 25)
    smtp.send_message msg, "nagios@#{monserver}", 'foo@boo.com'
    smtp.finish
  end

end

#!/usr/bin/ruby

#
# Query ingrapher clusters
#
# $ qinrg
# <cluster list>
# $ qinrg cluster-1
# <hosts in cluster-1
#
# define INGRAPHER_ROOT environment for this to work
# points to the ingrapher clusters directory
#

require 'pp'
require 'yaml'

def yaml_hosts(file)
  h = YAML.load_file file
  h['host']
end

ingr_root = ENV['INGRAPHER_ROOT']
ingr_root = "ingrapher_location" if not ingr_root

cluster_files = Dir.entries(ingr_root).collect { |e|
  File.expand_path(e, ingr_root)
}

cluster_files = cluster_files.reject { |f|
  not f.end_with? '.yaml'
}

clusters = {}

cluster_files.each { |f|
  cluster = File.basename f, '.yaml'
  clusters[cluster] = f
}

if ARGV.length == 0
  clusters.sort.each do |k,v| puts k end
end

hosts = []

ARGV.each do |q|
  if clusters[q]
    hosts.push yaml_hosts(clusters[q])
  else
    $stderr.puts "Invalid cluster #{q}"
  end
end

hosts.flatten!

hosts = hosts.reject { |host|
  host !~ /\.inmobi\.com$/
}

puts hosts.sort.uniq

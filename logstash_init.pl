#!/usr/bin/perl

#What does this script
#It takes the configuration filename as input, and writes the logstash service script, right into /etc/init.
#How to use it
#First of all, you must have a /opt/logstash/logstash.jar file (which can be a symbolic link to the real archive),
#and your configuration file must be in /etc/logstash/ (with a .conf extension).
#And that’s it : you just have to run (as root) ./initStash.pl <configName>
#(For example if you configuration file is /etc/logstash/shipper.conf, then run ./initStash.pl shipper)

use strict;
$\=$/;

sub erreur($) {
my ($txt)=@_;
print $txt;
exit(1);
}

#
# Generates logstash start script
#
my $SCRIPT=$ARGV[0] || erreur ("Usage : initStash.pl <script-name>");

my $INITFILE = "/etc/init/logstash-$SCRIPT.conf";
my $CONFFILE = "/etc/logstash/$SCRIPT.conf";

erreur("$CONFFILE n'existe pas.") unless (-f $CONFFILE);
erreur("$INITFILE existe deja.") if (-f $INITFILE);
erreur("Logstash not installed.") unless (-f '/opt/logstash/logstash.jar');
qx!mkdir /var/log/logstash! unless (-d '/var/log/logstash');

open I, ">$INITFILE" || erreur("Cannot write to $INITFILE");
while (<DATA>) {
chomp;
s/__NAME__/$SCRIPT/g;
print I;
}
close I;
print "Done. You can start with : service logstash-$SCRIPT start";

__DATA__
description "logstash-__NAME__"
author "suresh.prajapati@inmobi.com"

start on startup
stop on shutdown

script
CONFIGNAME="__NAME__"
CONFIGFILE="$CONFIGNAME.conf"
# ^ Change these values with yours
cd /opt/logstash/
exec 2>&1
# Need to set LOGSTASH_HOME and HOME so sincedb will work
LOGSTASH_HOME="/opt/logstash"
LOGSTASH_ETC="/etc/logstash"
LOGSTASH_LOG="/var/log/logstash"
LOGSTASH_TMP="/tmp"
JAVA_XMS="128m"
JAVA_XMX="196m"
GC_OPTS="-XX:+UseParallelOldGC"
JAVA_OPTS="-server -Xms${JAVA_XMS} -Xmx${JAVA_XMX} -Djava.io.tmpdir=$LOGSTASH_TMP/"
LOGSTASH_OPTS="agent -f $LOGSTASH_ETC/$CONFIGFILE -v -l $LOGSTASH_LOG/logstash-$CONFIGNAME.log"
HOME=$LOGSTASH_HOME exec java $JAVA_OPTS $GC_OPTS -jar $LOGSTASH_HOME/logstash.jar $LOGSTASH_OPTS
end script

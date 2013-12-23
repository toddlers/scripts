#!/usr/bin/env perl
#Monitoring Services with Nagios::Plugin

use Nagios::Plugin;
use Time::HiRes qw( gettimeofday tv_interval );

my $np = Nagios::Plugin->new(
    shortname => 'template',
);

 # We need a URL to test
 my $url = "http://www.myapp.com/";

 # Create the UserAgent
 my $ua = LWP::UserAgent->new (
    agent       =>  'Nagios Application Check',
    from        =>  'webmaster@myapp.com',
 );

 # Change the timeout to 10 seconds instead of 3 min (180 seconds)
 $ua->timeout( 10 );

 my $t0 = [gettimeofday];
 my $resp = $ua->get($url);
 my $total = tv_interval($t0, [gettimeofday]);

 my $content;
 if (!$resp->is_success) {
    $np->nagios_exit( CRITICAL, $resp->status_line);
 }

 # Time it
 $np->add_perfdata(
    label => "time",
    value => "${total}",
    uom => "s",
 );

 $np->nagios_exit(OK, "Test Suggestions Found");

exit 0;

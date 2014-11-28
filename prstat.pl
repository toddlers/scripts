#!/usr/bin/env perl

my @ps = `/bin/ps aux`;
my @headers = split(/\s+/, shift(@ps));
my %users;

foreach (@ps) {
  chomp;
  my $col = 0;
  my %ps_entry;
  foreach (split(/\s+/, $_, $#headers + 1)) {
    $ps_entry{$headers[$col]} = $_;
    $col++;
  }

  next unless exists $ps_entry{USER};
  $users{$ps_entry{USER}} = { nproc=>0, size=>0, rss=>0, mem=>0, time=>0, cpu=>0 } unless exists $users{$ps_entry{USER}};
  my $user = $users{$ps_entry{USER}};

  $user->{nproc}++;
  $user->{size} += $ps_entry{VSZ} if exists $ps_entry{VSZ};
  $user->{rss} += $ps_entry{RSS} if exists $ps_entry{RSS};
  $user->{mem} += $ps_entry{'%MEM'} if exists $ps_entry{'%MEM'};
  $user->{cpu} += $ps_entry{'%CPU'} if exists $ps_entry{'%CPU'};
  $user->{time} += ($1 * 60) + $2 if (exists $ps_entry{'TIME'} && $ps_entry{'TIME'} =~ /^([0-9]+):([0-9]+)$/);
}

print "NPROC\tUSER\tSIZE\tRSS\tMEMORY\tTIME\tCPU\n";
foreach (sort { $users{$b}{cpu} <=> $users{$a}{cpu} } keys(%users)) {
  printf("%d\t%s\t%d\t%d\t%.1f\%\t%.2d:%.2d\t%.1f\%\n", $users{$_}{nproc}, $_, $users{$_}{size}, $users{$_}{rss}, $users{$_}{mem}, ($users{$_}{time} / 60), ($users{$_}{time} % 60), $users{$_}{cpu});
}

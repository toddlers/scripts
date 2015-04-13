#!/usr/bin/perl
my $dateregext = "(\d{4}-\d{2}-\d{2}-\d{2})";
my @matched = glob "*.rb";

#print "Returned list of file @file_list\n";
foreach my $f (@matched) {
my ($dtpart);
if ($f =~ /$dateregext/) {
$dtpart = $1;
print "$dtpart\n";
}
else {
    print "skipping $f because $f didn't match $dateregext\n";
}
}




# This routine takes in a path to an LDIF and returns an associative array of associative arrays.
sub getLDAPEntriesFromLDIF($)
{
    # Get the input path.
    my ($inputLDAPPath) = @_;
    
    # Open the input file.
    open(INPUTFILE, $inputLDAPPath);

    # Get all the lines of the INPUTFILE as an array.
    my @lines;
    @lines = <INPUTFILE>;
    
    # Create an associative array to hold all the ldap entries.
    my %ldapEntries;
    
    # Create a counter.
    my $i = 0;
    
    # Start a loop that iterates through all the lines of the LDIF.
    GET_NEXT_ENTRY: for($i; $i < @lines; $i++)
    {
        # Create an associative array to hold an individual entry.
        my %ldapEntry;
        
        # Start iterating through the entry.
        GET_NEXT_LINE: for($i; $i < @lines; $i++)
        {
            # If the line is a line continuation, restart the GET_NEXT_LINE loop on the next line.
            if($lines[$i] =~ /^(\s+\S+)\n$/)
            {
                $i++;
                redo GET_NEXT_LINE;
            }

            # If the line is empty, it signifies the end of an entry.  In this case, cache the individual entry
            # to ldapEntries and restart the GET_NEXT_ENTRY loop on the next line.
            if($lines[$i] =~ /^(\s*)\n$/)
            {
                if($ldapEntry{"dn"})
                {
                    $ldapEntries{$ldapEntry{"dn"}[0]} = [%ldapEntry];
                }
                else
                {
                    print "-------------- ERROR - No \"dn\" for entry at line ".($i+1)." --------------\n";
                }
                $i++;
                redo GET_NEXT_ENTRY;
            }
            
            # Create a temp variable to hold the current line of the ldif
            my $tempLine = $lines[$i];
            
            # Create a variable and populate it with the split line.
            my @splitLine;
            if($tempLine =~ /::/)
            {
                @splitLine = split(/::/, $tempLine);
                print "-------------- WARN value for attribute \"".$splitLine[0]."\" at line ".($i+1)." is encrypted --------------\n";
            }
            elsif($tempLine =~ /:/)
            {
                @splitLine = split(/:/, $tempLine);
            }
            
            # If the next line is a continuation of the current line, go ahead and add its text to the end of the current line.
            if($lines[$i+1] =~ /^(\s+)/)
            {
                if($lines[$i+1] =~ /(\S+)$/)
                {
                    chomp($splitLine[1]);
                    $splitLine[1] = $splitLine[1].$&;
                }
            }

            # Remove the line break from the end of the value portion of the line.
            chomp($splitLine[1]);
            
            # Remove the space from the front of the value portion of the line.
            $splitLine[1] = substr($splitLine[1], 1);
            
            # If everything went okay, go ahead and cache the attribute and value to the ldap entry.
            if($splitLine[1])
            {
                if(!($ldapEntry{$splitLine[0]}))
                {
                    my @tempArray = ($splitLine[1]);
                    $ldapEntry{$splitLine[0]} = [@tempArray];
                }
                else
                {
                    if($splitLine[1])
                    {
                        push(@{$ldapEntry{$splitLine[0]}}, $splitLine[1])
                    }
                    else
                    {
                        print "ERROR!";
                    }
                }
            }
        }
    }
    close(INPUTFILE);
    return %ldapEntries;
}

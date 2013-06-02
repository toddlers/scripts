#!/usr/bin/python2.6

import errno
import optparse
import os
import resource
import string
import sys


class GetOptions:

    def __init__(self):
        parser = optparse.OptionParser()
        parser.add_option("--conf", action="store",
                dest="conf", type="string", help="citrusleaf config")
        parser.add_option("--fd-limit", action="store", default=100000,
                dest="fdlimit", type="int",
                help="set fd limit [default: %default]")
        parser.add_option("--bin", action="store", default="/usr/bin/cld",
                dest="bin", type="string",
                help="Path to cld binary [default: %default]")
        parser.add_option("--pid-file", action="store",
                default="/var/run/cld.pid",
                dest="pidfile", type="string",
                help="path to cld pid file [default: %default]")
        self.parser = parser
        self.parse()

    def parse(self):
        (self.options, self.args) = self.parser.parse_args()

    def __getattr__(self, k):
        return getattr(self.options, k)

    def print_help_exit(self):
        self.parser.print_help()
        exit(1)


def main():
    try:
        opts = GetOptions()

        if not opts.conf:
            print "Need citrusleaf conf path"
            opts.print_help_exit()

        if os.geteuid() != 0:
            print "Need to be root to run this."
            sys.exit(1)

        cld = opts.bin
        conf = opts.conf
        cldargs = [cld, "--config-file", conf]
        fdlimit = opts.fdlimit
        pidfile = opts.pidfile

        # set fdlimit
        resource.setrlimit(resource.RLIMIT_NOFILE, (fdlimit, fdlimit))

        # if pid file exists and the process is alive abort
        if os.path.exists(pidfile):
            cldpid = int(open(pidfile).read())
            try:
                if os.kill(cldpid, 0) is None:
                    print "- cld already running pid %d. exiting." % (cldpid)
                    exit(1)
            except OSError, exp:
                if exp.errno in [errno.ESRCH]:
                    pass
                else:
                    raise OSError(exp)

        print "+ %s" % (string.join(cldargs, " "))
        sys.stdout.flush()
        sys.stderr.flush()
        os.chdir("/")
        os.execvp(cld, cldargs)
    except Exception, exp:
        raise Exception(exp)

if __name__ == "__main__":
        main()

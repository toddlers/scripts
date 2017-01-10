#!/usr/bin/env python
import psutil
import optparse
import subprocess
from time import time

class GetOptions:

    def __init__(self):
        parser = optparse.OptionParser()
        parser.add_option("--threshold", action="store", default=60,
                dest="threshold", type="int", help="alert threshold [default: %default]")
        self.parser = parser
        self.parse()

    def parse(self):
        (self.options, self.args) = self.parser.parse_args()

    def __getattr__(self, k):
        return getattr(self.options, k)

    def print_help_exit(self):
        self.parser.print_help()
        exit(1)


def run_command(cmd):
        p = subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True)
        out,err = p.communicate()
        return out.strip()

def main():
  try:
    opts = GetOptions()

    if not opts.threshold:
        print "Please specify threshold"
        opts.print_help_exit()
    th = opts.threshold
    pcron = run_command(["/usr/bin/pgrep  -P 1 cron"]).split("\n")[0]
    pp = psutil.Process(int(pcron))
    for p in pp.children():
     for pc in p.children():
         x = pc.children()[0]
         if (time() - x.create_time()) > th:
            print("%s => %d" % (x.name(),x.create_time()))
            exit(1)
  except Exception, exp:
       raise Exception(exp)

if __name__ == "__main__":
   main()

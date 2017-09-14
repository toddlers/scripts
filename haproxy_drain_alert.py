#!/usr/bin/env python
import subprocess
import smtplib
import socket
from email.mime.text import MIMEText
def main():
  currentstat = {}
  readstats = subprocess.check_output(["echo show stat | socat unix-connect:/var/run/haproxy.sock stdio"], shell=True)
  vips = readstats.splitlines()
  for i in range(0, len(vips)):
    if "DRAIN" in str(vips[i]):
       bk = vips[i].split(",")
       bk_lsc = (int(bk[23])/60)
       currentstat[bk[0]] ="{},{}".format(bk[1],str(bk_lsc))
  for k,v in currentstat.iteritems():
     box_name,since = v.split(",")
     alerts  = "DRAIN Set on  Backend : {} , Box Name : {} , Since : {} Min".format(str(k),str(box_name), str(since))
     mail(alerts)

def mail(alert):
    msg=MIMEText(alert)
    hname = socket.gethostname()
    me=hname+"@example.com"
    you="admin@example.com"
    msg["Subject"] = "Drain Alert on " + hname
    msg["From"] = me
    msg["To"] = you

    s = smtplib.SMTP("127.0.0.1")
    s.sendmail(me,[you],msg.as_string())
    s.quit

main()

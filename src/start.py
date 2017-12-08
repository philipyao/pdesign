#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import shlex, subprocess

def do_start():
    pname = "gamesvr"
    index = 5
    clusterid = 7000
    port = 4201
    shcmd = "./%s -i %d -p %d -c %d" % (pname, index, port, clusterid)
    #shcmd = "./%s -i %d -p %d -c %d >> %s.output 2>&1 &" % (pname, index, port, clusterid, pname)
    print shcmd

    args = shlex.split(shcmd)
    output = open("%s.output" % pname, 'a')
    p = subprocess.Popen(args, stdout=output, stderr=output)

    print "start server[%s], pid %d" % (pname, p.pid)

def main():
    do_start()

if __name__ == '__main__':
    main()

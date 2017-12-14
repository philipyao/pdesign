#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import os,time
import hqpy
from hqpy import logger
from hqpy import hqenv
from hqpy import hqio
from hqpy import hqthread
from hqpy import hqhosts

def TestLogger():

    logger.EXIT("xxxxxxxxxxxxxxxxxx")

    logger.WARN('test %s %s', 1, 222.2222)
    logger.WARN('xxxxxxxxxxxxxxxxxxxxxxxxx')
    time.sleep(1)
    logger.DEBUG('xxxxxxxxxxxxxxxxxxxxxxxxx222:%s,%s,%s', 11, 22, '33')
    logger.INFO('xxxxxxxxxxxxxxxxxxxxxxxxx222:%s,%s,%s', 11, 22, '33')
    logger.WARN('xxxxxxxxxxxxxxxxxxxxxxxxx222:%s,%s,%s', 11, 22, '33')
    logger.ERR('xxxxxxxxxxxxxxxxxxxxxxxxx222:%s,%s,%s', 11, 22, '33')

    logger.WARN('english 中文 XX YY 混杂')

    # import hqpy.console
    from hqpy import console
    console.red("sssssssssssssss")
    console.magenta("sssssssssssssss")
    console.cyan("sssssssssssssss")
    console.red("sssssssssssssss")

    import logging
    progress = console.ColoramaConsoleHandler()
    co = logging.getLogger('test')
    co.setLevel(logging.DEBUG) 
    ch_format = logging.Formatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-7s] %(message)s')
    progress.setFormatter(ch_format)
    progress.setTaskMode()
    co.addHandler(progress)


    co.info('test1')
    co.debug('test1ddd')
    co.info('test2')
    co.warn("wwwwwwwwwwwwwwwww")
    co.error('test1ddd') 

def TestJob():
    from hqpy import hqjob
    job_file = os.path.join(hqenv.get_var('HQVAR_FILES_DIR'), 'jobs/test_job.json')
    ret, job = hqjob.load_job(job_file, None)
    ret = job.do()
    print ret.string()


def TestSsh():
    # sau = hqpy.SshAuth("10.1.164.43", "user200", "u0Rf++Q", 22)
    # hret = hqpy.sftpPut(sau, "./hqnet.tgz", "/tmp/hqnet2.tgz")
    # print(hret.string())

    # hret = hqpy.sftpGet(sau, "/tmp/hqnet.tgz", "./hq.tgz")
    # print(hret.string())

    # hret, stdout, stderr  = hqpy.sshRun(sau, 'ifconfig && sleep 3 && exit 0')
    # print(hret.string())
    # print(stdout)
    # print(stderr)
    pass

def TestIO():
    flist = [ 'bin/platsvr', 'conf2', 'scripts', 'lua', 'gocext/cmake', 'make.py']
    hqio.tar_zip('/home/user00/server/trunk/golang', 'test3.tgz', flist)

def test_thread_job(p):
    logger.INFO('thread job:%s start', p)
    import random
    ts = random.randint(1,5)
    time.sleep(ts)
    logger.DEBUG('-----------thread job:%s,%s end', p, ts)
    return p*p

def TestMultiThread():
    res = hqthread.RunMultiIoJob(test_thread_job, [1,2,3,4,5,6,7,8,9,10], 10)
    print res
    res = hqthread.RunMultiIoJob(test_thread_job, [1,2,3,4,5,6,7,8,9,10])
    print res

def TestHosts():
    print(hqhosts.get_host_byid('id001').info() )
    hs = hqhosts.get_host_bygroup('platsvr')
    for x in hs:
        print(x.info())

def run():
    TestLogger()
    #TestJob()
    #TestIO()
    #TestMultiThread()

    #TestHosts()

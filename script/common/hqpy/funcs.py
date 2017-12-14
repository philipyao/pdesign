#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import signal,time,platform
import subprocess,json,ConfigParser
import os,sys,traceback
import hqpy
from hqpy import logger


def __signal_run_timeout(signum, frame):
    logger.ERR('Run Shell Timeout:%s', frame)
    raise Exception("Run Shell Timeout")

### stdout stderr 输出到PIPE 可能导致程序堵塞的问题
###  http://noops.me/?p=92  
###  http://backend.blog.163.com/blog/static/2022941262014016710912/
###  实时获取输出 http://blog.chinaunix.net/uid-26000296-id-4461555.html
def run_shell(cmd, timeout=None, show_console=False):

    hret = hqpy.HqError()
    pipe = None
    ostd = None
    oerr = None

    if timeout == None:
        timeout = 180

    ## paramiko 内timeout不生效
    if platform.system() != "Windows":
        signal.signal(signal.SIGALRM, __signal_run_timeout)
        signal.alarm(timeout) 

    try:
        pipe = subprocess.Popen(cmd, stdout=subprocess.PIPE, 
            stderr=subprocess.PIPE, shell=True, universal_newlines=True)
        # st = pipe.wait()
        # oerr = pipe.stderr.read()
        # ostd = pipe.stdout.read()

        stdlines = []
        while True:
            line = pipe.stdout.readline()
            if not line:
                break
            stdlines.append(line)
            if show_console == True:
                logger.LOG("%s", line.rstrip() )
        ostd = ''.join(stdlines)
        st = pipe.wait()
        oerr = pipe.stderr.read()
        strip_err = oerr.strip()
        if show_console == True and len(strip_err)>0:
            logger.LOG("%s", oerr)

        hret.errno = st
        if platform.system() != "Windows":
            signal.alarm(0)
    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_CMD_RUN_ERR
        hret.msg = 'FUNCS Error run_shell:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret, ostd, oerr

def run_shell_std(cmd, timeout=None):

    if timeout == None:
        timeout = 180

    hret = hqpy.HqError()

    ## paramiko 内timeout不生效
    signal.signal(signal.SIGALRM, __signal_run_timeout)
    signal.alarm(timeout) 

    try:
        hret.errno = subprocess.call(cmd, shell=True)
        signal.alarm(0)
    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_CMD_RUN_ERR
        hret.msg = 'FUNCS Error run_shell:{0}'.format(message)
        logger.ERR(traceback.format_exc())

    finally:
        return hret

def check_exit(ret, msg="CHECK_EXIT"):
    # if ret.iserr():
    #     logger.EXIT("%s failed, err:%s", msg, ret.string())

    if ret.iserr():
        errmsg = "%s failed, err:%s" % (msg, ret.string() )
        if logger.plogger is not None:
            logger.plogger.error(errmsg)
        else:
            print errmsg 
        sys.exit(254)  ## PYRET_ERR_EXIT

def exit(msg, *arg):
    #logger.EXIT(msg, *arg)
    if logger.plogger is not None:
        logger.plogger.error(msg, *arg)
        logger.plogger.debug("")
    else:
        print msg % (arg) 
    sys.exit(254)  ## PYRET_ERR_EXIT

def format(fmt, *args):
    return fmt % args

def loadjs(jsfile):
    hret = hqpy.HqError()
    jsobj = None
    try:
        fd = open(jsfile)
        jsobj = json.load(fd)
        fd.close()
    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)

        hret.errno = hqpy.PYRET_TASK_LOAD_ERR
        hret.msg = 'FUNCS Error loadjs:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret, jsobj

def writejs(jsfile, jsdata):
    hret = hqpy.HqError()
    try:
        fd = open(jsfile, 'w')
        #jsstr = json.dumps(jsdata, sort_keys=True,indent=4)
        json.dump(jsdata, fd, sort_keys=True, indent=4)
        fd.close()
    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)

        hret.errno = hqpy.PYRET_TASK_LOAD_ERR
        hret.msg = 'FUNCS Error writejs:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret

def loadini(inifile):
    hret = hqpy.HqError()
    iniobj = ConfigParser.ConfigParser()
    try:
        iniobj.read(inifile)
    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)

        hret.errno = hqpy.PYRET_TASK_LOAD_ERR
        hret.msg = 'FUNCS Error loadini:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret, iniobj


## 预设文件编码方式为 utf-8
## 对字符串进行编码
def str_encode(val):
    if type(val) == unicode:
        return val.encode("utf-8")
    else:
        return unicode(val, "utf-8").encode("utf-8")

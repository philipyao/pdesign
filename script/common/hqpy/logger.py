#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import os,inspect,logging,sys,time,tempfile
import logging.handlers
import console

plogger = None
predirect = None

class HqFormatter(logging.Formatter):
    def __init__( self, fmt=None, datefmt=None ):
        logging.Formatter.__init__(self, fmt, datefmt)
    def format( self, rec ):
        if rec.levelno == 99:
            out = rec.msg % rec.args
        else:
            stack = inspect.stack()
            rec.filename = os.path.basename(stack[9][1])
            rec.lineno = stack[9][2]
            out = logging.Formatter.format(self, rec)
        return out

class HqFormatter2(logging.Formatter):
    def __init__( self, fmt=None, datefmt=None ):
        logging.Formatter.__init__(self, fmt, datefmt)
    def format( self, rec ):
        if rec.levelno == 99:
            out = rec.msg % rec.args
        else:
            out = logging.Formatter.format(self, rec)
        return out

def new_logger(file_name, file_level, console_level = None):
    log_name = os.path.basename(file_name)
    logger = logging.getLogger(log_name)
    logger.setLevel(logging.DEBUG) #By default, logs all messages

    if console_level != None:
        #ch = logging.StreamHandler(sys.stdout) #StreamHandler logs to console
        ch = console.ColoramaConsoleHandler(sys.stdout)
        ch.setLevel(console_level)
        #ch_format = logging.Formatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-5s] %(message)s')
        ch_format = HqFormatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-5s] %(message)s')
        ch.setFormatter(ch_format)
        logger.addHandler(ch)

    #fh = logging.FileHandler(file_name)
    fh = logging.handlers.TimedRotatingFileHandler(file_name, 'H', 1, 9)
    fh.setLevel(file_level)
    fh_format = HqFormatter2('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-5s] %(message)s')
    #fh_format = logging.Formatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-5s] %(message)s')
    #fh_format = HqFormatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-7s] %(message)s')
    fh.setFormatter(fh_format)
    logger.addHandler(fh)
    logger.setLevel(logging.DEBUG)

    return logger

#本地日志
def init_normal_logger(log_dir, name='local', show_console=True, tm_tag=True):
    global plogger

    local_name = name
    if tm_tag == True:
        ISOTIMEFORMAT='%y%m%d'
        tmstr = int(time.strftime( ISOTIMEFORMAT, time.localtime() ))
        local_name = '%s.%s.log' % (name, tmstr)
    else:
        local_name = '%s.log' % (name)

    file_name = os.path.join(log_dir, local_name)
    console_level = None
    if show_console == True:
        console_level = logging.DEBUG

    plogger = new_logger(file_name, logging.DEBUG, console_level)

#远程任务日志，只文件，DEBUG级别日志
def init_job_logger(log_dir, taskid):
    global plogger
    task_name = 'job.%s.res' % (taskid)
    file_name = os.path.join(log_dir, task_name)

    logger = logging.getLogger(task_name)
    logger.setLevel(logging.DEBUG)  #By default, logs all messages

    ch = logging.StreamHandler() #StreamHandler logs to console
    ch.setLevel(logging.DEBUG)
    #ch_format = logging.Formatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-5s] %(message)s')
    ch_format = HqFormatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-5s] %(message)s')
    ch.setFormatter(ch_format)
    logger.addHandler(ch)

    fh = logging.FileHandler(file_name)
    fh.setLevel(logging.DEBUG)
    fh_format = logging.Formatter('%(asctime)s %(filename)s:%(lineno)d [%(levelname)-7s] %(message)s')
    fh.setFormatter(fh_format)
    logger.addHandler(fh)

    plogger = logger

## 
def DEBUG(fmt, *arg):
    if plogger is not None:
        plogger.debug(fmt, *arg)

def INFO(fmt, *arg):
    if plogger is not None:
        plogger.info(fmt, *arg)

def WARN(fmt, *arg):
    if plogger is not None:
        plogger.warn(fmt, *arg)

def ERR(fmt, *arg):
    if plogger is not None:
        plogger.error(fmt, *arg)

def EXIT(fmt, *arg):
    if plogger is not None:
        plogger.error(fmt, *arg)
    sys.exit(254)  ## PYRET_ERR_EXIT

def LOG(fmt, *arg):
    if plogger is not None:
        plogger.log(99, fmt, *arg)

### 控制台日志
################################################################

def PRED(fmt, *arg):
    console.red(fmt, *arg)

def PGREEN(fmt, *arg):
    console.green(fmt, *arg)

def PYELLOW(fmt, *arg):
    console.yellow(fmt, *arg)

def PBLUE(fmt, *arg):
    console.blue(fmt, *arg)

def PMAGENTA(fmt, *arg):
    console.magenta(fmt, *arg)

def PCYAN(fmt, *arg):
    console.cyan(fmt, *arg)

def PRINT(fmt, *arg):
    txt = fmt % (arg)
    print txt
    sys.stdout.flush()

## IO REDIRECT 
###################################################################

class IORedirectTmp:
    def __init__(self):
        self.tmpFile = tempfile.SpooledTemporaryFile(bufsize=5*1024*1024)
        self.__console_stdout__= sys.stdout
        self.__console_stderr__= sys.stderr

        sys.stdout = self.tmpFile
        sys.stderr = self.tmpFile
        
    def reset(self):
        sys.stdout=self.__console_stdout__
        sys.stderr=self.__console_stderr__

class IORedirectLogger:
    def __init__(self):
        self.__console_stdout__= sys.stdout
        self.__console_stderr__= sys.stderr

        sys.stdout = self
        sys.stderr = self
        
    def reset(self):
        sys.stdout=self.__console_stdout__
        sys.stderr=self.__console_stderr__

    def write(self, output_stream):
        LOG(output_stream.strip())

    def flush(self):
        pass

def RedirectTmp():
    global predirect
    if predirect == None:
        predirect = IORedirectTmp()

def RedirectLogger():
    global predirect
    if predirect == None:
        predirect = IORedirectLogger()

def GetRedirBuf():
    global predirect
    if predirect is not None:
        predirect.tmpFile.seek(0)
        return predirect.tmpFile.read()
    else:
        return ""


def ResetStd():
    global predirect
    if predirect is not None:
        predirect.reset()

def GetRedirFile():
    global predirect
    if predirect is not None:
        return predirect.tmpFile
    else:
        return None


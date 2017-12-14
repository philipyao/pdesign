#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import signal,platform
import paramiko,socket,traceback
import hqpy
from hqpy import logger

class SshAuth:

    """docstring for SshAuth"""
    def __init__(self, host, u, p, port=22):
        self.user = u
        self.pwd = p
        self.host = host
        self.port = port

    def string(self):
        return 'host:%s,u:%s,p:%s,port:%s' % (self.host, self.user, self.pwd, self.port)

class SshClient(paramiko.SSHClient): 
    def __init__( self, show_console=False):
        super(SshClient, self).__init__()
        self.show_console = show_console

    def call(self, command, timeout=None, bufsize=-1):  
        chan = self._transport.open_session() 
        chan.settimeout(timeout)
        chan.exec_command(command)  
        stdin = chan.makefile('wb', bufsize)  
        # stdout = chan.makefile('rb', bufsize).read()
        # stderr = chan.makefile_stderr('rb', bufsize).read()

        fstdout = chan.makefile('rb', bufsize)
        fstderr = chan.makefile_stderr('rb', bufsize)
        stdlines = []
        while True:
            line = fstdout.readline()
            if not line:
                break
            stdlines.append(line)
            if self.show_console == True:
                #print line.rstrip() 
                logger.LOG("%s", line.rstrip() )
        stdout = ''.join(stdlines)
        stderr = fstderr.read()

        status = chan.recv_exit_status()  
        return status, stdout, stderr  

def signal_put_timeout(signum, frame):
    logger.ERR('Sftp Put Timeout:%s', frame)
    raise Exception("Sftp Put Timed out")

def signal_get_timeout(signum, frame):
    logger.ERR('Sftp Get Timeout:%s', frame)
    raise Exception("Sftp Get Timed out") 
    
def sshRun(sau, cmd, timeout=None, show_console=False):
    hret = hqpy.HqError()
    sshc = hqpy.SshClient(show_console=show_console)

    sret = None
    stdin = None
    stdout = None
    stderr = None

    if timeout is None:
        timeout = 180

    try:
        # sshc.load_system_host_keys()  
        # sshc._policy = paramiko.AutoAddPolicy()  
        sshc.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        sshc.connect(hostname=sau.host, 
            port=sau.port,
            username=sau.user, 
            password=sau.pwd,
            pkey=None,
            key_filename=None)
        status, stdout, stderr = sshc.call(cmd, timeout=timeout)
        hret.errno = status
        signal.alarm(0)
    except socket.timeout:
        hret.errno = hqpy.PYRET_TIMEOUT
        hret.msg = 'SSH Error: timeout'
    except Exception, e:
        hret.errno = hqpy.PYRET_SSH_RUNERR
        hret.msg = 'SSH Error:{0}'.format(e)
    finally:
        sshc.close() 
        return hret, stdout, stderr

def sftpPut(sau, flocal, fremote, timeout = None):
    hret = hqpy.HqError()
    #ts = None
    sftp = None

    if timeout is None:
        timeout = 120

    ## paramiko 内timeout不生效
    if platform.system() != "Windows":
        signal.signal(signal.SIGALRM, signal_put_timeout)
        signal.alarm(timeout) 

    try:
        # ts=paramiko.Transport((sau.host, sau.port))
        # ts.connect(username=sau.user, password=sau.pwd, timeout=timeout)
        # sftp=paramiko.SFTPClient.from_transport(ts)

        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(hostname = sau.host, 
                    port = sau.port,
                    username=sau.user, 
                    password=sau.pwd,
                    timeout=timeout)
        sftp = ssh.open_sftp()
        #sftp.get_channel().settimeout(timeout)

        sftp.put(flocal, fremote)
        if platform.system() != "Windows":
            signal.alarm(0)
    except socket.timeout:
        hret.errno = hqpy.PYRET_TIMEOUT
        hret.msg = 'SFTP Put Error: timeout'
    except Exception, e:
        hret.errno = hqpy.PYRET_SSH_PUTERR
        hret.msg = 'SFTP Put Error:{0}'.format(e)

    finally:
        if sftp is not None:
            sftp.close()
        # if ts is not None:
        #     ts.close()  

        return hret
 
def sftpGet(sau, fremote, flocal, timeout = None):
    hret = hqpy.HqError()
    #ts = None
    sftp = None

    if timeout is None:
        timeout = 120

    ## paramiko 内timeout不生效
    signal.signal(signal.SIGALRM, signal_get_timeout)
    signal.alarm(timeout) 

    try:
        # ts=paramiko.Transport((sau.host, sau.port))
        # ts.connect(username=sau.user, password=sau.pwd, timeout=timeout)
        # sftp=paramiko.SFTPClient.from_transport(ts)

        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(hostname = sau.host, 
                    port = sau.port,
                    username=sau.user, 
                    password=sau.pwd, 
                    timeout=timeout)
        sftp = ssh.open_sftp()

        sftp.get(fremote, flocal)
        signal.alarm(0)
    except socket.timeout:
        hret.errno = hqpy.PYRET_TIMEOUT
        hret.msg = 'SFTP Get Error: timeout'
    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_SSH_GETERR
        hret.msg = 'SFTP Get Error:{0}'.format(ex)
        logger.ERR(traceback.format_exc())

    finally:
        if sftp is not None:
            sftp.close()
        # if ts is not None:
        #     ts.close()  
            
        return hret

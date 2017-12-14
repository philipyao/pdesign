#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import tarfile,re,string
import os,json,traceback
import hqpy
from hqpy import logger
from hqpy import hqenv
from hqpy import hqhosts
from hqpy import hqio

# 任务模块

class TaskBase(object):
    # idx = 0
    # name = ""
    # conf = None
    # job = None

    """docstring for TaskBase"""
    def __init__(self, arg, job):
        self.name = arg['name']
        self.idx = arg['idx']
        self.conf = arg
        self.job = job

    def info(self):
        sau = self.get_param().get_auth()
        return 'task[%s:%s:%s]' % (sau.host, self.idx, self.name)

    def get_param(self):
        return self.job.jobp

    def check(self):
        return True, ''

    def get_arg(self, name):
        if self.conf.has_key(name) == False:
            return False, None

        val = self.conf[name]
        # if type(val) == unicode:
        #     val = val.encode("utf-8")
        btype = type(val) == str or type(val) == unicode

        if btype and val.startswith('JVAR_') == True:
            m = re.match('JVAR_(\d+)', val)
            if m is None:
                logger.ERR('%s invalid job var:%s', self.info(), val)
                return False, None
            aidx = int(m.group(1))
            args = self.get_param().get_args()
            if aidx >= len(args):
                logger.ERR('%s invalid job var idx:%s', self.info(), val)
                return False, None

            logger.DEBUG('replace arg:%s,%s', aidx, args[aidx])
            return True, args[aidx]
        else:
            return True, val

    def do(self):
        logger.DEBUG('hq default task do')
        return hqpy.HqError()

## put & run & get
## 功能说明
## 上传本地hat/script 目录下的到目标机运行，并可取回指定目录结果文件
## 脚本包括 python shell lua
## 上传脚本在目标机 /tmp/hat 目录，需要返回的结果文件也需要输出在该目录
## 返回结果在 hat/var/job/tmp 
## 参数说明
class TaskPRG(TaskBase):
    #url = ""
    """docstring for TaskPRG"""
    def __init__(self, arg, job):
        super(TaskPRG, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('script') == False:
            return False, 'TaskPRG need [script] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()
        script_name = self.conf['script']
        bname = os.path.basename(script_name)
        # script_name 相对于hat/script 路径名
        local_file = os.path.join(hqenv.get_var('HQVAR_SCRIPT_DIR'), script_name)
        # 远端置于 /tmp/hatjobs/
        remote_path = '/tmp/hat'
        sau = self.get_param().get_auth()
        jobid = self.get_param().get_jobid()

        remote_file = os.path.join(remote_path, '%s_%s' %(jobid, bname) )

        shbin = "sh"
        if bname.endswith('.py'):
            shbin = 'python'
        elif bname.endswith('.lua'):
            shbin = 'lua'

        tmo = None
        if self.conf.has_key('timeout'):
            tmo = self.conf['timeout']

        #先确保远端目录建立
        ret, ostd, oerr = hqpy.sshRun(sau, 'mkdir -p %s' % (remote_path), 15)
        if ret.iserr():
            logger.ERR('TaskPRG dir prepare failed:%s, std:%s, err:%s', ret.string(), ostd, oerr)
            return ret

        # 推脚本
        ret = hqpy.sftpPut(sau, local_file, remote_file, 30)
        if ret.iserr():
            logger.ERR('TaskPRG put script failed:%s, std:%s, err:%s', ret.string(), ostd, oerr)
            return ret

        # 运行脚本
        ret, ostd, oerr = hqpy.sshRun(sau, '%s %s' % (shbin, remote_file), tmo)
        if ret.iserr():
            logger.ERR('TaskPRG dir prepare failed:%s, std:%s, err:%s', ret.string(), ostd, oerr)
            return ret

        # 获取结果
        if self.conf.has_key('results'):
            ress = self.conf['results']
            local_tmp = hqenv.get_var('HQVAR_VAR_DIR')
            for r in ress:
                r_file = os.path.join(remote_path,  r)
                l_file = os.path.join(local_tmp, '%s_%s_%s' % (jobid, sau.host, r) )
                ret = hqpy.sftpGet(sau, r_file, l_file)
                if ret.iserr():
                    logger.ERR('TaskPRG get result %s failed:%s', r, ret.string())
                    return ret
                else:
                    logger.INFO('TaskPRG result get ok:%s', l_file)

        ## 任务结束
        return ret

## get
## 功能说明
## 从目标机取文件，本地文件在 var/job/tmp 目录
class TaskGet(TaskBase):
    #url = ""
    """docstring for TaskGet"""
    def __init__(self, arg, job):
        super(TaskGet, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('files') == False:
            return False, 'TaskGet need [files] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()
        r_files = self.conf['files']
        sau = self.get_param().get_auth()
        jobid = self.get_param().get_jobid()

        tmo = None
        if self.conf.has_key('timeout'):
            tmo = self.conf['timeout']

        local_tmp = hqenv.get_var('HQVAR_VAR_DIR')
        for r in r_files:
            l_file = os.path.join(local_tmp, '%s_%s_%s' % (jobid, sau.host, r) )
            ret = hqpy.sftpGet(sau, r, l_file, tmo)
            if ret.iserr():
                logger.ERR('TaskGet get file %s failed:%s', r, ret.string())
                return ret
            else:
                logger.INFO('TaskGet ok:%s', l_file)

        return ret

## put
## 功能说明
## 本机上传文件到 目标机指定目录
## 本地文件配置 files 指相对于 hat目录的相对路径
## 目标机目录配置 path 相对于 /hgame 的相对目录
class TaskPut(TaskBase):
    #url = ""
    """docstring for TaskPut"""
    def __init__(self, arg, job):
        super(TaskPut, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('files') == False:
            return False, 'TaskPut need [files] arg'

        if self.conf.has_key('path') == False:
            return False, 'TaskPut need [path] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()
        l_files = self.conf['files']
        sau = self.get_param().get_auth()
        jobid = self.get_param().get_jobid()

        tmo = None
        if self.conf.has_key('timeout'):
            tmo = self.conf['timeout']

        for x in self.conf['files']:
            local_file = os.path.join(hqenv.get_var('HQVAR_WORK_DIR'), x)
            bname = os.path.basename(x)
            remote_file = os.path.join(self.conf['path'], bname)
            ret = hqpy.sftpPut(sau, local_file, remote_file, tmo)
            if ret.iserr():
                logger.ERR('%s put file failed:%s', self.info(), x)
                return ret

        return ret

## sshcmd
## 功能说明
## 远程执行 cmd 配置的脚本
class TaskSshCmd(TaskBase):
    #url = ""
    """docstring for TaskSshCmd"""
    def __init__(self, arg, job):
        super(TaskSshCmd, self).__init__(arg, job)
        self.url = ""

    def check(self):
        # if self.conf.has_key('cmd') == False:
        #     return False, 'TaskSshCmd need [cmd] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()
        tmo = None
        if self.conf.has_key('timeout'):
            tmo = self.conf['timeout']

        sau = self.get_param().get_auth()
        jobid = self.get_param().get_jobid()

        if self.conf.has_key('cmd') == True:
            shcmd = self.conf['cmd']
            logger.DEBUG('%s run ssh cmd: %s', self.info(), shcmd)
            ret, ostd, oerr = hqpy.sshRun(sau, shcmd, tmo)
            logger.INFO('%s ssh cmd result: %s', self.info(), ret.string() )
            if ret.iserr():
                logger.INFO('========================================================std:\n%s', ostd)
                logger.ERR( '========================================================err:\n%s', oerr)  
                return ret
            else:
                logger.INFO('========================================================std:\n%s', ostd)
                logger.INFO('========================================================err:\n%s', oerr)  
                

        if self.conf.has_key('script') == True: 
            script_file = self.conf['script']
            if os.path.isabs(script_file) == False:
                bhas, val_path = self.get_arg('path')
                if bhas : 
                    script_file = os.path.join(val_path, script_file)
                else:
                    script_file = os.path.join(hqenv.get_var('HQVAR_SCRIPT_DIR'), script_file)

            fname = string.split(script_file)[0]
            bname = os.path.basename(fname)
            shbin = "sh"
            if bname.endswith('.py'):
                shbin = 'python'
            elif bname.endswith('.lua'):
                shbin = 'lua'

            shcmd = '%s %s' % (shbin, script_file)

            ret, ostd, oerr = hqpy.sshRun(sau, shcmd, tmo)
            if ret.iserr():
                logger.ERR('%s run ssh script error: %s', self.info(), ret.code())
                logger.INFO('========================================================std:\n%s', ostd)
                logger.ERR( '========================================================err:\n%s', oerr)  
                return ret
            else:
                logger.INFO('%s run ssh script ok: %s', self.info(), ret.code())
                logger.INFO('========================================================std:\n%s', ostd)
                logger.INFO('========================================================err:\n%s', oerr)  

        return ret

## sget TODO
class TaskSGet(TaskBase):
    #url = ""
    """docstring for TaskSGet"""
    def __init__(self, arg, job):
        super(TaskSGet, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('files') == False:
            return False, 'TaskSGet need [files] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()
        ms = hqhosts.get_host_byid('master')
        if ms is None:
            ret.errno = hqpy.PYRET_ERR
            ret.ms = 'master svr not found'
            return ret

        sau = ms.get()

        # pkg_file = self.get_param('file')

        # local_file = os.path.join(hqenv.get_var('HQVAR_VAR_DIR'), x)
        # remote_file = os.path.join(filep, x)
        # logger.DEBUG('%s get file:%s from:%s to:%s' % (self.info(), x, sau.ip, local_file))
        # ret = hqpy.sftpGet(sau, remote_file, local_file)
        # if ret.iserr():
        #     logger.DEBUG('%s get file failed:%s', self.info(), ret.string())
        #     return ret

        return ret

## sput TODO
class TaskSPut(TaskBase):
    #url = ""
    """docstring for TaskSPut"""
    def __init__(self, arg, job):
        super(TaskSPut, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('files') == False:
            return False, 'TaskSPut need [files] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()

        # ms = hqhosts.get_host_byid(msid)
        # if ms is None:
        #     ret.errno = hqpy.PYRET_ERR
        #     ret.ms = 'master svr not found'
        #     return ret

        # sau = ms.get()

        # local_file = os.path.join(hqenv.get_var('HQVAR_VAR_DIR'), x)
        # remote_file = os.path.join(filep, x)
        # logger.DEBUG('%s get file:%s from:%s to:%s' % (self.info(), x, s['host'], local_file))
        # ret = hqpy.sftpGet(sau, remote_file, local_file)
        # if ret.iserr():
        #     logger.DEBUG('%s get file failed:%s', self.info(), ret.string())
        #     return ret

        return ret


## cmd
## 功能说明
## 执行本地脚本 
## cmd 可配置一句脚本，script 执行的脚本文件，相对于script目录
class TaskCmd(TaskBase):
    #url = ""
    """docstring for TaskCmd"""
    def __init__(self, arg, job):
        super(TaskCmd, self).__init__(arg, job)
        self.url = ""

    def check(self):

        return True, ''

    def do(self):
        tmo = None
        if self.conf.has_key('timeout'):
            tmo = self.conf['timeout']

        ret = hqpy.HqError()
        if self.conf.has_key('cmd') == True:
            shcmd = self.conf['cmd']
            logger.DEBUG('%s run cmd: %s', self.info(), shcmd)
            ret, ostd, oerr = hqpy.run_shell(shcmd, tmo)
            if ret.iserr():
                logger.ERR('%s run shell error: %s', self.info(), ret.code())
                logger.INFO('========================================================std:\n%s', ostd)
                logger.ERR( '========================================================err:\n%s', oerr)  
                return ret
            else:
                logger.INFO('%s run shell ok: %s', self.info(), ret.code())
                logger.INFO('========================================================std:\n%s', ostd)
                logger.ERR( '========================================================err:\n%s', oerr)  
            

        if self.conf.has_key('script') == True: 
            script_file = self.conf['script']
            if os.path.isabs(script_file) == False:
                bhas, val_path = self.get_arg('path')
                if bhas : 
                    script_file = os.path.join(val_path, script_file)
                else:
                    script_file = os.path.join(hqenv.get_var('HQVAR_SCRIPT_DIR'), script_file)

            bname = os.path.basename(script_file)
            shbin = "sh"
            if bname.endswith('.py'):
                shbin = 'python'
            elif bname.endswith('.lua'):
                shbin = 'lua'

            shcmd = '%s %s' % (shbin, script_file)
            ret, ostd, oerr = hqpy.run_shell(shcmd, tmo)
            if ret.iserr():
                logger.ERR('%s run script error: %s', self.info(), ret.code())
                logger.INFO('========================================================std:\n%s', ostd)
                logger.ERR( '========================================================err:\n%s', oerr)  
                return ret
            else:
                logger.INFO('%s run script ok: %s', self.info(), ret.code())
                logger.INFO('========================================================std:\n%s', ostd)
                logger.ERR( '========================================================err:\n%s', oerr)  

        return ret

## yum
class TaskYum(TaskBase):
    #url = ""
    """docstring for TaskYum"""
    def __init__(self, arg, job):
        super(TaskYum, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('packages') == False:
            return False, 'TaskYum need [packages] arg'

        return True, ''

    def do(self):
        packs = self.conf['packages']
        logger.DEBUG('%s run yum to install packages:%s', self.info(), packs)
        cmdstr = 'yum install -y {0}'.format(packs)
        ret, ostd, oerr = hqpy.run_shell(cmdstr)
        if ret.iserr():
            logger.ERR("%s run yum cmd error:%s", self.info(), ret.string())

        return ret

## update
class TaskUpdate(TaskBase):
    #url = ""
    """docstring for TaskUpdate"""
    def __init__(self, arg, job):
        super(TaskUpdate, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('version') == False:
            return False, 'TaskUpdate need [version] arg'
        if self.conf.has_key('path') == False:
            return False, 'TaskUpdate need [path] arg'
        if self.conf.has_key('set') == False:
            return False, 'TaskUpdate need [set] arg'

        return True, ''

    def do(self):
        ret = hqpy.HqError()
        ms = hqhosts.get_host_byid('master')
        if ms is None:
            ret.errno = hqpy.PYRET_ERR
            ret.ms = 'master svr not found'
            return ret

        sau = ms.get()

        bhas, hat_path = self.get_arg('path')
        if bhas == False:
            ret.errno = hqpy.PYRET_ERR
            ret.ms = 'path arg not found'
            return ret

        bhas, dst_set = self.get_arg('set')
        if bhas == False:
            ret.errno = hqpy.PYRET_ERR
            ret.ms = 'set arg not found'
            return ret

        bhas, pkg_version = self.get_arg('version')
        if bhas == False:
            ret.errno = hqpy.PYRET_ERR
            ret.ms = 'version arg not found'
            return ret
        
        pkg_name = "hgame.svr.%s.tgz" % (pkg_version)
        local_file = os.path.join(hqenv.get_var('HQVAR_VAR_DIR'), pkg_name)
        remote_file = os.path.join("/hgame/hat/files/pkg_server", pkg_name)
        print remote_file
        logger.DEBUG('%s get file:%s from:%s to:%s' % (self.info(), pkg_name, sau.host, local_file))
        ret = hqpy.sftpGet(sau, remote_file, local_file)
        if ret.iserr():
            #logger.ERR('%s get file failed:%s', self.info(), ret.string())
            return ret

        ## TODO 
        dst_dir = os.path.join(hat_path, dst_set)
        hqio.mkdirs(dst_dir)

        tf = tarfile.open(local_file)
        logger.DEBUG('%s update all files to %s', self.info(), dst_dir)
        tf.extractall(path=dst_dir)

        # file list
        # else:
        #     names = tf.getnames()
        #     paths = []
        #     fnames = []
        #     for f in flist:
        #         if f.endswith("*") == True:
        #             paths.append(f[:(len(f)-1)])
            
        #     for name in names:
        #         if name in flist:
        #             fnames.append(name)
        #         else:
        #             for p in paths:
        #                 if name.startswith(p) == True:
        #                     fnames.append(name)
        #     logger.DEBUG('%s update files to %s:\n%s', self.info(), dpath, fnames)
        #     for f in fnames:
        #         tf.extract(f, path=dpath)

        return hqpy.HqError()


#cron
class TaskCron(TaskBase):
    #url = ""
    """docstring for TaskCron"""
    def __init__(self, arg, job):
        super(TaskCron, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('shell') == False:
            return False, 'TaskCron need [shell] arg'

        return True, ''

    def do(self):
        return ret

#bakup
class TaskBakup(TaskBase):
    #url = ""
    """docstring for TaskBakup"""
    def __init__(self, arg, job):
        super(TaskBakup, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('shell') == False:
            return False, 'TaskBakup need [shell] arg'

        return True, ''

    def do(self):
        return ret

#mysql
class TaskMysql(TaskBase):
    #url = ""
    """docstring for TaskMysql"""
    def __init__(self, arg, job):
        super(TaskMysql, self).__init__(arg, job)
        self.url = ""

    def check(self):
        if self.conf.has_key('shell') == False:
            return False, 'TaskMysql need [shell] arg'

        return True, ''

    def do(self):
        return ret

# Job boot 参数
class JobBootParam(object):
    # islocal = False
    # auth = None
    # jid = 0
    # jname = ""
    # args = None  # 命令行参数
    # in_file = None   # 输入文件
    # out_log = None   # 结果日志
    # out_file = None  # 结果文件

    """docstring for JobBootParam"""
    def __init__(self, auth, jid, jname, args, in_file=None, out_file=None):
        super(JobBootParam, self).__init__()
        self.islocal = False
        self.auth = auth
        self.jid = jid
        self.jname = jname
        self.args = args
        self.in_file = in_file
        self.out_file = out_file
        self.out_log = None

    def get_auth(self):
        ## TODO root/ user
        return self.auth.user

    def get_jobid(self):
        return self.jid

    def get_args(self):
        return self.args

## Job 对象
class JobEntry(object):
    # jname = ""
    # tlist = []
    # ptasks = None
    # jobp = None

    """docstring for JobEntry"""
    def __init__(self, jp):
        super(JobEntry, self).__init__()

        self.jobp = jp
        self.jname = ""
        self.ptasks = None
        self.tlist = []

    def load(self, arg):
        ret = hqpy.HqError()

        self.jname = arg['jname']
        self.ptasks = arg['tasks']

        self.jargn = 0
        if arg.has_key('jargn'):
            self.jargn = arg['jargn']

        if len(self.jobp.get_args()) < self.jargn:
            ret.errno = hqpy.PYRET_TASK_CHECK_ERR
            ret.msg = '%s need %s args' % (self.jname, self.jargn)
            return ret

        idxs = 1
        for x in self.ptasks:
            x['idx'] = idxs
            name = x['name']
            ts = None
            if name == 'sget':
                ts = TaskSGet(x, self)
            elif name == 'sput':
                ts = TaskSPut(x, self)
            elif name == 'yum':
                ts = TaskYum(x, self)
            elif name == 'cmd':
                ts = TaskCmd(x, self)
            elif name == 'update':
                ts = TaskUpdate(x, self)
            elif name == 'prg':
                ts = TaskPRG(x, self)
            elif name == 'put':
                ts = TaskPut(x, self)
            elif name == 'get':
                ts = TaskGet(x, self)
            elif name == 'sshcmd':
                ts = TaskSshCmd(x, self)
            else:
                ts = TaskBase(x, self)
            
            bret, sres = ts.check()
            if bret == False:
                logger.DEBUG('%s check failed:%s', ts.info(), sres)
                ret.errno = hqpy.PYRET_TASK_CHECK_ERR
                ret.msg = sres
                return ret

            self.tlist.append(ts)
            idxs = idxs + 1
            logger.DEBUG('%s load sucesss', ts.info())

        return ret

    def do(self):
        hqret = hqpy.HqError()
        sau = self.jobp.get_auth()
        host_info = '[%s@%s:%s]' % (sau.user, sau.host, self.jobp.get_jobid()) 
        logger.INFO('>>>>>>>>>>>>>>>>>> %s do all tasks in [%s] ', host_info, self.jname)
        for t in self.tlist:
            logger.INFO('%s start to do', t.info())
            hqret = t.do()
            if hqret.iserr():
                logger.ERR('%s do failed, err:%s', t.info(), hqret.string())
                return hqret
            else:
                logger.INFO('%s do sucesss', t.info())
        logger.INFO('<<<<<<<<<<<<<<<<< %s do all tasks in [%s] ok', host_info, self.jname)

        return hqret    

def load_job(job_file, param):
    fd = None
    hret = hqpy.HqError()
    je = None
    try:
        fd = open(job_file)
        jscfg = json.load(fd)

        je = JobEntry(param)
        hret = je.load(jscfg)
        if hret.iserr():
            je = None
            logger.ERR('load job failed:%s', hret.string())

    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_TASK_LOAD_ERR
        hret.msg = 'TASK Error:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret, je
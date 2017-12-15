#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import sys, os, platform, time, signal, traceback, string, subprocess, psutil

var_work_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir))
sys.path.append(os.path.join(var_work_dir, 'script', 'common'))
os.chdir(var_work_dir)

# import hqpy
try:
    import hqpy
    from hqpy import logger
except Exception, e:
    print 'HQPY IMPORT ERROR %s' % (e)
    sys.exit(255) ## PYRET_NEED_INIT
#############################################################################################


## 脚本变量及配置
##############################################################################################
var_work_dir        = var_work_dir.decode('utf-8')
var_bin_dir         = os.path.join(var_work_dir, "bin")
var_log_dir         = os.path.join(var_work_dir, "log")
var_pid_dir         = os.path.join(var_work_dir, "pid")
var_cluster_conf    = None

var_usage = """service.py [options] args
使用说明
示例：
./service.py -a           编译所有项目执行： 编译前预处理、清除编译结果、编译 
./service.py -i           """
var_usage = hqpy.str_encode(var_usage)
var_usage = "Usage"

################################################################################################

## 执行shell命令
def run_shell(cmd, tmo=30):
    hret, ostd, oerr = hqpy.run_shell(cmd, tmo)
    if hret.iserr():
        logger.ERR('bash:%s\n%s\n%s', cmd, ostd, oerr)
    else:
        # logger.DEBUG(ostd)
        # logger.DEBUG(oerr)
        #print hret.string()
        pass
    return hret

def gen_process_options(name):
    pgname = "none"
    idxs = []

    fields = name.split(":")
    pgname = fields[0]
    if len(fields) == 1:
        idxs.append(0)
    else:
        beg = 1 # 进程起始序号, 默认1
        if len(fields) > 2:
            beg = int(fields[2])
            
        num = int(fields[1])  # 启动多少个实例
        for i in range(beg, beg+num):
            idxs.append(i)
    
    return pgname, idxs

def get_pid_byname(name):
    global var_pid_dir

    pid = 0
    hret = hqpy.HqError()

    fpid = os.path.join(var_pid_dir, "run.%s.pid" % (name))
    if os.path.exists(fpid) == False:
        return hret, pid

    fp = None
    try:
        fp = open(fpid, "r")
        pidstr = fp.readline()
        pidstr = pidstr.strip()
        pid = int(pidstr)
    except Exception, ex:
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_SERVICE_ERR
        hret.msg = 'Service Error:{0}'.format(ex)
    finally:
        if fp is not None:
            fp.close()
        return hret, pid

def remove_pid_file(name):
    global var_pid_dir

    fpid = os.path.join(var_pid_dir, "run.%s.pid" % (name))
    if os.path.exists(fpid) == False:
        return

    os.remove(fpid)

def read_stderr_log(fcore):
    fd = open(fcore)
    if fd is None:
        return ""
    msgs = []
    while True:
        line = fd.readline()
        if not line:
            break
        if line.find("APP_START_DONE_TAG") >= 0 :
            msgs = []
        else:
            msgs.append(line)

    fd.close()
    msg = ''.join(msgs)
    return msg


def run_start_cmd(pname, cmd, tmo=18):
    import tempfile
    out_temp = tempfile.SpooledTemporaryFile(bufsize=30*1000)
    proc = subprocess.Popen(cmd, stdout=out_temp, stderr=out_temp, shell=True, universal_newlines=True)
    pid = proc.pid

    sleeps = 0
    while sleeps < tmo:
        #if sleeps > 3 and psutil.pid_exists(pid) == False:
        #    #进程挂了
        #    out_temp.close()
        #    logger.DEBUG('service[%s] breakdown', pname)

        #    core_log_file = './log/core.{0}.log'.format(pname)
        #    msg = read_stderr_log(core_log_file)
        #    return hqpy.HqError(hqpy.PYRET_SERVICE_ERR, "service: %s breakdown:\n%s" % (pname, msg) ), pid

        time.sleep(1)
        sleeps = sleeps + 1
        sret = proc.poll()

        if sret != 0 :
            out_temp.seek(0)
            msg = out_temp.read()
            out_temp.close()
            return hqpy.HqError(hqpy.PYRET_SERVICE_ERR, "run shell err:%s,%s" % (sret, msg)), pid

        ## 检查到pid文件存在
        hret, rpid = get_pid_byname(pname)
        if hret.iserr():
            out_temp.close()
            return hret, 0

        if rpid > 0:
            out_temp.close()
            return hqpy.OK(), pid
        else:
            if sleeps % 2 == 0:
                logger.DEBUG('service[%s] is starting', pname)
        
    out_temp.seek(0)
    msg = out_temp.read()
    out_temp.close()
    return hqpy.HqError(hqpy.PYRET_SERVICE_ERR, "timeout: %s" % (msg) ), pid

def fun_start_old(options, name):
    global var_work_dir
    os.chdir(var_work_dir)

    pname = name
    if options.pidx > 0:
        pname = '{0}{1}'.format(name, options.pidx)
    hret, pid = get_pid_byname(pname)
    hqpy.check_exit(hret, "get service pid")
    if psutil.pid_exists(pid):
        logger.WARN('service[%s:%s] is running', pname, pid)
    else:
        if pid > 0:
            ## pid file 存在，进程不在了
            logger.WARN('service[%s:%s] not running, but pid file exist, remove it.', pname, pid)
            remove_pid_file(pname)

        ## 启动程序
        cons = " "
        if options.console:
            cons = "-c"
        libexport = "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$HOME/hgame/lib:./lib"
        #os.system(libexport)
        spidx = ""
        if options.pidx > 0:
            spidx = "-p %s" % (options.pidx)

        core_log_file = './log/core.{0}.log'.format(name)
        shcmd = "%s ; ./bin/%s %s %s %s >> %s 2>&1 &" % (libexport, name, cons, spidx, options.xset, core_log_file)

        ## TODO 增强健壮性
        if os.path.exists('./bin/%s' % (name)) == False:
            hqpy.exit("service[%s] ./bin/%s not exist" % (pname, name))

        hret, rpid = run_start_cmd(pname, shcmd, options.tmo)
        hqpy.check_exit(hret, "service[%s:%s] run" % (pname, rpid))
        logger.DEBUG('service[%s:%s] start ok', pname, rpid)

def fun_start(name, clusterid, index, port):
    logger.DEBUG("%s %d %d %d", name, clusterid, index, port)
    global var_work_dir
    os.chdir(var_work_dir)

    pname = name
    if index > 0:
        pname = '{0}{1}'.format(name, index)
    hret, pid = get_pid_byname(pname)
    hqpy.check_exit(hret, "get service pid")
    if psutil.pid_exists(pid):
        logger.WARN('service[%s:%s] is running', pname, pid)
    else:
        if pid > 0:
            ## pid file 存在，进程不在了
            logger.WARN('service[%s:%s] not running, but pid file exist, remove it.', pname, pid)
            remove_pid_file(pname)

        ## 启动程序
        stridx = ""
        if index > 0:
            stridx = "-i %s" % (index)
        strport = "-p %s" % (port)
        strcluster = "-c %s" % (clusterid)
        output_file = './log/output.{0}.log'.format(name)
        shcmd = "./bin/%s %s %s %s >> %s 2>&1 &" % (name, strcluster, stridx, strport, output_file)
        logger.DEBUG("shcmd: %s", shcmd)
        ## TODO 增强健壮性
        if os.path.exists('./bin/%s' % (name)) == False:
            hqpy.exit("service[%s] ./bin/%s not exist" % (pname, name))

        hret, rpid = run_start_cmd(pname, shcmd, options.tmo)
        hqpy.check_exit(hret, "service[%s:%s] run" % (pname, rpid))
        logger.DEBUG('service[%s:%s] start ok', pname, rpid)

def fun_stop(options, name, kill=False):
    global var_work_dir
    os.chdir(var_work_dir)

    pname = name
    itmo = options.tmo
    itmo = 60 # 部分服停服时间较长
    if options.pidx > 0:
        pname = '{0}{1}'.format(name, options.pidx)
    hret, pid = get_pid_byname(pname)
    hqpy.check_exit(hret, "get service pid")

    if psutil.pid_exists(pid) == False:
        logger.WARN('service[%s:%s] is not running', pname, pid)
        return hret
    else:
        p = psutil.Process(pid)
        if kill == True:
            p.kill()
        else:
            p.terminate()

        sleeps = 0
        while sleeps < itmo:
            time.sleep(1)
            sleeps = sleeps + 1
            if p.is_running() == False:
                logger.DEBUG('service[%s:%s] stoped', pname, pid)
                return hret
            else:
                if sleeps % 2 == 0:
                    logger.DEBUG('service[%s:%s] is stopping', pname, pid)

        logger.EXIT('service[%s:%s] stop failed', pname, pid) 
        #return hqpy.HqError(hqpy.PYRET_SERVICE_ERR, "service [%s:%s] stop failed" % (pname, pid))

def fun_signal(options, name):
    global var_work_dir
    os.chdir(var_work_dir)

    pname = name
    sig = options.sigval
    if sig == None:
        sig = signal.SIGHUP
    else:
        if sig == "u1":
            sig = signal.SIGUSR1
        elif sig == "u2":
            sig = signal.SIGUSR2
        else:
            hqpy.exit("invalid signal:%s", sig)

    if options.pidx > 0:
        pname = '{0}{1}'.format(name, options.pidx)
    hret, pid = get_pid_byname(pname)
    hqpy.check_exit(hret, "get service pid")

    if psutil.pid_exists(pid) == False:
        logger.EXIT('service[%s:%s] is not running', pname, pid) 
    else:
        p = psutil.Process(pid)
        p.send_signal(sig)
        logger.DEBUG('service[%s:%s] signal:%s sended', pname, pid, sig)

def run_service_config(options, names):
    global var_work_dir
    logger.DEBUG('run service config cmd:%s', names)

    idx = options.idx 
    zoneid = options.zone
    groupid = options.group

    fdev = ""
    if options.dev:
        fdev = '-d'

    var_script_dir = os.path.join(var_work_dir, 'scripts/01_gen_cfg.lua')
    shcmd = "lua %s %s -p %s -c mkzone -i %s" % (var_script_dir, fdev, var_work_dir, idx)
    #print shcmd
    hret = hqpy.run_shell_std(shcmd)
    if hret.iserr():
        return hret

    args = ' '.join(names)
    shcmd = "lua %s %s -p %s -i %s -g %s -z %s %s" % (var_script_dir, fdev, var_work_dir, idx, groupid, zoneid, args)
    #print shcmd
    hret = hqpy.run_shell_std(shcmd)
    if hret.iserr():
        return hret  
    logger.DEBUG('run service config done')

def run_service_start(options, servers):
    logger.DEBUG('run service start cmd')
    for n in servers:
        logger.DEBUG("start: %s(@%d) %d-%d %d", n['server'], n['clusterid'], n['startidx'], n['endidx'], n['portbase'])
        if n['endidx'] == 0:
            fun_start(n['server'], n['clusterid'], 0, n['portbase'])
        else:
            logger.DEBUG("!!!")
            for idx in range(n['startidx'], n['endidx']+1):
                logger.DEBUG("idx: %d", idx)
                fun_start(n['server'], n['clusterid'], idx, n['portbase'] + idx - 1)
    logger.DEBUG('run service start done')

def run_service_stop(options, names, kill=False):
    rnames = list(names)
    logger.DEBUG('run service stop(kill=%s) cmd:%s', kill, rnames)
    rnames.reverse()
    for n in rnames:
        pgname, idxs = gen_process_options(n)
        for i in idxs:
            options.pidx = i
            fun_stop(options, pgname, kill)
    logger.DEBUG('run service stop done')

def run_service_restart(options, names, kill=False):
    run_service_stop(options, names, kill)
    run_service_start(options, names)

def run_service_signal(options, names):
    logger.DEBUG('run service signal cmd:%s', names)
    for n in names:
        pgname, idxs = gen_process_options(n)
        for i in idxs:
            options.pidx = i
            fun_signal(options, pgname)
    logger.DEBUG('run service signal done')

def run_service_monitor(options, names):
    logger.DEBUG('run service monitor cmd:%s', names)
    logger.DEBUG("TODO")
    logger.DEBUG('run service monitor done')

def load_cluster_conf():
    global var_cluster_conf
    logger.DEBUG("load cluster conf...")
    from yaml import load
    cluster_conf_path = os.path.join(var_work_dir, '.cluster.yml')
    stream = file(cluster_conf_path, 'r')
    var_cluster_conf = load(stream)
    print var_cluster_conf

def gen_service_targets(args):
    logger.DEBUG("gen service targets...")
    targets = []

    if len(args) == 0:
        args = ['all']

    if len(args) == 1 and args[0] == 'all':
        targets = var_cluster_conf['servers']
        return targets

    for arg in args:
        found = False
        for server in var_cluster_conf['servers']:
            if server['clusterlayer'] == arg:
                targets.append(server)
                found = True
            elif server['server'] == arg:
                targets.append(server)
                found = True
        if found == False:
            logger.ERR("unsupported service target: %s", arg)  
            sys.exit(0)    

    return targets

###
def run_main(): 
    global var_work_dir

    hret = hqpy.HqError()
    hqpy.init(module="service")

    # options
    from optparse import OptionParser
    parser = OptionParser(usage=var_usage) 
    parser.add_option("-s", "--start",      dest="start",   default=False, help="safe start   service", action="store_true") 
    parser.add_option("-k", "--stop",       dest="stop",    default=False, help="safe stop    service", action="store_true") 
    parser.add_option("-K", "--kill",       dest="kill",    default=False, help="kill    service", action="store_true") 
    parser.add_option("-r", "--restart",    dest="restart", default=False, help="safe restart service", action="store_true") 
    parser.add_option("-R", "--krestart",   dest="krestart",default=False, help="kill and restart service", action="store_true") 
    parser.add_option("-l", "--signal",     dest="signal",  default=False, help="send signal to  service", action="store_true")
    parser.add_option("-L", "--sigval",     dest="sigval",  default=None,  help="signal val[u1, u2]")
    parser.add_option("-m", "--monitor",    dest="monitor", default=False, help="monitor service", action="store_true")

    (options, args) = parser.parse_args()
    options.tmo = 60     

    print "options: ", options, "args: ", args

    load_cluster_conf()

    targets = gen_service_targets(args)

    names = targets

    ## start
    if options.start:
        run_service_start(options, names)
        return

    ## stop
    if options.stop:
        run_service_stop(options, names)
        return 

    ## restart 
    if options.restart:
        run_service_restart(options, names)
        return

    ## krestart 
    if options.krestart:
        run_service_restart(options, names, True)
        return

    ## signal 
    if options.signal:
        run_service_signal(options, names)
        return

    ## kill
    if options.kill:
        run_service_stop(options, names, True)
        return

    ## monitor
    if options.monitor:
        hret = run_service_monitor(options, names)
        hqpy.check_exit(hret, "service monitor")
        return

    ## 
    logger.DEBUG("Do nothing.")

def main():
    try:
        run_main()
    except Exception,e: 
        print "PROGRAM RUN FAILED:"
        print traceback.format_exc()
        sys.exit(254)

    finally:
        pass


if __name__ == "__main__":
    main()

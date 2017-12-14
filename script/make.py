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
    from hqpy import hqio
except Exception, e:
    print 'HQPY IMPORT ERROR %s' % (e)
    sys.exit(255) ## PYRET_NEED_INIT
#############################################################################################

## 脚本变量及配置
##############################################################################################

var_project_conf    = {}
var_projects_all    = {}
var_build_targets   = {}
var_targets         = {}

var_work_dir        = var_work_dir.decode('utf-8')
var_src_dir         = os.path.join(var_work_dir, "src")
var_src_proj_dir    = os.path.join(var_src_dir, "project")
var_bin_dir         = os.path.join(var_work_dir, "bin")
var_pkg_dir         = os.path.join(var_work_dir, "pkg")

var_usage = u""" smake.py [options] TARGET(all|zone|world|tool)
wgame编译脚本使用说明
示例：
./smake.py -a           编译所有项目执行： 编译前预处理、清除编译结果、编译、打包
./smake.py -c           清除编译结果
./smake.py -p           执行编译前预处理，包括：proto转pb描述文件，proto转golang，lua错误码转golang
./smake.py -b           编译项目，TARGET可选项：all|zone|world|tool 或者是对应svr的名字 默认 all
./smake.py -t           打包项目安装文件
./smake.py -V           版本号
./smake.py -C           只打包配置
./smake.py -D           打包文件目标目录
./smake.py -i           安装protoc-gen-go 编译lua的C扩展 
"""
################################################################################################

## 执行shell命令
def run_shell(cmd, tmo=60):
    hret, ostd, oerr = hqpy.run_shell(cmd, tmo)
    if hret.iserr():
        logger.ERR('bash:%s\n%s\n%s', cmd, ostd, oerr)

    return hret

def make_func(pcfg):
    pname = pcfg["name"]
    ppath = pcfg["path"]

    logger.DEBUG(">>>>>>>> build %s", pname)
    srcdir = ppath
    binpath = os.path.join(var_bin_dir, pname)
    #print "make project: %s in path:%s" % (pname, srcdir)

    cmdstr = "go build -o %s %s" % (binpath, srcdir)
    hret = run_shell(cmdstr)
    if hret.iserr():
        return hqpy.HqError(13, "builder: target[%s] gobuild failed:%s" % (pname, hret.string()) )

    return hqpy.HqError()

def run_make_clean():
    global var_work_dir
    global var_build_targets

    os.chdir(var_work_dir)

    for k,p in var_build_targets.items():
        logger.DEBUG("<<<<<<<< clean %s", k)
        pname = p["name"]
        bindir = os.path.join(var_bin_dir, pname)
        if os.path.exists(bindir) == True:
            os.remove(bindir)

    return hqpy.HqError()

def run_make_build():
    global var_build_targets
    global var_work_dir
    os.chdir(var_work_dir)

    for k,p in var_build_targets.items():
        hret = make_func(p)
        if hret.iserr():
            #logger.ERR('build %s failed:%s', k, hret.string())
            return hret

        # # todo thread mode
    return hqpy.HqError()

def list_dirs(root):
    children = []
    for _, dirs, _ in os.walk(root):
        for name in dirs:
            children.append(name)
        break
    return children

def run_set_target(options, args):
    global var_build_targets
    global var_targets
    global var_src_proj_dir

    var_targets = {}
    var_targets['servers'] = {}
    var_targets['clusters'] = {}
    for cname in list_dirs(var_src_proj_dir):
        var_targets['clusters'][cname] = {}
        for sname in list_dirs(os.path.join(var_src_proj_dir, cname)):
            svr = {}
            svr['name'] = sname
            svr['path'] = os.path.join('project', cname, sname)
            var_targets['servers'][sname] = svr
            var_targets['clusters'][cname][sname] = svr

    targets = args

    nlen = len(targets)
    if nlen == 0:
        ## 默认全部
        targets = ["all"]
        nlen = 1

    var_build_targets = {}    
    if nlen == 1:
        target = targets[0]

        if target == "all":
            var_build_targets = var_targets['servers']
        elif var_targets['clusters'].has_key(target) == True:
            var_build_targets = var_targets['clusters'][target]
        elif var_targets['servers'].has_key(target) == True:
            var_build_targets[target] = var_targets['servers'][target]
        else:
            print var_targets['servers'].keys()
            return hqpy.HqError(hqpy.PYRET_SVR_MAKER, "invalid build target:%s" % (target))
    else:
        for target in targets:
            if var_targets['servers'].has_key(target):
                var_build_targets[target] = var_targets['servers'][target]
            else:
                return hqpy.HqError(hqpy.PYRET_SVR_MAKER, "invalid build target:%s" % (tarname))

    return hqpy.HqError()

def run_main(): 
    global var_build_projects
    global var_work_dir

    hret = hqpy.HqError()
    hqpy.init(module="smake")

    ## set GOPATH
    os.environ['GOPATH']=var_work_dir
    subprocess.call("echo GOPATH:$GOPATH", shell=True)

    # options
    from optparse import OptionParser
    parser = OptionParser(usage=var_usage) 
    parser.add_option("-t", "--total", dest="total", default=False, help="make clean; make build", action="store_true")
    parser.add_option("-c", "--clean", dest="clean", default=False, help="make clean", action="store_true") 
    parser.add_option("-b", "--build", dest="build", default=False,  help="make build", action="store_true")  
    (options, args) = parser.parse_args() 

    if options.total:
        options.clean = True
        options.build = True

    hret = run_set_target(options, args)
    hqpy.check_exit(hret, "set target")

    logger.DEBUG('hgame maker start.')

    print "options: ", options
    ## clean
    if options.clean:
        logger.DEBUG("make clean start")
        hret = run_make_clean()
        hqpy.check_exit(hret, "make clean")
        logger.DEBUG("make clean done")

    ## build
    if options.build:
        logger.DEBUG("make build start")
        hret = run_make_build()
        hqpy.check_exit(hret, "make build")
        logger.DEBUG("make build done")

    ## 
    logger.DEBUG('hgame maker done.')

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

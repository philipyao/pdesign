#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import traceback, shutil
import os, tarfile, time, zipfile
import hqpy
from hqpy import logger

def tar_zip(src, file_name, flist=None):

    cmdstr = ""
    if os.path.isfile(src):
        cmdstr = "tar cfz %s %s" % (file_name, src)
    else:
        cmdstr = "cd %s && tar cfz %s * " % (src, file_name)
        if flist is not None:
            strlist = ' '.join(flist)
            cmdstr = "cd %s && tar cfz %s %s " % (src, file_name, strlist)

    logger.DEBUG("tar_zip cmd:%s", cmdstr)
    hret, ostd, oerr = hqpy.run_shell(cmdstr, show_console=True, timeout=180)
    if hret.iserr():
        logger.ERR("tar_zip failed:%s,%s,%s \n%s", src, file_name, flist, oerr)

    return hret

    # 通过python实现的方法 速度太慢了
    # hret = hqpy.HqError()
    # try:
    #     tf = tarfile.open(file_name, 'w:gz')
    #     if os.path.isfile(src):
    #         #print os.path.basename(src)
    #         tf.add(src, os.path.basename(src))
    #     else:
    #         path = src
    #         if flist == None:
    #             flist = os.listdir(path)

    #         for f in flist:
    #             ff = os.path.join(path, f)
    #             tf.add(ff, f)

    #     tf.close()
    # except Exception,ex: 
    #     template = "An exception of type {0} occured. Arguments:{1!r}"
    #     message = template.format(type(ex).__name__, ex.args)
    #     hret.errno = hqpy.PYRET_IO_LIB_ERR
    #     hret.msg = 'IO tar Error:{0}'.format(message)
    #     logger.ERR(traceback.format_exc())
    # finally:
    #     return hret

def mkdirs(dirname):
    if os.path.exists(dirname):
        return

    os.makedirs(dirname)

def remove(fname):
    hret = hqpy.HqError()
    if os.path.exists(fname) == False:
        return hret

    try:
        if os.path.isfile(fname):
            os.remove(fname)
        else:
            shutil.rmtree(fname)
    except Exception, ex:
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_IO_LIB_ERR
        hret.msg = 'IO remove Error:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret

def copy_from_to(src_dir, fname, dst_dir):
    hret = hqpy.HqError()
    if type(fname) == str:
        fname = [fname]

    for f in fname:
        fsrc = os.path.join(src_dir, f)
        if os.path.exists(fsrc) == False:
            hret.errno = hqpy.PYRET_IO_LIB_ERR
            hret.msg = 'IO copyft Error:{0} not exist'.format(fsrc)
            return hret

        if os.path.isfile(fsrc):
            dst = os.path.join(dst_dir, f)
            ndir = os.path.dirname(dst)
            if os.path.exists(ndir) == False:
                mkdirs(ndir)
            shutil.copyfile(fsrc, dst)
        else:
            dst = os.path.join(dst_dir, f)
            shutil.copytree(fsrc, dst)

    return hret

def zip(src, dst, pwd=None, tar=False, ziped=False):
    hret = hqpy.HqError()

    compress = zipfile.ZIP_STORED
    if ziped == True:
        compress = zipfile.ZIP_DEFLATED

    zipf = zipfile.ZipFile(dst, 'w', compress)

    if os.path.isfile(src):
        zipf.write(src)
    else:
        # ziph is zipfile handle
        for root, dirs, files in os.walk(src):
            for f in files:
                zipf.write(os.path.join(root, f))

    zipf.close()

    if pwd is not None:
        import pyminizip
        fenc = dst + ".tpwd"
        pyminizip.compress(dst, fenc, pwd, 0)
        os.remove(dst)
        shutil.move(fenc, dst)

    if tar == True:
        tar_name = dst + ".tgz"
        hret = tar_zip(dst, tar_name)
        if hret.iserr():
            return hret

        os.remove(dst)

    return hret

    

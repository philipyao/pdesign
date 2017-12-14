# Copyright (C) 2015-2015  tonyliu <tonyliu@hanqugame.com>
#
# This file is part of hqpy.

import sys,os
from hqpy._version import __version__, __version_info__

if sys.version_info < (2, 6):
    raise RuntimeError('You need Python 2.6+ for this module.')

__author__ = "hanqugame <tonyliu@hanqugame.com>"
__license__ = ""

from hqpy.funcs     import *
from hqpy.consts    import *
from hqpy.hqerr     import *
from hqpy.ssh       import *

__all__ = [ 
    'SshAuth',
    'SshClient',
    'HqError',
    'funcs' ]

def init(module="hat", load_host=False, jid=None, jname=None, show_console=True, tm_tag=True):
    from hqpy import hqenv
    from hqpy import logger
    from hqpy import hqhosts
    from hqpy import hqio

    if os.environ.has_key('HATROOT'):
        hat_work_dir = os.environ['HATROOT']
    else:
        hat_work_dir = "/var/tmp/hat"
        hqio.mkdirs(hat_work_dir)


    hqenv.set_var('HQVAR_WORK_DIR', hat_work_dir)
    hqenv.set_var('HQVAR_LOG_DIR', os.path.join(hat_work_dir, 'log'))
    hqenv.set_var('HQVAR_SRC_DIR', os.path.join(hat_work_dir, 'src'))
    hqenv.set_var('HQVAR_BIN_DIR', os.path.join(hat_work_dir, 'bin'))
    hqenv.set_var('HQVAR_CONF_DIR', os.path.join(hat_work_dir, 'conf'))
    hqenv.set_var('HQVAR_VAR_DIR', os.path.join(hat_work_dir, 'var'))
    hqenv.set_var('HQVAR_RLS_DIR', os.path.join(hat_work_dir, 'files/release'))
    hqenv.set_var('HQVAR_FILES_DIR', os.path.join(hat_work_dir, 'files'))
    hqenv.set_var('HQVAR_SCRIPT_DIR', os.path.join(hat_work_dir, 'script'))
    hqenv.set_var('HQVAR_HGAME_DIR', os.path.abspath(os.path.join(hat_work_dir, "..")))

    hqio.mkdirs(hqenv.get_var('HQVAR_LOG_DIR'))

    if jid is not None:
        #job mode
        hqenv.set_var('HQVAR_JOBID', jid)
        logger.init_job_logger(hqenv.get_var('HQVAR_LOG_DIR'), jid)
    else:
        #normal mode
        hqenv.set_var('HQVAR_JOBID', 0)
        logger.init_normal_logger(hqenv.get_var('HQVAR_LOG_DIR'), module, show_console, tm_tag)       

    # load env.json
    fini = os.path.join(hqenv.get_var('HQVAR_CONF_DIR'), 'env.ini')
    hret, envcfg = hqpy.loadini(fini)
    if hret.iserr():
        logger.ERR("load env.json failed:%s", hret.string())
        return False
    else:
        hqenv.set_var('HQVAR_ENV_CFG', envcfg)
    
    # load host cfg
    if load_host == True:
        hret = hqhosts.init_hosts()
        if hret.iserr():
            logger.ERR('hqhosts init failed:%s', hret.string())
            return False

    #logger.DEBUG('ALL HQPY VAR:%s', hqenv.get_all())
    return True
#!/usr/bin/env python
# -*- coding: UTF-8 -*-

### PY 状态码定义
###############################
PYRET_OK    = 0     # 成功
PYRET_ERR   = 1     # 通用失败

## 远程本地状态码
PYRET_RMIN          = 1     # 最小
PYRET_RERR          = 1 


PYRET_TIMEOUT           = 100   # SSH 超时
PYRET_SSH_RUNERR        = 101 
PYRET_SSH_PUTERR        = 102   # PUT 异常
PYRET_SSH_GETERR        = 103   # GET 异常
PYRET_TASK_LOAD_ERR     = 120
PYRET_TASK_CHECK_ERR    = 121
PYRET_CMD_RUN_ERR       = 122

PYRET_INSTALL_ERR       = 248
PYRET_IO_LIB_ERR        = 249
PYRET_SERVICE_ERR       = 250
PYRET_SVR_MAKER         = 251
PYRET_JOB_LOCK_ERR      = 252
PYERT_JOB_FAILED        = 253
PYRET_ERR_EXIT          = 254
PYRET_NEED_INIT         = 255
PYRET_RMAX              = 255   # 最大

class Object(object):
    pass

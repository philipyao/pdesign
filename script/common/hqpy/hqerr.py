#!/usr/bin/env python
# -*- coding: UTF-8 -*-

class HqError():  
    def __init__(self, errno=0, msg = ""):
        self.errno = errno
        self.msg = msg

    def check(self, err):
        return self.errno == err

    def isok(self):
        return self.errno == 0

    def iserr(self):
        return self.errno != 0

    def string(self):
        return '[Hqerr {0} {1}]'.format(self.errno, self.msg)

    def code(self):
        return self.errno 

def OK():
    return HqError(errno=0)

def Err(msg=""):
    return HqError(errno=254, msg=msg)
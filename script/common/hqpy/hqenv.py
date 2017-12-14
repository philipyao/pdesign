#!/usr/bin/env python
# -*- coding: UTF-8 -*-

global_env_var = {}

def set_var(key, val):
    global_env_var[key] = val

def has_var(key):
    return global_env_var.has_key(key)

def get_var(key):
    if has_var(key) == True:
        return global_env_var[key]
    else:
        return None

def get_all():
    return global_env_var
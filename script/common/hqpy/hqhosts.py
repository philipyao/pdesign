#!/usr/bin/env python
# -*- coding: UTF-8 -*-

import os,json,traceback

import hqpy
from hqpy import logger
from hqpy import hqenv


__hosts_by_id = {}
__hosts_by_group = {}
__deploy_by_id= {}
__deploy_by_role= {}
__deploy_by_kind={}

class HostAuthObj(object):

    def __init__(self):
        self.id = ""
        self.ip = ""
        self.port = ""
        self.kind = ""
        self.islocal = False

        self.user = None
        self.root = None
    def __repr__(self):
        return "%s" % (self.ip)

    def __str__(self):
        return "%s" % (self.ip)

    def info(self):
        return 'id:{0}, host:{1}:{2}, kind:{3}'.format(self.id, self.ip, self.port, self.kind)

    def get(self, broot=False):
        if broot:
            return self.root
        else:
            return self.user

def init_hosts():
    #logger.DEBUG('host cfg load start.')
    fd = None
    hret = hqpy.HqError()
    try:
        cfg_file = os.path.join(hqenv.get_var('HQVAR_CONF_DIR'), 'hosts.json')

        fd = open(cfg_file)
        jscfg = json.load(fd)

        def_auth = jscfg['default']
        default_user = def_auth['user']
        default_passwd = def_auth['pwd']
        default_root_user = def_auth['root_user']
        default_root_passwd = def_auth['root_pwd']

        hosts = jscfg['hosts']
        for h in hosts:
            hid = h['id']
            ip = h['ip']
            kind = h['kind']
            port = 22
            if h.has_key('port'):
                port = h['port']
            user = default_user
            pwd = default_passwd
            if h.has_key('user'):
                user = h['user']
            if h.has_key('pwd'):
                pwd = h['pwd']
            ruser = default_root_user
            rpwd = default_root_passwd
            if h.has_key('ruser'):
                ruser = h['ruser']
            if h.has_key('rpwd'):
                ruser = h['rpwd'] 
                
            hauth = HostAuthObj()
            user_auth = hqpy.SshAuth(ip, user, pwd, port)
            root_auth = hqpy.SshAuth(ip, ruser, rpwd, port)
            hauth.user = user_auth
            hauth.root = root_auth

            hauth.id = hid
            hauth.ip = ip
            hauth.port = port
            hauth.kind = kind

            if __hosts_by_id.has_key(hid):
                return hqpy.Err("load host, exist hid:%s" % (hid))

            __hosts_by_id[hid] = hauth
            for g in kind:
                if __hosts_by_group.has_key(g) == False:
                    __hosts_by_group[g] = []
                    __deploy_by_kind[g] = []
                __hosts_by_group[g].append(hauth)
                __deploy_by_kind[g].append(hid)

            __deploy_by_id[hid] = {}
            if h.has_key('zone'):
                v = h['zone']
                for x in v:
                    rid = x[0]
                    rseq = x[1]
                    rname = 'zone' + str(rid)
                    if __deploy_by_role.has_key(rname):
                        hret = hqpy.Err("load host, exist zone name:%s" % (rname) )
                        return hret
                    __deploy_by_role[rname] = {'sau':hauth, 'info':h, 'role':'zone', 'roles':[[rid, rseq]] }
                __deploy_by_id[hid]['zone'] = {'sau':hauth, 'info':h, 'role':'zone', 'roles':h['zone'] }

            if h.has_key('group'):
                v = h['group']
                for x in v:
                    rid = x[0]
                    rseq = x[1]
                    rname = 'group' + str(rid)
                    if __deploy_by_role.has_key(rname):
                        hret = hqpy.Err("load host, exist group name:%s" % (rname))
                        return hret
                    __deploy_by_role[rname] = {'sau':hauth, 'info':h, 'role':'group', 'roles':[[rid, rseq]] }
                __deploy_by_id[hid]['group'] = {'sau':hauth, 'info':h, 'role':'group', 'roles':h['group'] } 

            if h.has_key('global'):
                v = h['global']
                for x in v:
                    rid = x[0]
                    rseq = x[1]
                    rname = 'global' + str(rid)
                    if __deploy_by_role.has_key(rname):
                        hret = hqpy.Err("load host, exist global name:%s" % (rname) )
                        return hret

                    __deploy_by_role[rname] = {'sau':hauth, 'info':h, 'role':'global', 'roles':[[rid, rseq]] }
                __deploy_by_id[hid]['global'] = {'sau':hauth, 'info':h, 'role':'global', 'roles':h['global'] }  

            if __deploy_by_id.has_key('hid'):
                hret = hqpy.Err("load host, exist hostid :%s" % (hid))
                return hret
            
        #print __hosts_by_id
        #print __hosts_by_group
        #print __deploy_by_id
        #print __deploy_by_role
        #logger.DEBUG('host cfg load done.')

    except Exception,ex: 
        template = "An exception of type {0} occured. Arguments:{1!r}"
        message = template.format(type(ex).__name__, ex.args)
        hret.errno = hqpy.PYRET_TASK_LOAD_ERR
        hret.msg = 'TASK Error:{0}'.format(message)
        logger.ERR(traceback.format_exc())
    finally:
        return hret

def get_host_cfg():
    dkind = []
    idx=1000
    for k,v in __deploy_by_kind.iteritems():
        kidx = idx
        children=[]
        for hid in v:
            kidx=kidx+1
            children.append({'id':kidx, 'text':hid})
        dkind.append({'id':idx, 'text':k, 'children':children})
        idx = idx + 1000

    dhid = []
    idx=1000
    for k,v in __deploy_by_id.iteritems():
        hidx=idx
        children=[]
        for _,r in v.iteritems():
            for x in r['roles']:
                hidx = hidx +1
                rname = "%s.%s" % (r['role'], x[0])
                children.append({'id':hidx, 'text':rname})
        dhid.append({'id':idx, 'text':k, 'children':children})
        idx = idx + 1000

    return dkind, dhid

def get_host_byid(hid):
    if __hosts_by_id.has_key(hid) == True:
        return __hosts_by_id[hid]
    return None

def get_host_bykind(kind):
    if __hosts_by_group.has_key(kind) == True:
        return __hosts_by_group[kind]
    return None

def get_local_host():
    obj = HostAuthObj()
    obj.id = "localhost"
    obj.ip = "127.0.0.1"
    obj.port = "22"
    obj.kind = "local"
    obj.islocal = True

    obj.user = hqpy.SshAuth(obj.ip, "user", "", 22)
    obj.root = hqpy.SshAuth(obj.ip, "root", "", 22)

    return obj

def get_dst_byrole(role, rid):
    rname = "%s%s" % (role, rid)
    if __deploy_by_role.has_key(rname):
        return __deploy_by_role[rname]
    return None

def get_dst_byid(role, hid):
    if __deploy_by_id.has_key(hid):
        rs = __deploy_by_id[hid]
        if rs.has_key(role):
            return rs[role]
    return None
    
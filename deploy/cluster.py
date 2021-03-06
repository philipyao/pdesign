#!/usr/bin/env python
# -*- coding=utf-8 -*-

import sys
import MySQLdb as mdb
import yaml

db_cnx = None

def init_db_handle():
    try:
        global db_cnx
        db_cnx = mdb.connect('10.1.164.20', 'hgame', 'Hgame188', 'db_new_oms');
    except mdb.Error, e:
        print "Error %d: %s" % (e.args[0],e.args[1])
        sys.exit(1)

def load_host():
    try:
        cur = db_cnx.cursor(mdb.cursors.DictCursor)
        cur.execute("SELECT * FROM tbl_host")

        hosts = {}
        dbhosts = cur.fetchall()
        if len(dbhosts) == 0:
            print "no hosts"
            sys.exit(1)
        for hentry in dbhosts:
            hosts[hentry["hostid"]] = hentry
        return hosts
             
    except mdb.Error, e:
  
        print "Error %d: %s" % (e.args[0],e.args[1])
        sys.exit(1)
    
def load_cluster():
    try:
        cur = db_cnx.cursor(mdb.cursors.DictCursor)
        sqlcmd = "SELECT * FROM tbl_cluster"
        cur.execute(sqlcmd)
        dbcluster = cur.fetchall()
        if len(dbcluster) == 0:
            print "no clusters"
            sys.exit(1)
        clusters = {}
        for c in dbcluster:
            clusters[c["cluster"]] = c
        return clusters
    except mdb.Error, e:
        print "Error %d: %s" % (e.args[0],e.args[1])
        sys.exit(1)

def load_server(cluster_name):
    try:
        cur = db_cnx.cursor(mdb.cursors.DictCursor)
        sqlcmd = "SELECT * FROM tbl_server where cluster = '%s'" % cluster_name
        cur.execute(sqlcmd)
        dbserver = cur.fetchall()
        if len(dbserver) == 0:
            print "no servers found for cluster %s" % cluster_name
            sys.exit(1)
        return dbserver
    except mdb.Error, e:
        print "Error %d: %s" % (e.args[0],e.args[1])
        sys.exit(1)

if __name__ == '__main__':
    from optparse import OptionParser
    parser = OptionParser()
    parser.add_option("-p", "--sshpass",  
                  action="store", dest="sshpass", type="string",  
                  help="password for ssh to remote host")  
    parser.add_option("-t", "--target",  
                  action="store", dest="target", type="string",  
                  help="center host target")  
    (options, args) = parser.parse_args() 

    init_db_handle()
    hosts = load_host()
    clusters = load_cluster()
    cluster_name = 'zone2001'
    cluster = clusters[cluster_name]
    if cluster == None:
        print "cluster %s not found in db" % cluster_name
        sys.exit(1)

    print "cluster: ", cluster
    host = hosts[cluster['host']]
    if host == None:
        print "cluster %s host=%s not found!" % (cluster_name, cluster['host'])

    print "host: ", host

    servers = load_server(cluster_name)
    print "servers: ", servers

    # 生成yaml格式cluster配置
    cfgname = '../.cluster.yml'
    yaml_obj = {}
    yaml_obj['ip'] = host['ip']
    yaml_obj['wanip'] = host['wanip']
    yaml_obj['servers'] = []
    for s in servers:
        svr = {}
        svr['clusterlayer'] = cluster['clusterlayer']
        svr['clusterid'] = int(cluster['clusterid'])
        svr['server'] = s['typename']
        svr['startidx'] = int(s['startidx'])
        svr['endidx'] = int(s['endidx'])
        svr['portbase'] = int(s['portbase'])

        yaml_obj['servers'].append(svr)

    print yaml_obj
    with open(cfgname, 'w') as outfile:
        yaml.dump(yaml_obj, outfile, default_flow_style=False)    

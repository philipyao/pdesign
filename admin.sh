#!/bin/bash
##################################################################################
workdir=$(cd `dirname $0`; pwd)  
source $HOME/.bashrc

export HAT_DIR=$workdir
export HAT_HQSH_DIR=$HAT_DIR/scripts/common/hqshell
source $HAT_HQSH_DIR/import.sh
## import VARS: WORKDIR  LLLOCALIP
##################################################################################
cd $workdir
export GOPATH=$workdir
CEXTPATH=$workdir/c/
RESPATH=$workdir/resource
LUARES=$RESPATH/luares
JSONRES=$RESPATH/jsonres
PBRES=$RESPATH/pb
BTLSTAGERES=$RESPATH/btlstage
ORGLUARES=$RESPATH/luacfg

var_user_opt=$1
fenv=$workdir/cfg/envrc

if [[ "$var_user_opt" != "cfg" ]] && [[ "$var_user_opt" != "help" ]] &&[[ "$var_user_opt" != "loc" ]] && [[ "$var_user_opt" != "init" ]]; then
    if [[ ! -f $fenv ]]; then
        hq_exit "cfg/envrc not found."
    fi
    source $fenv
fi

hq_debug "HGAME_TARGETS:$HGAME_TARGETS"

target=""

usage()
{
    echo  "init        初始公用代码及配置目录"
    echo  "loc         设置版本地区 cn tw hk"
    echo  "db|cleardb  清空数据库 请谨慎操作"
    echo  "ossdb       OSS数据库操作 请谨慎操作"
    echo  "osssql      OSS数据库SQL生成"
    echo  "cfg         生成开发配置"
    echo  "make | m    编译项目 调用smake.py     make TARGET "
    echo  "start | s   启服务                    s  TARGET"
    echo  "stop  | k   停服务                    k  TARGET"
    echo  "stop9 | K   强制停服务                k  TARGET"
    echo  "restart | r  重启服务                 r  TARGET"
    echo  "restart9| r9 强制停服并重启           r9 TARGET"
    echo  "redis       启动redis server" 
    echo  "externals   初始化svn:externals"

    echo  "TARGET： 默认 zone 可选[zone world global all] 或者是多个svr名字 例如 gamesvr gatesvr gmsvr"
}

var_set=
if [[ -n "$HGAME_WORLDID" ]]; then
    var_set="-x world$HGAME_WORLDID"
fi
if [[ -n "$HGAME_ZONEID" ]]; then
    var_set="-x zone$HGAME_ZONEID"
fi

loc_conf()
{
    loc=$1
    if [[ -z $loc ]]; then
        loc="cn"
    fi
    var_conf="conf"
    var_proto="proto"
    var_btlstage="btlstage"
    if [[ "$loc" != "cn" ]]; then
        var_conf="conf-$loc"
        var_proto="proto-$loc"
        var_btlstage="btlstage-$loc"
    fi

    if [[ ! -d $workdir/common/$var_conf ]]; then
        hq_exit "location conf error:$loc"
    fi 
    if [[ ! -d $workdir/common/$var_proto ]]; then
        hq_exit "location proto error:$loc"
    fi 

    echo "loc conf:$loc, $var_conf, $var_proto, $var_btlstage"
    cd $workdir
    rm -rf $LUARES
    rm -rf $JSONRES
    rm -rf $PBRES
    rm -rf $BTLSTAGERES
    rm -rf $ORGLUARES

    ln -s $workdir/common/$var_conf/export/cfg_server/lua $LUARES
    ln -s $workdir/common/$var_conf/export/cfg_server/json $JSONRES
    ln -s $workdir/common/$var_proto/export/bin $PBRES
    ln -s $workdir/common/$var_btlstage $BTLSTAGERES
    ln -s $workdir/common/$var_conf/original/lua $ORGLUARES

    echo "loc conf done."
}

gen_target()
{
    option=$1
    if [[ -z "$option" ]]; then 
        target=$HGAME_TARGETS
    elif [ "$option"x = "global"x ]; then
        target=$HGAME_GLOBAL_TARS
    elif [ "$option"x = "world"x ]; then
        target=$HGAME_WORLD_TARS
    elif [ "$option"x = "zone"x ]; then
        target=$HGAME_ZONE_TARS
    else
        OIFS=$IFS
        IFS=" "
        tarlist=$HGAME_TARGETS
        for x in $tarlist; do
            [[ $x = $option* ]] && target=$x && break
        done
        IFS=$OIFS
    fi
    if [ -z "$target" ]; then
        hq_exit "no target found in envrc: $@"
    fi
}

main()
{
    chmod +x *.sh *.py
    shift 1

    case $var_user_opt in

        init)
            ## developer only
            cd $workdir

            loc_conf "cn"

            python ./smake.py -i
            cd $workdir
            rm -rf $CEXTPATH/build/*
            ;;
        loc)
            loc_conf "$@"
            ;;

        ldoc)
            source ~/.bashrc
            lua $HATROOT/LDoc-master/ldoc.lua -d lua/rpcdoc lua/gamesvr/rpc
            ;;

        f | find)
            ## developer only
            cd $workdir
            sh ./scripts/05_grep.sh "$@" 
            ;;

        db | cleardb )
            cd $workdir
            chmod +x $workdir/cfg/mysql/sql.sh 
            vopt=$1
            if [[ "$vopt" != "QUIET" ]]; then
                read -p "!!! Enter YES to continue:" uinput
                if [[ "$uinput" != "YES" ]]; then
                    hq_debug "cancel cleardb"
                    exit 0
                fi
            fi

            dothing=0
            if [[ -n "$HGAME_ZONEID" ]] && [[ -n "$HGAME_AU_ZONE" ]]; then
                $workdir/cfg/mysql/sql.sh -r zone -c create -z $HGAME_ZONEID -a "${HGAME_AU_ZONE}"
                $workdir/cfg/mysql/sql.sh -r zone -c drop -z $HGAME_ZONEID -a "${HGAME_AU_ZONE}"
                $workdir/cfg/mysql/sql.sh -r zone -c create -z $HGAME_ZONEID -a "${HGAME_AU_ZONE}"
                dothing=1
            fi
            if [[ -n "$HGAME_WORLDID" ]] && [[ -n "$HGAME_AU_WORLD" ]]; then
                $workdir/cfg/mysql/sql.sh -r world -c create -w $HGAME_WORLDID -a "${HGAME_AU_WORLD}"
                $workdir/cfg/mysql/sql.sh -r world -c drop -w $HGAME_WORLDID -a "${HGAME_AU_WORLD}"
                $workdir/cfg/mysql/sql.sh -r world -c create -w $HGAME_WORLDID -a "${HGAME_AU_WORLD}"
                $workdir/cfg/mysql/sql.sh -r world -c custom -w $HGAME_WORLDID -a "${HGAME_AU_WORLD}"
                dothing=1
            fi
            if [[ -n "$HGAME_GLOBALIDX" ]] && [[ -n "$HGAME_AU_GLOBAL" ]]; then
                $workdir/cfg/mysql/sql.sh -r global -c create  -a "${HGAME_AU_GLOBAL}"
                $workdir/cfg/mysql/sql.sh -r global -c drop  -a "${HGAME_AU_GLOBAL}"
                $workdir/cfg/mysql/sql.sh -r global -c create -a "${HGAME_AU_GLOBAL}"
                dothing=1
            fi

            if [ $dothing -eq 0 ]; then
                hq_debug "warning: nothing to do"
            fi

            ;;

        ossdb )
            cd $workdir
            chmod +x $workdir/cfg/mysql/sql.sh 
            cmd=$1
            role=$2
            zoneid=$3

            if [[ $# -lt 2 ]]; then
                hq_debug "ossdb create|drop|delete  oss_global|oss_zone|oss_stat [ZONEID]"
                exit 0
            fi

            read -p "!!! Enter YES to continue:" uinput
            if [[ "$uinput" != "YES" ]]; then
                hq_debug "cancel cleardb"
                exit 0
            fi

            if [[ "$role" == "oss_zone" ]] && [[ -n "$zoneid" ]] && [[ -n "$HGAME_AU_OSS_ZONE" ]]; then
                $workdir/cfg/mysql/sql.sh -r oss_zone -c $cmd -z $zoneid -a "${HGAME_AU_OSS_ZONE}"
            fi

            if [[ "$role" == "oss_global" ]] && [[ -n "$HGAME_AU_OSS_GLOBAL" ]]; then
                $workdir/cfg/mysql/sql.sh -r oss_global -c $cmd -a "${HGAME_AU_OSS_GLOBAL}"
            fi

            if [[ "$role" == "oss_stat" ]] && [[ -n "$HGAME_AU_OSS_GLOBAL" ]]; then
                $workdir/cfg/mysql/sql.sh -r oss_stat -c $cmd -a "${HGAME_AU_OSS_GLOBAL}"
            fi

            hq_debug "ossdb done"
            ;;

        osssql )
            cd $workdir
            lua ./scripts/11_oss.lua "$@"
            hq_debug "osssql done"
            ;;

        cfg )
            cd $workdir
            lua ./scripts/08_gen_cfg.lua -p $workdir
            hq_check_exit "cfg"
            ;;

        make | m )
            cd $workdir
            python ./smake.py -c -p -b "$@"
            ;;

        start | s )
            cd $workdir
            gen_target "$@"
            python ./service.py $var_set  -s $target 
            hq_check_exit "start"
            ;;

        stop | k)
            gen_target "$@"
            python ./service.py  -k $target 
            hq_check_exit "stop"
            ;;

        stop9 | K)
            gen_target "$@"
            python ./service.py -K $target 
            hq_check_exit "kill"
            ;;

        restart | r)
            gen_target "$@"
            python ./service.py $var_set -r $target 
            hq_check_exit "restart"
            ;;

        restart9 | r9)
            gen_target "$@"
            python ./service.py  -K $target 
            python ./service.py $var_set  -s $target 
            hq_check_exit "restart9"
            ;;

        initstart )
            flog=./log/reboot.log
            echo "$(date): reboot start" >> $flog
            # 先启动本服其它服务
            if [[ -e ./init.local ]]; then
                echo "$(date): init local" >> $flog
                sh ./init.local
            fi
            # 带监控启动OMC
            export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$HOME/hgame/lib:./lib ; ./bin/omc -on >> ./log/core.omc.log 2>&1 &

            echo "$(date): reboot done" >> $flog
            ;;
        redis )
            cd $workdir 
            mkdir -p ../redisdb
            cp  ./cfg/ctpls/redis.conf.tpl ./cfg/redis.conf
            redis-server ./cfg/redis.conf > /dev/null &
            exit 1
            ;;
        externals )
            sh ./scripts/13_externals.sh -p $workdir
            ;;
        help | h)
            usage
            ;;
        
        * )  
            echo "Invalid opt: $var_user_opt"
            usage
            exit 1
            ;;  
    esac
}

main "$@"

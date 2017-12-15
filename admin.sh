#!/bin/bash
##################################################################################
workdir=$(cd `dirname $0`; pwd)  
source $HOME/.bashrc
##################################################################################
cd $workdir
export GOPATH=$workdir

var_user_opt=$1

target=""

usage()
{
    echo  "make | m    编译项目 调用make.py     make TARGET "
    echo  "start | s   启服务                    s  TARGET"
    echo  "stop  | k   停服务                    k  TARGET"
    echo  "stop9 | K   强制停服务                k  TARGET"
    echo  "restart | r  重启服务                 r  TARGET"
    echo  "restart9| r9 强制停服并重启           r9 TARGET"

    echo  "TARGET： 默认 zone 可选[zone world all] 或者是多个svr名字 例如 gamesvr gatesvr gmsvr"
}

main()
{
    chmod +x *.sh *.py
    shift 1

    case $var_user_opt in

        make | m )
            cd $workdir
            python ./script/make.py -t "$@"
            ;;

        start | s )
            cd $workdir
            python ./script/service.py -s "$@"
            ;;

        stop | k)
            python ./script/service.py -k "$@"
            ;;

        stop9 | K)          
            python ./script/service.py -K "$@"
            ;;

        restart | r)
            python ./script/service.py -r "$@"
            ;;
            
        restart9 | r9)
            python ./script/service.py  -K $target
            python ./script/service.py  -s $target 
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

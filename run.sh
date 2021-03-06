#!/bin/sh

DIR="./"

typeset -l COMMAND
COMMAND=${1}

typeset -l PROCESS_NAME
PROCESS_NAME="tuyue_gateway"

startProcess() {
    echo "Starting service ..."
    log="./log/log_`date +%Y%m%d_%H%M%S`.log"
    nohup ${DIR}${PROCESS_NAME} > ${log} 2>&1 &
    echo "Starting service success"
    Pid=$(queryProcessPid ${PROCESS_NAME})
    if [[ ${Pid} ]]; then
        echo pid:[${Pid}]
    fi
}

stopProcess() {
    Pid=${1}
    if [[ ${Pid} ]]; then
        for id in ${Pid}
        do
            kill -s 9 ${id}
        done
    fi
}

queryProcessPid() {
    if [[ ${1} ]]; then
        Pid=`ps -ef | grep ${1} | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
        if [[ "$Pid" ]]; then
            echo ${Pid}
        fi
    fi
}

if [[ ${COMMAND} ]]; then
    Pid=$(queryProcessPid ${PROCESS_NAME})
    if [[ ${Pid} ]]; then
        echo pid:[${Pid}]
    fi

    if [[ ${COMMAND} = "restart" ]]; then
        stopProcess ${Pid}
        sleep 1s
        startProcess
    elif [[ ${COMMAND} = "start" ]]; then
        if [[ "$Pid" ]]; then
            echo "service running"
        else
            startProcess
        fi
    elif [[ ${COMMAND} = "stop" ]]; then
        stopProcess ${Pid}
        echo "ok"
    else
        echo "please input [start/restart/stop] command"
    fi
else
    echo "please input [start/restart/stop] command"
fi

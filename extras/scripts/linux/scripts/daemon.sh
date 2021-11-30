#!/usr/bin/env bash

# Copyright (c) 2019-2029 Dunyu All Rights Reserved.
#
# Author      : yangping
# Email       : ping.yang@wengold.net
# Version     : 1.0.1
# Description :
#   This script to start, stop and show status of server.
#
# Prismy.No | Date       | Modified by. | Description
# -------------------------------------------------------------------
# 00001       2021/08/29   yangping       New version
# -------------------------------------------------------------------

usage="
Usage: daemon.sh (start|stop|status) (server_name) \
    start  : start server.
    stop   : stop server.
               -f (-force) : force to stop server process.
    status : to show the process number of running server.
"

# if no args specified, show usage
if [ $# -le 1 ]; then
  echo $usage
  exit 1
fi

# enter service root path
scripts=`dirname "$0"`
scripts=`cd "$scripts"; pwd`
SERVER_ROOT=`cd "${scripts}/../.."; pwd`

# get (start|stop|status) arguments
commands=$1
shift

# get server app name
SERVER_APP=$1
shift

umask 0000

# parse the '-f' option of stop command
forced=$1
FORCED_KILL=false
shopt -s extglob
if [ ! -z ${forced} ]; then
  case ${forced} in
    *(-)f )
      FORCED_KILL=true
      ;;
    *(-)force )
      FORCED_KILL=true
      ;;
  esac
fi

# create log directory
LogDir="${SERVER_ROOT}/logs"
mkdir -p "$LogDir"

pid=${LogDir}/${SERVER_APP}.pid

# dispath (start|stop|status) commands
case $commands in

  (start)
    if [ -f $pid ]; then
      if ps -p `cat $pid` > /dev/null 2>&1; then
        echo server running as process `cat $pid`.
        exit 1
      fi
    fi

    cd ${SERVER_ROOT}
    ./${SERVER_APP} >&1 </dev/null &
    echo $! > ${pid}
    ;;

  (stop)
    if [ -f $pid ]; then
      if ps -p `cat $pid` > /dev/null 2>&1; then
        if ! kill -0 `cat $pid` > /dev/null 2>&1; then
          echo cannot stop server with pid `cat $pid` - permission denied
        elif ${FORCED_KILL}; then
          kill -9 `cat $pid` > /dev/null 2>&1;
          sleep 1;
          echo server killed
        else
          echo -n stopping server
          kill `cat $pid` > /dev/null 2>&1
          while ps -p `cat $pid` > /dev/null 2>&1; do
            echo -n "."
            sleep 1;
          done
          echo
          echo stoped server
        fi
        if ! ps -p `cat $pid` > /dev/null 2>&1; then
          rm $pid
        fi
      else
        echo server not runging
      fi
    else
      echo server not runging
    fi
    ;;

  (status)
    if [ -f $pid ]; then
      if ps -p `cat $pid` > /dev/null 2>&1; then
        echo server running as process `cat $pid`.
      else
        echo server not running.
      fi
    else
      echo server not running.
    fi
    ;;

  (*)
    echo $usage
    exit 1
    ;;

esac

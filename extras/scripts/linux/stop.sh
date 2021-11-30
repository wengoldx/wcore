#!/usr/bin/env bash

# Copyright (c) 2019-2029 Dunyu All Rights Reserved.
#
# Author      : yangping
# Email       : ping.yang@wengold.net
# Version     : 1.0.1
# Description :
#   Stop server.
#
# Prismy.No | Date       | Modified by. | Description
# -------------------------------------------------------------------
# 00001       2021/08/29   yangping       New version
# -------------------------------------------------------------------

bin=`dirname "$0"`
bin=`cd "$bin"; pwd`

source ${bin}/scripts/exports.sh
${bin}/scripts/daemon.sh stop ${SERVICE_APP_NAME}

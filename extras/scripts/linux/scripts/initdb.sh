#!/usr/bin/env bash

# Copyright (c) 2019-2029 Dunyu All Rights Reserved.
#
# Author      : yangping
# Email       : ping.yang@wengold.net
# Version     : 1.0.1
# Description :
#   Create database for server.
#
# Prismy.No | Date       | Modified by. | Description
# -------------------------------------------------------------------
# 00001       2021/08/29   yangping       New version
# -------------------------------------------------------------------

# enter service bin folder
bin=`dirname "$0"`
bin=`cd "$bin"; pwd`
source ${bin}/exports.sh
echo "Starting init database..."

# check database user, receive input if empty
if [ -z "$SERVICE_DB_USER" ]; then
  read -p "Enter database user: " SERVICE_DB_USER

  # check the input data
  if [ -z "$SERVICE_DB_USER" ]; then
    echo "Database user must config or input!"
    exit 1
  fi
fi

# enter script path, init database tables
cd  ${bin} 
for file in `ls *.sql`
do 
  mysql -h $SERVICE_DB_HOST -P 3306 -u $SERVICE_DB_USER -p < ${file} --default-character-set=utf8mb4
  if [[ $? == 1 ]]; then
    echo "Failed init ${file} database!"
    exit 1
  fi
done

# exit script
echo "Finished init database."
exit 0

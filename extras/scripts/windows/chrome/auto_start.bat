@echo off
:: Copyright (c) 2018-2028 Dunyu All Rights Reserved.
::
:: Author      : https://www.wengold.net
:: Email       : support@wengold.net
::
:: Prismy.No | Date       | Modified by. | Description
:: -------------------------------------------------------------------
:: 00001       2020/05/08   yangping       New version
:: 00002       2020/08/16   yangping       Support for windows
:: -------------------------------------------------------------------

set BINPATH=%~dp0
call %BINPATH%\exports.bat

cd /d %BINPATH%\..\..
call %SERVICE_APP_NAME%.exe
echo started server...


:: waiting for 5 seconds later
ping -n 5 127.1>nul

:: start browser and fullscreen to show dashbord page
:: start chrome.exe --start-fullscreen "http://%SERVICE_HOST_PORT%/%SERVICE_APP_NAME%/"
start chrome.exe --kiosk "http://%SERVICE_HOST_PORT%/%SERVICE_APP_NAME%/"

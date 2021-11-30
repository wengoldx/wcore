@echo off

:: Copyright (c) 2019-2029 Dunyu All Rights Reserved.
::
:: Author      : yangping
:: Email       : ping.yang@wengold.net
:: Version     : 1.0.1
:: Description :
::   Start server.
::
:: Prismy.No | Date       | Modified by. | Description
:: -------------------------------------------------------------------
:: 00001       2021/08/29   yangping       New version
:: -------------------------------------------------------------------

set BINPATH=%~dp0
call %BINPATH%\scripts\exports.bat

cd /d %BINPATH%\..
call %SERVICE_APP_NAME%.exe
echo started server...

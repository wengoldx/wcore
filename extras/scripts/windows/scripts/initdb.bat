@echo off

:: Copyright (c) 2019-2029 Dunyu All Rights Reserved.
::
:: Author      : yangping
:: Email       : ping.yang@wengold.net
:: Version     : 1.0.1
:: Description :
::   Create database for server.
::
:: Prismy.No | Date       | Modified by. | Description
:: -------------------------------------------------------------------
:: 00001       2021/08/29   yangping       New version
:: -------------------------------------------------------------------

set BINPATH=%~dp0
call %BINPATH%\exports.bat

:: check database user and password input data
if "%SERVICE_DB_USER%"=="" goto ERROR

:: init database from .sql file
mysql -u%SERVICE_DB_USER% -p < %BINPATH%\db_create.sql --default-character-set=utf8mb4
echo Inited database %SERVICE_DATABASE% for mysql user : %SERVICE_DB_USER%

:: upgrade database to version-1
echo Next to upgrade database to version-1
mysql -u%SERVICE_DB_USER% -p < %BINPATH%\db_upgrade_v1.sql --default-character-set=utf8mb4
goto END

:: print out error or end messages
:ERROR
echo "Invalid database user or password!"
pause
exit 0

:: exit script
:END
exit 0
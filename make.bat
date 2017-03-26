@ECHO off
IF "%GOROOT%"=="" (
	ECHO GOHOME not defined
	GOTO exit
)

IF NOT EXIST "%GOROOT%bin\go.exe" (
	ECHO GO not found
	GOTO exit
) 

SET GO="%GOROOT%bin\go.exe"

ECHO Generating SparrowDb binaries
IF EXIST "dist" (
	RMDIR /s /q dist
)
MKDIR dist

ECHO Building ...
%GO% build -o dist/sparrow.exe .
%GO% build -o dist/commander.exe tools/commander/commander.go
%GO% build -o dist/datafile.exe tools/datafile/datafile.go

ECHO Copying ...
ROBOCOPY scripts dist/scripts
ROBOCOPY config dist/config
XCOPY /f README.md dist/
XCOPY /f LICENCE dist/

:exit
	ECHO Done
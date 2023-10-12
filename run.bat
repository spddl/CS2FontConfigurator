@echo off

set filename=CS2FontConfigurator

:loop
cls

rem gocritic check -enableAll -disable="#experimental,#opinionated,#commentedOutCode" ./...
go build -tags debug -buildvcs=false -o %filename%.exe

IF %ERRORLEVEL% EQU 0 %filename%.exe

pause
goto loop
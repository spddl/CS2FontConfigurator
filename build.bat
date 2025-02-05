@echo off

SET GOOS=windows
SET GOARCH=amd64
set filename=CS2FontConfigurator

:loop
CLS

gocritic check -enableAll -disable="#experimental,#opinionated,#commentedOutCode" ./...

IF exist %filename%.exe (
    FOR /F "usebackq" %%A IN ('%filename%.exe') DO SET /A beforeSize=%%~zA
) ELSE (
    SET /A beforeSize=0
)

rem rsrc.exe -manifest CS2FontConfigurator.exe.manifest -ico icon.ico

:: Build https://golang.org/cmd/go/
go build -buildvcs=false -ldflags="-w -s -H windowsgui" -o %filename%.exe
go build -buildvcs=false -o %filename%_debug.exe

FOR /F "usebackq" %%A IN ('%filename%.exe') DO SET /A size=%%~zA
SET /A diffSize = %size% - %beforeSize%
SET /A size=(%size%/1024)+1
IF %diffSize% EQU 0 (
    ECHO %size% kb
) ELSE (
    IF %diffSize% GTR 0 (
        ECHO %size% kb [+%diffSize% b]
    ) ELSE (
        ECHO %size% kb [%diffSize% b]
    )
)

:: Run
IF %ERRORLEVEL% EQU 0 %filename%.exe

PAUSE
GOTO loop
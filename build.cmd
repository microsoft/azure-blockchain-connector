@echo off

set GOARCH=amd64
cd .\cmd\abc

start /w /b go build -o ..\..\build\abc.exe

set CGO_ENABLED=0
set GOOS=darwin
start /w /b go build -o ..\..\build\abc_darwin

set GOOS=linux
start /w /b go build -o ..\..\build\abc_linux

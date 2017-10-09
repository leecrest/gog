echo off

echo build gate
cd src\gate
go build
copy gate.exe ..\..\bin\
del /f gate.exe

echo build game
cd ..\game
go build
copy game.exe ..\..\bin\
del /f game.exe

echo build account
cd ..\account
go build
copy account.exe ..\..\bin\
del /f account.exe

pause
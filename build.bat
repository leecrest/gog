echo build gateserver
cd src\gateserver
go build
copy gateserver.exe ..\bin\
del /f gateserver.exe

cd ..\src\gameserver
go build
copy gameserver.exe ..\bin\
del /f gameserver.exe

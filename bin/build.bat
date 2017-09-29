cd ..\server\gateserver
go build

cd ..\..\
copy server\gateserver\gateserver.exe bin\
del /f server\gateserver\gateserver.exe

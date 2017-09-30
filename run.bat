cd bin

echo start gateserver
gateserver.exe -cfg=setting\server.ini -id=0

echo start gameserver0
gameserver.exe -cfg=setting\server.ini -id=0
cd bin

echo start gate
start gate.exe -cfg=setting\server.ini -id=0

echo start account
start account.exe -cfg=setting\server.ini -id=0 -name=account

echo start game0
start game.exe -cfg=setting\server.ini -id=1 -name=user
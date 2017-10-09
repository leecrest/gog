
echo build gate
cd src/gate
go build
cp gate ../../bin/
rm -f gate

echo build game
cd ../game
go build
cp game ../../bin/
rm -f game

echo build account
cd ../account
go build
cp account ../../bin/
rm -f account

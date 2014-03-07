
#!/bin/sh

CURRENT=`pwd`
BASENAME=`basename $CURRENT`

echo $CURRENT
echo $GOPATH
echo $PATH

cd $CURRENT/src/gate
go get
go install

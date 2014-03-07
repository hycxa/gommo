
#!/bin/sh

CURRENT=`pwd`
BASENAME=`basename $CURRENT`

cd $CURRENT/src/gate
go get
go install

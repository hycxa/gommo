
#!/bin/sh

BASENAME=`basename $PWD`

echo "PWD" $PWD
echo "GOPATH" $GOPATH
echo "PATH" $PATH
echo "BUILD_NUMBER" $BUILD_NUMBER
echo "BUILD_ID" $BUILD_ID

cd $PWD/src/gate
go get
go install

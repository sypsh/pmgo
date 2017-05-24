#!/bin/sh

version="0.5.0"

UNAME=$(uname)
# Check to see if it starts with MINGW.
if [ "$UNAME" ">" "MINGW" -a "$UNAME" "<" "MINGX" ] ; then
    echo "current release do not support windows"
    exit 1
fi
if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
    echo "Sorry, this OS is not supported yet via this installer."
    exit 1
fi


wget https://github.com/struCoder/pmgo/archive/v${version}.tar.gz
tar -zxvf v${version}.tar.gz
cd pmgo-${version}

echo "build pmgo..."
go build -o /usr/local/bin/pmgo pmgo.go

echo "build done"

cd ..
echo "rm build files"

rm -rf pmgo-${version}
rm v${version}.tar.gz

echo "all down."

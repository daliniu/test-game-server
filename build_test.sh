#!/bin/sh

git pull

if [ $1 == "all" ]; then
    for srv in `ls src`; do
        if [ ${srv:0-4} == "-srv" ]; then
            echo "开始更新$srv"
            GOROOT=/data/services/go GOBIN=/data/goapp/wolf/bin GOPATH=`pwd`:`pwd`/vendor /data/services/go/bin/go install $srv && sudo supervisorctl restart wolf-$srv:*
        fi
    done
else
    for srv in "$@"
    do
        if [ "${srv:0-4}" != "-srv" ]; then
            srv="${srv}-srv"
        fi
        echo "开始更新$srv"
        GOROOT=/data/services/go GOBIN=/data/goapp/wolf/bin GOPATH=`pwd`:`pwd`/vendor /data/services/go/bin/go install $srv && sudo supervisorctl restart wolf-$srv:*
    done
fi

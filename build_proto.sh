#!/bin/sh
CURRENTDIR=`pwd`
protoc --proto_path=$CURRENTDIR/src/share/proto --go_out=plugins=micro:$CURRENTDIR/src/share/pb $CURRENTDIR/src/share/proto/*.proto
protoc --proto_path=$CURRENTDIR/src/share/proto --js_out=/$CURRENTDIR/src/share/jspb $CURRENTDIR/src/share/proto/*.proto

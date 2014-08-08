#!/bin/bash

set -e

mkdir -p android/libs/armeabi-v7a
mkdir -p android/obj/local/armeabi-v7a
CC="$NDK_ROOT/bin/arm-linux-androideabi-gcc"

if [ -z $GOANDROID ]
then
  GOANDROID="../go"
fi

CC=$CC GOPATH="`pwd`:$GOPATH" GOROOT="" GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 $GOANDROID/bin/go install $GOFLAGS -v -ldflags="-android -shared -extld $CC -extldflags '-march=armv7-a -mfloat-abi=softfp -mfpu=vfpv3-d16'" -tags android nativeactivity
cp bin/linux_arm/nativeactivity android/libs/armeabi-v7a/libnativeactivity.so
cp bin/linux_arm/nativeactivity android/obj/local/armeabi-v7a/libnativeactivity.so

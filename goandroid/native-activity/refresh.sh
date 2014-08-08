#!/bin/bash

APKNAME=nativeactivity

cd android
ln -sT libs lib
zip -u bin/$APKNAME-debug.apk lib/armeabi-v7a/lib$APKNAME.so
rm lib
cd ..
./signapk.sh "android/bin/$APKNAME-debug.apk"
mv signed_$APKNAME-debug.apk "android/bin/$APKNAME-debug.apk"

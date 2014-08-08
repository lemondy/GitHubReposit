goandroid
=========

Patches to the Go tools and runtime to enable Android apps to interface directly with a shared library written in Go. Goandroid also includes two demos, `hellogl2` and `nativeactivity`, both written in Go. Goandroid currently runs on ARMv7 or later CPUs.

Running [Go](http://golang.org) code from Android apps is currently not possible, because the Go tools can only output executables while Android requires any foreign code in shared library (.so) format. This repository contains patches for the Go tools and runtime to enable shared library output, including workarounds to Android specific limitations.

Note that only Go 1.2 is supported by goandroid. Go 1.3 will probably never be supported now that official Android su

*Disclaimer*: Please note that the patches are unofficial and only works for Go 1.2. Goandroid will probably never be updated for Go 1.3 now that official support is being worked on for [Go 1.4](https://groups.google.com/forum/#!topic/golang-dev/P1ATVp1mun0). Also note that goandroid only includes the bare minimum to let Go run in Android apps. If you need a more comprehensive framework, check out [Mandala](https://github.com/remogatto/mandala).

### Set up ###

This guide is tested on linux/amd64 and assumes you have an android device connected through USB and that you meet the requirements for [building Go from source](http://golang.org/doc/install/source).

1. Download and install the [NDK](http://developer.android.com/tools/sdk/ndk/index.html) at least version r8e. These instructions assumes the NDK is installed in `$NDK`.
2. Create a standalone NDK toolchain (as described in $NDK/docs/STANDALONE-TOOLCHAIN.html):

	`$NDK/build/tools/make-standalone-toolchain.sh --platform=android-9 --install-dir=ndk-toolchain`

	You might need to add `--system=linux-x86_64` or `--system=darwin-x86_64` depending on your system.

	Set `$NDK_ROOT` to point at the `ndk-toolchain` directory

3. Clone the golang repository:

	`hg clone -u go1.2.2 https://code.google.com/p/go`

4. Copy the `patches` directory  to the `go/.hg` directory:

	`cp -a patches go/.hg`

5. Enable the `mq` extension by adding the following lines to `go/.hg/hgrc`:

	```
	[extensions]  
	mq =
	codereview = !

	[ui]  
	username = me<me@mail.com>
	```

6. In the `go/src` directory apply the patches and build go:

	```
	cd go/src  
	hg qpush -a  
	CGO_ENABLED=0 GOOS=linux GOARCH=arm ./make.bash \
	  CC="$NDK_ROOT/bin/arm-linux-androideabi-gcc" GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 ../bin/go install -tags android -a -v std  
	cd ../..
	```

### Building and installing the example apps ###

If everything is set up correctly, you should be able to `cd hello-gl2` and run `build.sh` to build and copy `libgoandroid.so` to android/libs. Then, running `ant -f android/build.xml clean debug install` will build and install the final apk to the connected device. Running the app should display a simple color animated triangle that you can move around the screen with your finger.

A more complicated but also more useful example is `nativeactivity`. It mimicks the C/C++ `native_app_glue` library and uses the NativeActivity API to completely avoid Java code and gain control of the input and render loop as well as context creation through EGL. It can be compiled and installed in the same way as `hellogl2` with `./build.sh` and `ant -f android/build.xml clean debug install`. The nativeactivity sample requires Android 2.3, which is the version where NativeActivity were introduced.

An interesting artifact of Go apps is that the compile-deploy cycle can be shorter than the ant scripts. If an existing debug apk is already present after an `ant -f android/build.xml clean debug install` you can use `./upload.sh` to replace the Go library and upload the apk to the device. On my system, `touch src/nativeactivity/main.go android/AndroidManifest.xml && time ./build.sh && time ant -f android/build.xml debug install` takes 17 seconds, while `touch src/nativeactivity/main.go && time ./upload.sh` takes 10 seconds. This difference will only be more exaggerated if the apk contains resources.

### Go patches ###

All patches except `android-tls` and `android-build-hacks` correspond to the patches for linux/arm external linking and shared library support discussed on the [golang-nuts mailing list](https://groups.google.com/d/msg/golang-nuts/zmjXkGrEx6Q/L4R8qyw7WW4J).

The `android-tls` patch is a workaround for the missing support for the `R_ARM_TLS_IE32` relocation in the Android linker.

The `android-build-hacks` patch contains various changes to account for the difference between a vanilla linux/arm system and Android.

package net.goandroid.hellogl2;

public class Engine {
	static native void resize(int w, int h);
	static native void init();
	static native void drawFrame();
	static native void onTouch(int action, float x, float y);

	static {
		System.loadLibrary("hellogl2");
	}
}

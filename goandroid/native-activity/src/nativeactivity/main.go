// +build android

package main

// #include <stdlib.h>
// #include <jni.h>
// #include <android/native_activity.h>
// #include <android/input.h>
// #include <EGL/egl.h>
// #include "main.h"
//
// #cgo LDFLAGS: -landroid -lEGL
import "C"

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"unsafe"
)

type mainLoop struct {
	looper *C.ALooper

	quit   chan struct{}
	resume chan struct{}
	pause  chan struct{}
	focus  chan bool
	render chan *renderState
	input  chan *C.AInputQueue
	ack    chan struct{}

	inputQ           *C.AInputQueue
	renderState      *renderState
	running, focused bool
	width, height    int

	game *game
}

type activityState struct {
	renderState *renderState
	mLoop       *mainLoop
}

var states map[*C.ANativeActivity]*activityState = make(map[*C.ANativeActivity]*activityState)

const (
	LOOPER_ID_INPUT = iota
)

func (s *activityState) Destroy() {
	s.mLoop.Quit()
	s.renderState.Destroy()
}

type renderState struct {
	disp C.EGLDisplay
	conf C.EGLConfig
	ctx  C.EGLContext
	surf C.EGLSurface
}

func (s *renderState) Destroy() {
	if s == nil {
		return
	}
	if s.disp != nil {
		if s.ctx != nil {
			C.eglDestroyContext(s.disp, s.ctx)
		}
		C.eglTerminate(s.disp)
	}
}

func newMainLoop() (m *mainLoop) {
	m = &mainLoop{
		quit:   make(chan struct{}, 1),
		resume: make(chan struct{}, 1),
		pause:  make(chan struct{}, 1),
		focus:  make(chan bool, 1),
		render: make(chan *renderState, 1),
		input:  make(chan *C.AInputQueue, 1),
		ack:    make(chan struct{}, 1),
		game:   &game{},
	}
	init := make(chan struct{})
	go m.loop(init)
	<-init
	return m
}

func (m *mainLoop) Resume() {
	m.resume <- struct{}{}
	m.wakeAndAck()
}

func (m *mainLoop) Focused(focused bool) {
	m.focus <- focused
	m.wakeAndAck()
}

func (m *mainLoop) Pause() {
	m.pause <- struct{}{}
	m.wakeAndAck()
}

func (m *mainLoop) Quit() {
	m.quit <- struct{}{}
	m.wakeAndAck()
}

func (m *mainLoop) UpdateRenderState(rs *renderState) {
	m.render <- rs
	m.wakeAndAck()
}

func (m *mainLoop) isRunning() bool {
	return m.focused && m.running && m.renderState != nil
}

func (m *mainLoop) loop(init chan<- struct{}) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	looper := C.ALooper_prepare(C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	if looper == nil {
		panic("ALooper_prepare returned nil")
	}
	m.looper = looper
	init <- struct{}{}
	for {
		select {
		case <-m.quit:
			if m.renderState != nil && m.renderState.ctx != nil {
				if C.eglMakeCurrent(m.renderState.disp, C.EGLSurface(unsafe.Pointer(nil)), C.EGLSurface(unsafe.Pointer(nil)), C.EGLContext(unsafe.Pointer(nil))) != C.EGL_TRUE {
					panic("Error: eglMakeCurrent() failed\n")
				}
			}
			m.ack <- struct{}{}
			break
		case <-m.resume:
			m.running = true
			m.ack <- struct{}{}
		case <-m.pause:
			m.running = false
			m.width, m.height = 0, 0
			m.ack <- struct{}{}
		case m.focused = <-m.focus:
			m.width, m.height = 0, 0
			m.ack <- struct{}{}
		case m.renderState = <-m.render:
			m.ack <- struct{}{}
		case inputQ := <-m.input:
			if inputQ != nil {
				C.AInputQueue_attachLooper(inputQ, m.looper, LOOPER_ID_INPUT, nil, nil)
			} else {
				C.AInputQueue_detachLooper(m.inputQ)
			}
			m.inputQ = inputQ
			m.ack <- struct{}{}
		default:
			m.frame()
		}
	}
}

func (m *mainLoop) frame() {
	var timeout C.int = 0
	if !m.isRunning() {
		timeout = -1
	}
	ident := C.ALooper_pollAll(timeout, nil, nil, nil)
	switch ident {
	case LOOPER_ID_INPUT:
		if m.inputQ != nil {
			m.processInput(m.inputQ)
		}
	case C.ALOOPER_POLL_ERROR:
		log.Fatalf("ALooper_pollAll returned ALOOPER_POLL_ERROR\n")
	}
	if m.isRunning() {
		m.checkSize()
		createCtx := m.renderState.ctx == nil
		if createCtx {
			log.Printf("Creating context\n")
			ctx_attribs := [...]C.EGLint{
				C.EGL_CONTEXT_CLIENT_VERSION, 2,
				C.EGL_NONE,
			}

			m.renderState.ctx = C.eglCreateContext(m.renderState.disp, m.renderState.conf, C.EGLContext(unsafe.Pointer(nil)), (*C.EGLint)(unsafe.Pointer(&ctx_attribs[0])))
			if m.renderState.ctx == nil {
				panic("Error: eglCreateContext failed\n")
			}
		}

		if C.eglMakeCurrent(m.renderState.disp, m.renderState.surf, m.renderState.surf, m.renderState.ctx) != C.EGL_TRUE {
			panic("Error: eglMakeCurrent() failed\n")
		}
		if createCtx {
			m.game.initGL()
		}
		m.game.drawFrame()
		C.eglSwapBuffers(m.renderState.disp, m.renderState.surf)
	}
}

func (m *mainLoop) checkSize() {
	var w, h C.EGLint
	C.eglQuerySurface(m.renderState.disp, m.renderState.surf, C.EGL_WIDTH, &w)
	C.eglQuerySurface(m.renderState.disp, m.renderState.surf, C.EGL_HEIGHT, &h)
	width, height := int(w), int(h)
	if width != m.width || height != m.height {
		m.width = width
		m.height = height
		m.game.resize(m.width, m.height)
	}
}

func (m *mainLoop) inputQueue(inputQ *C.AInputQueue) {
	m.input <- inputQ
	m.wakeAndAck()
}

func (m *mainLoop) wakeAndAck() {
	C.ALooper_wake(m.looper)
	<-m.ack
}

func handleCallbackError(act *C.ANativeActivity, err interface{}) {
	if err == nil {
		return
	}
	errStr := fmt.Sprintf("callback panic: %s stack: %s", err, debug.Stack())
	errStrC := C.CString(errStr)
	defer C.free(unsafe.Pointer(errStrC))
	if C.throwException(act, errStrC) == 0 {
		log.Fatalf("%v\n", errStr)
	}
}

//export onWindowFocusChanged
func onWindowFocusChanged(act *C.ANativeActivity, focusedC C.int) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onWindowFocusChanged %v...\n", focusedC)
	focused := false
	if focusedC != 0 {
		focused = true
	}
	states[act].mLoop.Focused(focused)
	log.Printf("onWindowFocusChanged done\n")
}

//export onConfigurationChanged
func onConfigurationChanged(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onConfigurationChanged\n")
}

//export onNativeWindowResized
func onNativeWindowResized(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onNativeWindowResized\n")
}

//export onInputQueueDestroyed
func onInputQueueDestroyed(act *C.ANativeActivity, queue unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onInputQueueDestroy...\n")
	states[act].mLoop.inputQueue(nil)
	log.Printf("onInputQueueDestroy done\n")
}

func (m *mainLoop) dispatchEvent(event *C.AInputEvent) bool {
	switch C.AInputEvent_getType(event) {
	case C.AINPUT_EVENT_TYPE_MOTION:
		return m.game.onTouch(event)
	}
	return false
}

func (m *mainLoop) processInput(inputQueue *C.AInputQueue) {
	var event *C.AInputEvent
	for {
		if ret := C.AInputQueue_getEvent(inputQueue, &event); ret < 0 {
			break
		}
		if C.AInputQueue_preDispatchEvent(inputQueue, event) != 0 {
			continue
		}
		handled := m.dispatchEvent(event)
		var handledC C.int
		if handled {
			handledC = 1
		}
		C.AInputQueue_finishEvent(inputQueue, event, handledC)
	}
}

//export onInputQueueCreated
func onInputQueueCreated(act *C.ANativeActivity, queue unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onInputQueueCreated...\n")
	inputQ := (*C.AInputQueue)(queue)
	state := states[act]
	state.mLoop.inputQueue(inputQ)
	log.Printf("onInputQueueCreated done\n")
}

//export onPause
func onPause(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("Pausing...\n")
	states[act].mLoop.Pause()
	log.Printf("Paused...\n")
}

//export onResume
func onResume(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("Resuming...\n")
	states[act].mLoop.Resume()
	log.Printf("Resumed...\n")
}

//export onCreate
func onCreate(act *C.ANativeActivity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onCreate...\n")
	state := &activityState{
		mLoop: newMainLoop(),
	}
	states[act] = state
	log.Printf("onCreate done\n")
}

//export onDestroy
func onDestroy(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onDestroy...\n")
	state := states[act]
	delete(states, act)
	state.Destroy()
	log.Printf("onDestroy done\n")
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onWindowDestroy...\n")
	state := states[act]
	state.mLoop.UpdateRenderState(nil)
	C.eglDestroySurface(state.renderState.disp, state.renderState.surf)
	state.renderState.surf = nil
	log.Printf("onWindowDestroy done\n")
}

func getEGLDisp(disp C.EGLNativeDisplayType) C.EGLDisplay {
	if C.eglBindAPI(C.EGL_OPENGL_ES_API) == C.EGL_FALSE {
		panic("Error: eglBindAPI() failed")
	}

	egl_dpy := C.eglGetDisplay((C.EGLNativeDisplayType)(disp))
	if egl_dpy == nil {
		panic("Error: eglGetDisplay() failed\n")
	}

	var egl_major, egl_minor C.EGLint
	if C.eglInitialize(egl_dpy, &egl_major, &egl_minor) != C.EGL_TRUE {
		panic("Error: eglInitialize() failed\n")
	}
	return egl_dpy
}

func EGLCreateWindowSurface(eglDisp C.EGLDisplay, config C.EGLConfig, win C.EGLNativeWindowType) C.EGLSurface {
	eglSurf := C.eglCreateWindowSurface(eglDisp, config, win, (*C.EGLint)(unsafe.Pointer(nil)))
	if eglSurf == nil {
		panic("Error: eglCreateWindowSurface failed\n")
	}
	return eglSurf
}

func getEGLNativeVisualId(eglDisp C.EGLDisplay, config C.EGLConfig) C.EGLint {
	var vid C.EGLint
	if C.eglGetConfigAttrib(eglDisp, config, C.EGL_NATIVE_VISUAL_ID, &vid) != C.EGL_TRUE {
		panic("Error: eglGetConfigAttrib() failed\n")
	}
	return vid
}

func chooseEGLConfig(eglDisp C.EGLDisplay) C.EGLConfig {
	eglAttribs := [...]C.EGLint{
		C.EGL_RED_SIZE, 4,
		C.EGL_GREEN_SIZE, 4,
		C.EGL_BLUE_SIZE, 4,
		//C.EGL_DEPTH_SIZE, 1,
		C.EGL_RENDERABLE_TYPE, C.EGL_OPENGL_ES2_BIT,
		C.EGL_SURFACE_TYPE, C.EGL_WINDOW_BIT,
		C.EGL_NONE,
	}

	var config C.EGLConfig
	var num_configs C.EGLint
	if C.eglChooseConfig(eglDisp, (*C.EGLint)(unsafe.Pointer(&eglAttribs[0])), &config, 1, &num_configs) != C.EGL_TRUE {
		panic("Error: couldn't get an EGL visual config\n")
	}

	return config
}

//export onNativeWindowCreated
func onNativeWindowCreated(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onNativeWindowCreated...\n")
	state := states[act]
	if state.renderState == nil {
		state.renderState = &renderState{
			disp: getEGLDisp(C.EGLNativeDisplayType(nil)),
		}
		state.renderState.conf = chooseEGLConfig(state.renderState.disp)
	}
	vid := getEGLNativeVisualId(state.renderState.disp, state.renderState.conf)
	C.ANativeWindow_setBuffersGeometry((*[0]byte)(win), 0, 0, C.int32_t(vid))
	state.renderState.surf = EGLCreateWindowSurface(state.renderState.disp, state.renderState.conf, C.EGLNativeWindowType(win))

	state.mLoop.UpdateRenderState(state.renderState)
	log.Printf("onNativeWindowCreated done\n")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

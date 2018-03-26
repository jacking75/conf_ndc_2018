package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	mWaitGroupGoroutine sync.WaitGroup
	// server close signal
	mServerTerminate = make(chan struct{})
)

func IncrementWait() {
	mWaitGroupGoroutine.Add(1)
}

func DecrementWait() {
	mWaitGroupGoroutine.Done()
}

func GetServerTerminateChennel() chan struct{} {
	return mServerTerminate
}

// handle unix signals
func Sighandler(mainloop chan<- struct{}) {
	defer PrintPanicStack()
	//utils.Logger.Info("Exit sighandler")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	for {
		msg := <-ch

		switch msg {
		case syscall.SIGTERM: // os 명령어 kill로 종료 시켰음
			Logger.Info("sigterm received: syscall.SIGTERM")
			sighandlerProcessExit(mainloop)
		case syscall.SIGINT: // ctrl + c 로 종료 시켰음
			Logger.Info("sigterm received: syscall.SIGINT")
			sighandlerProcessExit(mainloop)
		}

		Logger.Info("Exit sighandler")
		close(ch)
		return
	}

}

func sighandlerProcessExit(mainloop chan<- struct{}) {
	close(mainloop)

	close(mServerTerminate)

	Logger.Info("waiting for session close, please wait...")
	mWaitGroupGoroutine.Wait()
}

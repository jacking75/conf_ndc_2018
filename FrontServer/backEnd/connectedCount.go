package backEnd

import "sync/atomic"

var mConnectedBackendCount int32

func IncrementConnectCount() int32 {
	newValue := atomic.AddInt32(&mConnectedBackendCount, 1)
	return newValue
}

func DecrementConnectCount() int32 {
	newValue := atomic.AddInt32(&mConnectedBackendCount, -1)
	return newValue
}

func CurrentConnectCount() int32 {
	value := atomic.LoadInt32(&mConnectedBackendCount)
	return value
}

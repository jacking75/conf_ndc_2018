package monitoring

import (
	"sync/atomic"

	"../utils"
)

type ServerStatus struct {
	ChannelCount int32
}

var serverStatus ServerStatus

//TODO: utils.IncrementWait() 분리하기
func ServerStatusCreateChannel() int32 {
	newValue := atomic.AddInt32(&serverStatus.ChannelCount, 1)
	utils.IncrementWait()
	return newValue
}

//TODO: utils.DecrementWait() 분리하기
func ServerStatusDestoryChannel() int32 {
	newValue := atomic.AddInt32(&serverStatus.ChannelCount, -1)
	utils.DecrementWait()
	return newValue
}

func ServerStatusChannelCount() int32 {
	value := atomic.LoadInt32(&serverStatus.ChannelCount)
	return value
}

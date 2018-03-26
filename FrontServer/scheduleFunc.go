package main

import (
	"fmt"
	"time"

	"./monitoring"
	"./utils"

	"github.com/carlescere/scheduler"
	"github.com/emirpasic/gods/lists/arraylist"
)

var mShedulList *arraylist.List

func schedulStart() {
	mShedulList = arraylist.New()

	job1, _ := scheduler.Every(3).Seconds().Run(jobPrintServerStatus) //scheduler.Every(0.1).Minutes().Run(jobPrintServerStatus)
	mShedulList.Add(job1)
	utils.Logger.Info("schedul Count: ", mShedulList.Size())
}

func schedulEnd() {
	utils.Logger.Info("schedul Count: ", mShedulList.Size())

	mShedulList.Each(func(index int, value interface{}) {
		job, _ := mShedulList.Get(index)
		job.(*scheduler.Job).Quit <- true
		utils.Logger.Info("job Quit")
	})

	time.Sleep(1 * time.Second)
	mShedulList.Clear()

	utils.Logger.Info("schedulEnd")
}

func jobPrintServerStatus() {
	channelCount := monitoring.ServerStatusChannelCount()
	fmt.Println("Current Channel Count", channelCount)
}

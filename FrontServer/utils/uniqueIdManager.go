package utils

import (
	"github.com/bwmarrin/snowflake"
)

var generator *snowflake.Node

func UniqueIdManagerInit(machineID int64) {
	NewGenerator, err := snowflake.NewNode(machineID)

	if err != nil {
		Logger.Fatal("sonyflake not created: ", err)
	}

	generator = NewGenerator
}

func UniqueIdManagerNewId() int64 {
	id := generator.Generate()
	return id.Int64()
}

package utils

import (
	"runtime"

	"github.com/davecgh/go-spew/spew"
)

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		Logger.Error(x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			Logger.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			Logger.Errorf("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
		}
	}
}

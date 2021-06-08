package errorz

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Error struct {
	E error
	StackTraces []string
}


// Pack return a format error
// example : " 2021-05-29 15:33:06 | /User/zonst/path/controller.go:20 | error is nil \n "
func Pack(e error) error {
	if e == nil {
		return nil
	}

	//获取正在运行的go的数据
	_, file, line, _ := runtime.Caller(1)
	format := PrintFormat(file, line, e.Error())
	//errorZ.StackTraces = append(errorZ.StackTraces,format)

	return errors.New(format)
}

func NewString(str string) error {


	return nil
}

func Empty() Error {
	return Error{
		E:           nil,
		StackTraces: make([]string, 0, 10),
	}
}

func PrintFormat(file string,line int, eStr string) string {
	var formatGroup = make([]string,0,3)

	formatGroup = append(formatGroup,time.Now().Format("2006-01-02 15:04:05"))

	trac := fmt.Sprintf("%s:%d",file,line)
	formatGroup = append(formatGroup,trac)

	formatGroup = append(formatGroup,eStr)

	return fmt.Sprintf(strings.Join(formatGroup," | ")+"\n")
}
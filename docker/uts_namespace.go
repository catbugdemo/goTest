package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// uts_namespace 目的，用来隔离 nodename 和 domainname
// 每个namespace允许有自己的hostname

type Syscall struct {
	syscall.SysProcAttr
	Cloneflags uintptr
}

func main() {
	//指定被 fork 出来新进程的初始命令, 默认使用sh
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{

	}
	//指向标准输入
	cmd.Stdin  = os.Stdin
	//指向标准输出
	cmd.Stdout = os.Stdout
	//指向标准错误
	cmd.Stderr = os.Stderr

	//Run执行c包含的命令，并阻塞直到完成。命令成功表示stdin,stdout,stderr转交没有错误
	//err一般表示I/O出现错误
	if e := cmd.Run();e!= nil{
		log.Fatal(e)
	}
}




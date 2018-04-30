/********************************************************************************
* log.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package app

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	. "github.com/Azraid/pasque/core"
)

type logWriteCloser struct {
	path   string
	prefix string
	w      *bufio.Writer
	fp     *os.File
}

func (lwc *logWriteCloser) Write(p []byte) (int, error) {

	if Config.Global.UseStdOut {
		fmt.Printf("[%s]%s\r\n", lwc.prefix, string(p))
	}

	if err := os.MkdirAll(lwc.path, 0777); err != nil {
		fmt.Println("Can not create log directory", err)
		return 0, err
	}

	fn := fmt.Sprintf("%s/%s.%s.%s.%s", lwc.path, App.Hostname, App.Eid, time.Now().Format("20060102.15"), lwc.prefix)

	if lwc.w == nil || lwc.fp == nil || lwc.fp.Name() != fn {
		if nfp, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err == nil {
			if lwc.w != nil {
				lwc.w.Flush()
			}

			if lwc.fp != nil {
				lwc.fp.Close()
			}

			lwc.fp = nfp
			lwc.w = bufio.NewWriter(lwc.fp)
		} else {
			return 0, IssueErrorf("can not open log file%s", fn)
		}
	}

	n, err := lwc.w.Write(p)
	if Config.Global.Log.Immediate {
		lwc.w.Flush()
	}
	return n, err
}

//여기서 프로그램이 종료될때, close가 불리지 않는다면.
// SetFinalizer()로 flush 시켜야 한다.
// https://golang.org/pkg/runtime/
func (lwc *logWriteCloser) Close() error {
	if lwc.fp != nil {
		lwc.w.Flush()
		return lwc.fp.Close()
	}

	return nil
}

var elog *log.Logger
var ilog *log.Logger
var dlog *log.Logger
var dump *log.Logger
var plog *log.Logger

var logDConn *net.UDPConn

func connectLogD() error {
	if len(Config.Global.LogDAddr) == 0 {
		return IssueErrorf("logDAddr is wrong")
	}

	addr, err := net.ResolveUDPAddr("udp", Config.Global.LogDAddr)
	if err != nil {
		return err
	}

	logDConn, err = net.DialUDP("udp", nil, addr)

	if err != nil {
		return err
	}

	return nil
}

func initLog(path string) {
	if len(path) == 0 {
		path, _ = os.Getwd()
		path += "/log"
	}

	elog = log.New(&logWriteCloser{path: path, prefix: "error"}, "", log.LstdFlags|log.Lmicroseconds)
	ilog = log.New(&logWriteCloser{path: path, prefix: "info"}, "", log.LstdFlags|log.Lmicroseconds)
	dlog = log.New(&logWriteCloser{path: path, prefix: "debug"}, "", log.LstdFlags|log.Lmicroseconds)
	dump = log.New(&logWriteCloser{path: path, prefix: "dmp"}, "", log.Ldate|log.Ltime)
	plog = log.New(&logWriteCloser{path: path, prefix: "packet"}, "", log.LstdFlags|log.Lmicroseconds)

	connectLogD()
}

//ErrorLog 는 오류에 대해서 기록한다.
func ErrorLog(a string, v ...interface{}) {
	if Config.Global.Log.Error {
		elog.Output(2, fmt.Sprintf(a, v...))
	}
}

//InfoLog BIZ 정보에 대한 것을 기록한다.
func InfoLog(a string, v ...interface{}) {
	if Config.Global.Log.Info {
		ilog.Output(2, fmt.Sprintf(a, v...))
	}
}

//DebugLog calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func DebugLog(a string, v ...interface{}) {
	if Config.Global.Log.Debug {
		dlog.Output(2, fmt.Sprintf(a, v...))
	}
}

func PacketLog(a string, v ...interface{}) {
	if Config.Global.Log.Packet {
		plog.Output(2, fmt.Sprintf(a, v...))
	}

	if logDConn != nil {
		msg := "[" + App.Eid + "]" + time.Now().Format("2006/01/02 15:30:30.000 ") + fmt.Sprintf(a, v...) + "\r\n"
		logDConn.Write([]byte(msg))
	}
}

//Dump is log for all goroutine stack
func Dump(r interface{}) {
	if r != nil {
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		dump.Output(3, fmt.Sprintf("%v\r\n%s", r, buf[0:stackSize]))

	}
}

func DumpRecover() {
	if Config.Global.DumpRecover {
		if r := recover(); r != nil {
			Dump(r)
		}
	}
}

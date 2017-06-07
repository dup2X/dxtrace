package dxtrace

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

func init() {
	logFile, _ := os.OpenFile(tmpPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 2)
	go func() {
		tk := time.NewTicker(time.Second * 1)
		for range tk.C {
			data, err := ioutil.ReadFile(tmpPath)
			if err == nil && len(data) != 0 {
				datas := bytes.Split(data, []byte("SCHED"))
				for i := range datas {
					if len(datas[i]) == 0 {
						continue
					}
					var head = []byte("SCHED")
					report(append(head, datas[i]...))
				}
			}
			os.Truncate(tmpPath, 0)
		}
	}()
}

var tmpPath = "/tmp/xdebug"

func report(data []byte) {
	parser := &textParser{}
	r, err := parser.Parse(data)
	if err != nil {
		fmt.Printf("err=%s\n", err)
		return
	}
	fmt.Printf("%+v\n", r)
}

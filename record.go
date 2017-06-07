package dxtrace

import (
	"bytes"
	"fmt"
	"strconv"
)

type record struct {
	cost            string
	maxProc         int64
	idleprocs       int64
	threads         int64
	spinningthreads int64
	idlethreads     int64
	runqueue        int64
	gcwaiting       int64
	nmidlelocked    int64
	stopwait        int64
	sysmonwait      int64
	allp            []p
	allm            []m
	allg            []g
}

type p struct {
	status      int64
	schedtick   int64
	syscalltick int64
	mid         int64
	runqsize    int64
	gfreecnt    int64
}

type m struct {
	id         int64
	pid        int64
	curg       int64
	mallocing  int64
	throwing   int64
	preemptoff bool
	locks      int64
	dying      int64
	helpgc     int64
	spinning   bool
	blocked    bool
	lockedg    int64
}

type g struct {
	goid       int64
	status     int64
	waitreason string
	mid        int64
	lockedm    int64
}

func (r *record) preFill(line []byte) {
	secs := bytes.Split(line, []byte(" "))
	for i := 1; i < len(secs); i++ {
		if i == 1 {
			r.cost = string(secs[i][:len(secs[i])-1])
		} else {
			idx := bytes.Index(secs[i], []byte("="))
			if idx > 0 {
				switch string(secs[i][:idx]) {
				case "gomaxprocs":
					r.maxProc, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "idleprocs":
					r.idleprocs, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "threads":
					r.threads, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "spinningthreads":
					r.spinningthreads, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "idlethreads":
					r.idlethreads, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "runqueue":
					r.runqueue, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "gcwaiting":
					r.gcwaiting, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "nmidlelocked":
					r.nmidlelocked, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "stopwait":
					r.stopwait, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				case "sysmonwait":
					r.sysmonwait, _ = strconv.ParseInt(string(secs[i][idx+1:]), 10, 64)
				}
			}
		}
	}
}

func (r *record) pfill(plines [][]byte) {
	var ps []p
	for i := range plines {
		if len(plines[i]) < 3 {
			continue
		}
		if !bytes.Equal(plines[i][:3], []byte("  P")) {
			fmt.Printf("bad format for P:%s\n", string(plines[i]))
			continue
		}
		secs := bytes.Split(plines[i], []byte(" "))
		var _p = p{}
		for idx := range secs {
			index := bytes.Index(secs[idx], []byte("="))
			if index > 0 {
				switch string(secs[idx][:index]) {
				case "status":
					_p.status, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "schedtick":
					_p.schedtick, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "syscalltick":
					_p.syscalltick, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "m":
					var err error
					_p.mid, err = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
					if err != nil {
						fmt.Printf("err=%s\n", err)
					}
				case "runqsize":
					_p.runqsize, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "gfreecnt":
					_p.gfreecnt, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				}
			}
		}
		ps = append(ps, _p)
	}
	r.allp = ps
}

func (r *record) mfill(mlines [][]byte) {
	var ms []m
	for i := range mlines {
		if len(mlines[i]) < 3 {
			continue
		}
		if !bytes.Equal(mlines[i][:3], []byte("  M")) {
			fmt.Printf("bad format for M:%s\n", string(mlines[i]))
			continue
		}
		secs := bytes.Split(mlines[i], []byte(" "))
		var _m = m{}
		for idx := range secs {
			index := bytes.Index(secs[idx], []byte("="))
			if index > 0 {
				switch string(secs[idx][:index]) {
				case "p":
					_m.pid, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "curg":
					_m.curg, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "mallocing":
					_m.mallocing, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "throwing":
					_m.throwing, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "preemptoff":
				case "locks":
					var err error
					_m.locks, err = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
					if err != nil {
						fmt.Printf("err=%+v\n", err)
					}
				case "dying":
					_m.dying, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "helpgc":
					_m.helpgc, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "spinning":
				case "blocked":
				case "lockedg":
					_m.lockedg, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				}
			}
		}
		ms = append(ms, _m)
	}
	r.allm = ms
}

func (r *record) gfill(glines [][]byte) {
	var gs []g
	for i := range glines {
		if len(glines[i]) < 3 {
			continue
		}
		if !bytes.Equal(glines[i][:3], []byte("  G")) {
			fmt.Printf("bad format for G:%s\n", string(glines[i]))
			continue
		}
		var _g = g{}
		fmt.Printf("gline=%s\n", string(glines[i]))
		index1 := bytes.Index(glines[i], []byte("("))
		index2 := bytes.LastIndex(glines[i], []byte(")"))
		if index2 > index1 {
			_g.waitreason = string(glines[i][index1+1 : index2])
			glines[i] = append(glines[i][:index1], glines[i][index2+1:]...)
		}
		fmt.Printf("gline=%s\n", string(glines[i]))
		secs := bytes.Split(glines[i], []byte(" "))
		for idx := range secs {
			fmt.Printf("sec=%s\n", string(secs[idx]))
			index := bytes.Index(secs[idx], []byte("="))
			if index > 0 {
				switch string(secs[idx][:index]) {
				case "status":
					_g.status, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "m":
					_g.mid, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				case "lockedm":
					_g.lockedm, _ = strconv.ParseInt(string(secs[idx][index+1:]), 10, 64)
				}
			} else if index == 0 {
				index = bytes.Index(secs[idx], []byte(":"))
				_g.goid, _ = strconv.ParseInt(string(secs[idx][1:index]), 10, 64)
			}
		}
		gs = append(gs, _g)
	}
	r.allg = gs
}

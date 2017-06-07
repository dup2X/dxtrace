package dxtrace

import (
	"bytes"
	"errors"
)

var errBadFormat = errors.New("not expected format")

type textParser struct{}

func (tp *textParser) Parse(data []byte) (*record, error) {
	lines := bytes.Split(data, []byte("\n"))
	if len(lines) == 0 || !bytes.Equal(lines[0][:5], []byte("SCHED")) {
		return nil, errBadFormat
	}
	var r = &record{}
	r.preFill(lines[0])
	lines = lines[1:]
	r.pfill(lines[:r.maxProc])
	lines = lines[r.maxProc:]
	r.mfill(lines[:r.threads])
	lines = lines[r.threads:]
	r.gfill(lines)
	return r, nil
}

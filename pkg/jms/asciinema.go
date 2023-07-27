package jms

import (
	"encoding/json"
	"io"
	"time"
)

const (
	Version        = 2
	Width          = 80
	Height         = 40
	DefaultShell   = "/bin/bash"
	DefaultTerm    = "xterm"
	NewLine        = "\n"
	DateTimeFormat = "2006-01-02T15:04:05.999Z"
)

type AsciinemaWriter struct {
	Config        map[string]interface{}
	Writer        io.Writer
	TimestampNano int64
}

func NewAsciinemaWriter(writer io.Writer) *AsciinemaWriter {
	return &AsciinemaWriter{
		Config: map[string]interface{}{
			"width":     Width,
			"height":    Height,
			"envShell":  DefaultShell,
			"envTerm":   DefaultTerm,
			"timestamp": time.Now().Unix(),
			"title":     nil,
		},
		Writer:        writer,
		TimestampNano: time.Now().UnixNano(),
	}
}

func (aw *AsciinemaWriter) writeHeader() {
	header := map[string]interface{}{
		"version":   Version,
		"width":     aw.Config["width"],
		"height":    aw.Config["height"],
		"timestamp": aw.Config["timestamp"],
		"title":     aw.Config["title"],
		"env": map[string]interface{}{
			"shell": aw.Config["envShell"],
			"term":  aw.Config["envTerm"],
		},
	}
	jsonData, _ := json.Marshal(header)
	jsonData = append(jsonData, NewLine...)
	aw.Writer.Write(jsonData)
}

func (aw *AsciinemaWriter) writeRow(p []byte) {
	now := time.Now().UnixNano()
	ts := float64(now-aw.TimestampNano) / 1e9
	aw.writeStdout(ts, p)
}

func (aw *AsciinemaWriter) writeStdout(ts float64, data []byte) {
	row := []interface{}{ts, "o", string(data)}
	jsonData, _ := json.Marshal(row)
	jsonData = append(jsonData, NewLine...)
	aw.Writer.Write(jsonData)
}

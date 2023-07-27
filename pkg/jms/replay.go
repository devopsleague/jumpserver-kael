package jms

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	DefaultEncoding = "utf-8"
	ReplayDir       = "data/replay" // You need to set the correct path to the replay directory
)

type ReplayHandler struct {
	Session      *Session
	ReplayWriter *AsciinemaWriter
	FileWriter   *os.File
	File         *os.File
}

func NewReplayHandler(session *Session) *ReplayHandler {
	handler := &ReplayHandler{
		Session: session,
	}
	handler.buildFile()
	return handler
}

func (rh *ReplayHandler) buildFile() {
	rh.ensureReplayDir()

	replayFilePath := filepath.Join(ReplayDir, fmt.Sprintf("%s.cast", rh.Session.ID))
	file, err := os.Create(replayFilePath)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create replay file: %s -> %s", file.Name(), err)
		// Handle the error
		return
	}

	rh.File = file
	rh.FileWriter = file
	rh.ReplayWriter = NewAsciinemaWriter(file)
	rh.ReplayWriter.WriteHeader()
}

func (rh *ReplayHandler) ensureReplayDir() {
	err := os.MkdirAll(ReplayDir, os.ModePerm)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create replay directory: %s", ReplayDir)
		// Handle the error
	}
}

func (rh *ReplayHandler) writeRow(row string) {
	row = strings.ReplaceAll(row, "\n", "\r\n")
	row = strings.ReplaceAll(row, "\r\r\n", "\r\n")
	row = fmt.Sprintf("%s \r\n", row)

	_, err := rh.ReplayWriter.WriteRow([]byte(row))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write replay row: %s", err)
		// Handle the error
	}
}

func (rh *ReplayHandler) WriteInput(inputStr string) {
	// TODO: Convert the time to the desired format
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	inputStr = fmt.Sprintf("[%s]#: %s", formattedTime, inputStr)
	rh.writeRow(inputStr)
}

func (rh *ReplayHandler) WriteOutput(outputStr string) {
	// Wrap the output text to the desired width
	wrapper := tabwriter.NewWriter(rh.FileWriter, 0, 0, 1, ' ', 0)
	_, err := fmt.Fprintln(wrapper, outputStr)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write output: %s", err)
		// Handle the error
	}

	// Flush the tabwriter to ensure all data is written
	err = wrapper.Flush()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to flush tabwriter: %s", err)
		// Handle the error
	}
}

func (rh *ReplayHandler) Upload() {
	defer rh.FileWriter.Close()

	replayRequest := &ReplayRequest{
		SessionID:      rh.Session.ID,
		ReplayFilePath: rh.File.Name(),
	}

	resp, err := rh.Stub.UploadReplayFile(replayRequest)
	if err != nil || !resp.Status.Ok {
		errorMessage := fmt.Sprintf("Failed to upload replay file: %s %s", rh.File.Name(), resp.Status.Err)
		// Handle the error
	}
}

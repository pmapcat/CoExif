package main

import (
	"os/exec"
	"testing"
	// "unsafe"
)

func Test_long_running_proc(t *testing.T) {
	stdin_chan := make(chan string)
	stdout_chan := make(chan string)
	stderr_chan := make(chan string)
	cmd := exec.Command("cat")
	go long_running_proc(stdin_chan, stdout_chan, stderr_chan, cmd)
	stdin_chan <- "hello world!\n"
	msg := <-stdout_chan
	if msg != "hello world!" {
		t.Error(msg, "!=", "hello world!")
	}

	cmd = exec.Command("ping", "google.com")
	go long_running_proc(stdin_chan, stdout_chan, stderr_chan, cmd)
	msg = <-stdout_chan
	if msg == "hello world!" {
		t.Error(msg, "==", "hello world!")
	}

}

package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
)

func long_running_proc(stdin_chan,
	stdout_chan,
	stderr_chan chan string,
	cmd *exec.Cmd) (error, int) {
	var wg sync.WaitGroup
	wg.Add(3)
	// STDIN
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return err, 0
	}
	defer stdin.Close()

	// STDOUT
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err, 0
	}
	defer stdout.Close()

	// STDERR
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return err, 0
	}
	defer stderr.Close()

	// START
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err, 0
	}

	// WRITER ROUTINE
	go func() {
		defer wg.Done()
		for {
			pass_to_func := <-stdin_chan
			io.WriteString(stdin, pass_to_func)
		}
	}()

	// READER ROUTINE
	go func() {
		defer wg.Done()
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			stdout_chan <- in.Text()
		}
		if err := in.Err(); err != nil {
			log.Println(err)
		}
	}()
	// ERROR ROUTINE
	go func() {
		defer wg.Done()
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			stderr_chan <- in.Text()
		}
		if err := in.Err(); err != nil {
			log.Println(err)
		}
	}()
	wg.Wait()
	return nil, 1
}

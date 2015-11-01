package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	// "time"
)

// Algorithm:
//   Start process
// Execute:
//  Block

type LongProcess struct {
	Proc    *exec.Cmd
	stdin   io.WriteCloser
	Blocked bool
	Payload string
}

func (l *LongProcess) ExecSync(stdin string) string {
	// Will block until payload is received
	io.WriteString(l.stdin, stdin)
	for l.Blocked {
	}
	return l.Payload
}

func (l *LongProcess) Init(name string, params ...string) {
	cmd := exec.Command(name, params...)
	l.Proc = cmd
	l.Blocked = true
	// open stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	defer stdout.Close()

	// open stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
	}
	l.stdin = stdin
	defer stdin.Close()
	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	// attach writer to Payload
	// We assume that output is always one line(true, in case of ExifTool)
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		l.Payload = in.Text()
		log.Println("blab")
		l.Blocked = false
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
}

// func (e *ExifProcess) Exec (data string) string {
//   e.Blocked = true

// }

func main_log() (int, error) {
	cmd := exec.Command("cat")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	// open stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return 0, err
	}
	defer stdin.Close()

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
		return 0, err
	}

	io.WriteString(stdin, "fuck you in the ass\n")
	// read command's stdout line by line
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		log.Printf(in.Text()) // write each line to your log, or anything you need

	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	return 0, nil
}
func main() {
	main_log()
	// x := LongProcess{}
	// log.Println("init cat")
	// go x.Init("cat")
	// log.Println("sending hello")
	// result := x.ExecSync("Hello world!\n")
	// log.Println(result)

}

// func main_redirect() {
// 	// Replace `ls` (and its arguments) with something more interesting
// 	cmd := exec.Command("ping")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.Start()
// 	cmd.Stdin.Read([]byte("ya.ru"))
// 	log.Println("blab")
// 	for {
// 	}

// }
// func main_2() {
// 	cmd := exec.Command("ping", "ya.ru")
// 	randomBytes := &bytes.Buffer{}
// 	cmd.Stdout = randomBytes
// 	// Start command asynchronously
// 	err := cmd.Start()
// 	printError(err)

// 	// Create a ticker that outputs elapsed time
// 	go func() {
// 		for {
// 			b, err := randomBytes.ReadByte()
// 			if err != nil {
// 				log.Println(err)
// 			} else {
// 				log.Println(b)
// 			}

// 		}
// 	}()

// 	// Create a timer that will kill the process
// 	// timer := time.NewTimer(time.Second * 4)
// 	// go func(timer *time.Timer, ticker *time.Ticker, cmd *exec.Cmd) {
// 	// 	for _ = range timer.C {
// 	// 		err := cmd.Process.Signal(os.Kill)
// 	// 		printError(err)
// 	// 		ticker.Stop()
// 	// 	}
// 	// }(timer, ticker, cmd)

// 	// Only proceed once the process has finished
// 	cmd.Wait()
// 	printOutput(
// 		[]byte(fmt.Sprintf("%d bytes generated!", len(randomBytes.Bytes()))),
// 	)
// }

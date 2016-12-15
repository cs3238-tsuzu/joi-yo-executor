package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"sync/atomic"
	"syscall"
	"text/template"
	"time"
)

func main() {
	var stopped int32
	atomic.StoreInt32(&stopped, 0)

	inputTemp := os.Getenv("JOI_YO_IN_TEMPLATE")

	if inputTemp == "" {
		inputTemp = "2017-yo-t{{.Prob}}-in{{.Case}}.txt"
	}

	outTemp := os.Getenv("JOI_YO_OUT_TEMPLATE")

	if outTemp == "" {
		outTemp = "2017-yo-t{{.Prob}}-out{{.Case}}.txt"
	}

	wd, _ := os.Getwd()

	prob := path.Base(wd)

	if env := os.Getenv("JOI_YO_PROB"); env != "" {
		prob = env
	}

	exe := os.Getenv("JOI_YO_EXECUTABLE_PATH")

	if exe == "" {
		exe = "./a.out"
	}

	if len(os.Args) >= 2 {
		exe = os.Args[1]
	}

	for i := 1; i <= 5; i++ {
		cmd := exec.Command(exe)

		buf := bytes.NewBuffer(nil)

		template.Must(template.New("").Parse(inputTemp)).Execute(buf, map[string]string{"Prob": prob, "Case": strconv.FormatInt(int64(i), 10)})

		var err error
		cmd.Stdin, err = os.Open(buf.String())

		if err != nil {
			log.Println(i, ": Failed opening input file. Path:", buf.String(), " Reason:", err.Error())

			return
		}

		stderrBuf := bytes.NewBuffer(nil)
		cmd.Stderr = stderrBuf

		buf.Reset()
		template.Must(template.New("").Parse(outTemp)).Execute(buf, map[string]string{"Prob": prob, "Case": strconv.FormatInt(int64(i), 10)})

		output, err := os.Create(buf.String())

		if err != nil {
			log.Println(i, ": Failed opening output file. Path:", buf.String(), " Reason:", err.Error())

			return
		}

		cmd.Stdout = output

		startTime := time.Now()
		err = cmd.Start()

		if err != nil {
			log.Println(i, ": Failed executing. Reason:", err.Error())

			return
		}

		fin := make(chan bool, 1)
		signalManager := make(chan bool, 1)
		go func() {
			ch := make(chan os.Signal, 10)

			signal.Notify(ch,
				syscall.SIGHUP,
				syscall.SIGKILL,
				syscall.SIGTERM,
				syscall.SIGQUIT,
				syscall.SIGINT,
			)

			defer signal.Stop(ch)

			select {
			case <-ch:
				log.Println("Canceled")
				atomic.StoreInt32(&stopped, 1)
				cmd.Process.Kill()
			case <-fin:
			}
			signalManager <- true
		}()

		_, err = cmd.Process.Wait()
		finishTime := time.Now()

		fin <- true

		<-signalManager

		if atomic.LoadInt32(&stopped) == 1 {
			os.Exit(1)

			return
		}

		if err != nil {
			log.Println(i, ": Failed processing. Reason:", err.Error())

			return
		}
		log.Println(i, ": Finished", finishTime.Sub(startTime).String())

		if stderrBuf.Len() != 0 {
			log.Printf("%d: Standard Error:\r\n%s", i, stderrBuf.String())
		}
	}
}

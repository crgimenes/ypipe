package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/crgimenes/goconfig"
)

var (
	in  = os.Stdin
	out = os.Stdout
	buf = make([]byte, 4*1024)
)

func main() {
	type configFlags struct {
		Input   string `json:"i" cfg:"i" cfgDefault:"stdin" cfgHelper:"input from"`
		Output  string `json:"o" cfg:"o" cfgDefault:"stdout" cfgHelper:"output to"`
		OutFile string `json:"f" cfg:"f" cfgRequired:"true" cfgHelper:"output file name"`
	}

	cfg := configFlags{}
	goconfig.PrefixEnv = "YPIPE"
	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)

	f, err := os.Create(cfg.OutFile)
	fWriter := bufio.NewWriter(f)

	n := 0
	for {
		n, err = reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				printErrorln(err.Error())
				os.Exit(-1)
			}
			if n == 0 {
				break
			}
		}
		n, err = writer.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
		if n < len(buf[:n]) {
			fmt.Println(io.ErrShortWrite)
			return
		}
		n, err = fWriter.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}
		if n < len(buf[:n]) {
			fmt.Println(io.ErrShortWrite)
			return
		}
	}
	writer.Flush()
	fWriter.Flush()

}

func printErrorln(msg string) {
	fmt.Fprintf(os.Stderr, "%v\n", msg) // nolint
}

func write(w *bufio.Writer, buf []byte) {
	n, err := w.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if n < len(buf) {
		fmt.Println(io.ErrShortWrite)
		os.Exit(-1)
	}
}

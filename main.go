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
		OutFile string `json:"f" cfg:"f" cfgDefault:"" cfgHelper:"output file name"`
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
		writer.Write(buf[:n])
		fWriter.Write(buf[:n])
	}
	writer.Flush()
	fWriter.Flush()
}

func printErrorln(msg string) {
	fmt.Fprintf(os.Stderr, "%v\n", msg) // nolint
}

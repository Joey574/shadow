package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Joey574/shadow/v2/pkg/shell"
	"github.com/jessevdk/go-flags"
)

type Args struct {
	Output string `short:"o" long:"output" description:"output file, if none is specified stdout will be used"`

	// Shell args
	B64     bool   `short:"b" long:"b64" description:"encodes shell script to base 64, basic obfuscation"`
	Encrypt bool   `short:"e" long:"encrypt" description:"encrypts shell script, requiring a key to run it, if both encrypt and encode are selected, the script will be encrypted first"`
	Shell   string `long:"shell" description:"set the shell to use for the script" default:"sh"`
	Stride  int    `long:"stride" description:"set the stride length to use for shell encoding" default:"5"`
}

func main() {
	var args Args
	files, err := flags.Parse(&args)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalln(err)
	}

	// gather input
	if len(files) == 0 {
		log.Fatalln("no input files, see help (-h)")
	}

	var src string
	for _, f := range files {
		bytes, err := os.ReadFile(f)
		if err != nil {
			log.Fatalln(err)
		}

		src += string(bytes)
	}

	// set output
	out := os.Stdout
	if args.Output != "" {
		out, err = os.OpenFile(args.Output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}

	s := shell.NewObfuscator(
		shell.EncodeSrc(args.B64),
		shell.EncryptSrc(args.Encrypt),
		shell.SetShell(args.Shell),
		shell.SetStride(args.Stride),
	)

	script := s.Obfuscate(src)
	fmt.Fprint(out, script.Dump())
	if out == os.Stdout {
		fmt.Println()
	}

	if script.Key != nil {
		fmt.Println("Key:", string(script.Key))
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/winebarrel/sqsf"
)

var (
	version string
)

type flags struct {
	*sqsf.SqsfOpts
}

func init() {
	cmdLine := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	cmdLine.Usage = func() {
		fmt.Fprintf(cmdLine.Output(), "Usage: %s [OPTION] QUEUE\n", cmdLine.Name())
		cmdLine.PrintDefaults()
	}

	flag.CommandLine = cmdLine
}

func parseFlags() *flags {
	flags := &flags{SqsfOpts: &sqsf.SqsfOpts{}}
	flag.BoolVar(&flags.DecodeBody, "decode-body", false, "print decoded message body")
	flag.BoolVar(&flags.Delete, "delete", false, "delete received message")
	flag.IntVar(&flags.Limit, "limit", 0, "maximum number of received messages")
	flag.StringVar(&flags.MessageId, "message-id", "", "message ID to receive")
	showVersion := flag.Bool("version", false, "print version and exit")
	visibilityTimeout := flag.Int("vis-timeout", 600, "visibility timeout")
	flag.Parse()

	if *showVersion {
		printVersionAndExit()
	}

	args := flag.Args()

	if len(args) < 1 {
		printUsageAndExit()
	} else if len(args) > 1 {
		log.Fatal("too many arguments")
	}

	if flags.Delete && isFlagPassed("vis-timeout") {
		log.Fatal("cannot pass both '-delete=true' and `-vis-timeout`")
	}

	flags.QueueName = args[0]
	flags.VisibilityTimeout = int32(*visibilityTimeout)

	return flags
}

func printVersionAndExit() {
	v := version

	if v == "" {
		v = "<nil>"
	}

	fmt.Fprintln(flag.CommandLine.Output(), v)
	os.Exit(0)
}

func printUsageAndExit() {
	flag.CommandLine.Usage()
	os.Exit(0)
}

func isFlagPassed(name string) (passed bool) {
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			passed = true
		}
	})

	return
}

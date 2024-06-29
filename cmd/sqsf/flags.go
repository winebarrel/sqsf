package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/itchyny/gojq"
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
	query := ""
	flags := &flags{SqsfOpts: &sqsf.SqsfOpts{}}
	flag.BoolVar(&flags.DecodeBody, "decode-body", false, "print decoded message body")
	flag.BoolVar(&flags.Delete, "delete", false, "delete received message")
	flag.IntVar(&flags.Limit, "limit", 0, "maximum number of received messages")
	flag.StringVar(&flags.Region, "region", os.Getenv("AWS_REGION"), "AWS region ($AWS_REGION)")
	flag.StringVar(&flags.EndpointUrl, "endpoint-url", os.Getenv("AWS_ENDPOINT_URL"), "AWS endpoint URL ($AWS_ENDPOINT_URL)")
	flag.StringVar(&flags.MessageId, "message-id", "", "message ID to receive")
	flag.StringVar(&query, "query", "", "jq expression to filter the output")
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

	if query != "" {
		q, err := gojq.Parse(query)

		if err != nil {
			log.Fatalf("failed to parse jq expression: %s", err)
		}

		flags.Query = q
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

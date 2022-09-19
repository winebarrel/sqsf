package main

import (
	"context"
	"log"

	"github.com/winebarrel/sqsf"
)

func init() {
	log.SetFlags(0)
}

func main() {
	flags := parseFlags()
	ctx := context.Background()
	sqs, err := sqsf.NewClient(ctx, flags.SqsfOpts)

	if err != nil {
		log.Fatal(err)
	}

	err = sqs.Follow(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"fmt"
	cron "github.com/tgmendes/cron_parser"
	"log"
)

func main() {
	cronExpr := flag.String("e", "", "a string with a cron expression")

	flag.Parse()

	if *cronExpr == "" {
		log.Fatal("no cron expression was specified")
	}

	cronTable, err := cron.Parse(*cronExpr)
	if err != nil {
		log.Fatalf("error parsing the cron expression: %s", err.Error())
	}

	fmt.Print(cronTable)
}

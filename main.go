package main

import (
	"fmt"
	topfew "github.com/timbray/topfew/internal"
	"os"
)

func main() {
	var err error

	config, err := topfew.Configure(os.Args[1:]) // skip whatever go puts in os.Args[0]
	if err != nil {
		fmt.Println("Problem (tf -h for help): " + err.Error())
		os.Exit(1)
	}

	counts, err := topfew.Run(config, os.Stdin)
	if err != nil {
		os.Exit(1)
	}
	for _, kc := range counts {
		fmt.Printf("%d %s\n", *kc.Count, kc.Key)
	}
}

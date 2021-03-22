package main

import (
	"flag"
	"fmt"
	topfew "github.com/timbray/topfew/internal"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

func main() {
	size := flag.Uint("n", 10, "how many of the top results to display")
	fieldSpec := flag.String("fields", "", "fields (one or more comma-separated numbers)")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	var err error

	flag.Parse()
	fname := flag.Arg(0)
	var fields []uint
	if *fieldSpec == "" {
		fields = nil
	} else {
		fields = parseFields(*fieldSpec)
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "can't create profiler: %s", err.Error())
			return
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "can't start profiler: %s", err.Error())
			return
		}
		defer pprof.StopCPUProfile()
	}

	var kf = topfew.NewKeyFinder(fields)
	var topList []*topfew.KeyCount

	if fname == "" {
		topList, err = topfew.FromStream(os.Stdin, kf, *size)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error reading stream: %s\n", err.Error())
			return
		}
	} else {
		counter := topfew.NewCounter(*size)
		err = topfew.ReadFileInSegments(fname, counter, kf)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err.Error())
			return
		}
		topList = counter.GetTop()
	}

	for _, kc := range topList {
		fmt.Printf("%d %s\n", kc.Count, kc.Key)
	}
}

func parseFields(spec string) []uint {
	parts := strings.Split(spec, ",")
	var fields []uint
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "illegal field spec '%s'", part)
			os.Exit(1)
		}
		fields = append(fields, uint(num))
	}
	return fields
}

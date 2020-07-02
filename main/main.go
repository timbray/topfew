package main

import (
	"flag"
	"fmt"
	topfew "github.com/timbray/topfew/internal"
	"io"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
)

func main() {
	size := flag.Uint("few", 10, "how many is a few?")
	fieldSpec := flag.String("fields", "", "which fields?")
	mmap := flag.Bool("mmap", false, "use mmap rather than file reader")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

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
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "can't start profiler: %s", err.Error())
		}
		defer pprof.StopCPUProfile()
	}

	var reader io.Reader
	var err error
	if fname == "" {
		reader = os.Stdin
	} else {
		if *mmap {
			reader, err = topfew.NewMmap(fname)
		} else {
			reader, err = os.Open(fname)
		}
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Can’t open %s: %s", fname, err.Error())
		}
	}

	kf := topfew.NewKeyFinder(fields)
	topList, err := topfew.FromStream(reader, kf, *size)
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

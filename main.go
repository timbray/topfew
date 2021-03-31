package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"

	topfew "github.com/timbray/topfew/internal"
)

const instructions = `
Usage: tf
	-n, --number (output line count) [default 10]
	-f, --fields (field list) [default is the whole record]
	-g, --grep (regexp) [may repeat]
	-v, --vgrep (regexp) [may repeat]
	-s, --sed (regexp) (replacement) [may repeat]
	-w, --width (segment count) [default is result of runtime.numCPU()]
	--sample
	-h, -help, --help
	(filename) [optional, stdin if omitted]

Field list is comma-separated integers, e.g. -f 3 or --fields 1,3,7

The regexp-valued fields work as follows:
-g/--grep discardsrecords that don't match the regexp (g for grep)
-v/--vgrep discards records that do match the regexp (v for grep -v)
-s/--sed works on extracted fields, replacing regexp with replacement

The regexp-valued fields can be supplied multiple times; the filtering
will be performed in the order supplied.

It can be difficult to get the regular expressions right. "-sample"
causes topfew to read records and print out the results of the 
filtering activities. It only works on standard input.
`

func usage(err error) {
	fmt.Println(instructions)
	if err != nil {
		fmt.Println("Problem: " + err.Error())
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func main() {

	size := 10
	var err error
	var fields []uint
	var cpuprofile string
	var tracefname string
	var fname string
	var filter topfew.Filters
	sample := false
	width := 0

	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]
		switch {
		case arg == "-n" || arg == "--number":
			i++
			size, err = strconv.Atoi(os.Args[i])
			if err == nil && size < 1 {
				err = fmt.Errorf("invalid size %d", size)
			}
		case arg == "-f" || arg == "--fields":
			i++
			fields, err = parseFields(os.Args[i])
		case arg == "--cpuprofile":
			i++
			cpuprofile = os.Args[i]
		case arg == "--trace":
			i++
			tracefname = os.Args[i]
		case arg == "-g" || arg == "--grep":
			i++
			err = filter.AddGrep(os.Args[i])
		case arg == "-v" || arg == "--vgrep":
			i++
			err = filter.AddVgrep(os.Args[i])
		case arg == "-s" || arg == "--sed":
			err = filter.AddSed(os.Args[i+1], os.Args[i+2])
			i += 2
		case arg == "--sample":
			sample = true
		case arg == "-h" || arg == "-help" || arg == "--help":
			usage(nil)
		case arg == "-w" || arg == "--width":
			i++
			width, err = strconv.Atoi(os.Args[i])
			if err == nil && size < 1 {
				err = fmt.Errorf("invalid size %d", size)
			}

		default:
			if arg[0] == '-' {
				err = fmt.Errorf("Unexpected flag argument %v", arg)
			} else {
				fname = os.Args[i]
			}
		}
		if err != nil {
			usage(err)
		}
		i++
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
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
	if tracefname != "" {
		f, err := os.Create(tracefname)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "can't create trace output file: %s", err.Error())
			return
		}
		// The generated trace can be analyzed with: go tool trace <tracefile>
		trace.Start(f)
		defer trace.Stop()
	}
	var kf = topfew.NewKeyFinder(fields)
	var topList []*topfew.KeyCount

	if fname == "" {
		if sample {
			for i, sed := range filter.Seds {
				fmt.Printf("SED %d: s/%s/%s/\n", i, sed.ReplaceThis, sed.WithThat)
			}
			err = topfew.Sample(os.Stdin, &filter, kf)
		} else {
			topList, err = topfew.FromStream(os.Stdin, &filter, kf, size)
		}
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error reading stream: %s\n", err.Error())
			return
		}
	} else {
		counter := topfew.NewCounter(size)
		err = topfew.ReadFileInSegments(fname, &filter, counter, kf, width)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", fname, err.Error())
			return
		}
		topList = counter.GetTop()
	}

	for _, kc := range topList {
		fmt.Printf("%d %s\n", *kc.Count, kc.Key)
	}
}

func parseFields(spec string) ([]uint, error) {
	parts := strings.Split(spec, ",")
	var fields []uint
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("Illegal field spec: %v", err)
		}
		fields = append(fields, uint(num))
	}
	return fields, nil
}

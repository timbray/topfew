package topfew

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type config struct {
	size           int
	fields         []uint
	fieldSeparator *regexp.Regexp
	Fname          string
	filter         filters
	width          int
	sample         bool
	quotedFields   bool
}

func Configure(args []string) (*config, error) {
	// lifted out of main.go to facilitate testing
	config := config{size: 10}
	var err error

	i := 0
	for i < len(args) {
		arg := args[i]
		switch {
		case arg == "-n" || arg == "--number":
			if (i + 1) >= len(args) {
				err = errors.New("insufficient arguments for --number")
			} else {
				i++
				config.size, err = strconv.Atoi(args[i])
				if err == nil && config.size < 1 {
					err = fmt.Errorf("invalid size %d", config.size)
				}
			}
		case arg == "-f" || arg == "--fields":
			if (i + 1) >= len(args) {
				err = errors.New("insufficient arguments for --fields")
			} else {
				i++
				config.fields, err = parseFields(args[i])
			}
		case arg == "-p" || arg == "--fieldseparator":
			if (i + 1) >= len(args) {
				err = errors.New("insufficient arguments for --fieldseparator")
			} else {
				i++
				config.fieldSeparator, err = regexp.Compile(args[i])
			}
		case arg == "-g" || arg == "--grep":
			if (i + 1) >= len(args) {
				err = errors.New("insufficient arguments for --grep")
			} else {
				i++
				err = config.filter.addGrep(args[i])
			}
		case arg == "-v" || arg == "--vgrep":
			if (i + 1) >= len(args) {
				err = errors.New("insufficient arguments for --vgrep")
			} else {
				i++
				err = config.filter.addVgrep(args[i])
			}
		case arg == "-s" || arg == "--sed":
			if (i + 2) >= len(args) {
				err = errors.New("insufficient arguments for --sed")
			} else {
				err = config.filter.addSed(args[i+1], args[i+2])
				i += 2
			}
		case arg == "--sample":
			config.sample = true
		case arg == "--quotedfields" || arg == "-q":
			config.quotedFields = true
		case arg == "-h" || arg == "-help" || arg == "--help":
			fmt.Println(instructions)
			os.Exit(0)
		case arg == "-w" || arg == "--width":
			if (i + 1) >= len(args) {
				err = errors.New("insufficient arguments for --width")
			} else {
				i++
				config.width, err = strconv.Atoi(args[i])
				if err == nil && config.width < 1 {
					err = fmt.Errorf("invalid width %d", config.width)
				}
			}

		default:
			if arg[0] == '-' {
				err = fmt.Errorf("unexpected flag argument %v", arg)
			} else {
				config.Fname = args[i]
			}
		}
		if err != nil {
			return nil, err
		}
		i++
	}
	if (config.fieldSeparator != nil) && config.quotedFields {
		err = errors.New("only one of -p/--fieldseparator and -q/--quotedfields may be specified")
	}

	return &config, err
}

func parseFields(spec string) ([]uint, error) {
	parts := strings.Split(spec, ",")
	var fields []uint
	lastNum := -1
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("illegal field spec: %w", err)
		}
		if num <= lastNum {
			return nil, fmt.Errorf("field-number list must be in order; problem at \"%d\"", num)
		} else {
			lastNum = num
		}
		fields = append(fields, uint(num))
	}
	return fields, nil
}

const instructions = `
tf (short for "topfew") finds the most common values in a line-structured input
and prints the top few of them out, with their occurrence counts, in decreasing
order of occurrences.

Usage: tf
	-n, --number (output line count) [default is 10]
	-f, --fields (field list) [default is the whole record]
    -p, --fieldseparator (field separator regex) [default is white space]
	-q, --quotedfields [default is false]
	-g, --grep (regexp) [may repeat, default is accept all]
	-v, --vgrep (regexp) [may repeat, default is reject none]
	-s, --sed (regexp) (replacement) [may repeat, default is no changes]
	-w, --width (segment count) [default is result of runtime.numCPU()]
	--sample
	-h, -help, --help
	filename [default is stdin]

All the arguments are optional; if none are provided, tf will read records 
from the standard input and list the 10 which occur most often.

Field list is comma-separated integers, e.g. -f 3 or --fields 1,3,7. The fields
must be provided in order, so 3,1,7 is an error.

Fields are separated by white space (spaces or tabs) by default.
This can be overridden with the --fieldseparator option, at some cost in
performance.

Some files, for example Apache httpd logs, use space-separation but also
allow spaces within fields which are quoted with ("). The -q/--quotedfields
allows tf to process these correctly. It is an error to specify both
-p and -q.

The regexp-valued fields work as follows:
-g/--grep discards records that don't match the regexp (g for grep)
-v/--vgrep discards records that do match the regexp (v for grep -v)
-s/--sed works on extracted fields, replacing regexp with replacement

The regexp-valued fields can be supplied multiple times; the filtering
and substitution will be performed in the order supplied.

If the input is a named file, tf will process it in multiple parallel
threads, which can dramatically improve performance. The --width argument
allows you to specify the number of threads. The default value is not always 
optimal; experience with particular data on a particular computer may lead 
to finding a better value.

It can be difficult to get the regular expressions right. "--sample"
causes topfew to read records and print out the results of the 
filtering activities. It only works on standard input.`

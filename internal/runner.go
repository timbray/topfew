package topfew

import (
	"fmt"
	"io"
	"os"
)

func Run(config *config, instream io.Reader) ([]*keyCount, error) {
	// lifted out of main.go to facilitate testing
	var kf = newKeyFinder(config.fields)
	var topList []*keyCount
	var err error

	if config.Fname == "" {
		if config.sample {
			for i, sed := range config.filter.seds {
				fmt.Printf("SED %d: s/%s/%s/\n", i, sed.ReplaceThis, sed.WithThat)
			}
			err = sample(instream, &config.filter, kf)
		} else {
			topList, err = fromStream(instream, &config.filter, kf, config.size)
		}
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error reading stream: %s\n", err.Error())
			return nil, err
		}
	} else {
		counter := newCounter(config.size)
		err = readFileInSegments(config.Fname, &config.filter, counter, kf, config.width)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error processing %s: %s\n", config.Fname, err.Error())
			return nil, err
		}
		topList = counter.getTop()
	}

	return topList, err
}

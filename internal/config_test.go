package topfew

import "testing"

func TestArgSyntax(t *testing.T) {
	bads := [][]string{
		{"-donkey"},
		{"-n", "0"}, {"--number", "-3"}, {"-n", ""}, {"--number", "two"}, {"-n"},
		{"--fields"}, {"--sample", "-f"},
		{"-f", "a"}, {"-f", "1,2,z,4"}, {"-f", "1,3,2"},
		{"--sample", "--cpuprofile"}, {"--cpuprofile"},
		{"--grep"}, {"--sample", "-g"},
		{"--vgrep"}, {"--sample", "-vg"},
		{"--sample", "--trace"}, {"--trace"},
		{"--sed"}, {"-s", "x"}, {"--sample", "--sed", "1"},
		{"--width", "a"}, {"-w", "0"}, {"--sample", "-w"},
		// COMMENT OUT FOLLOWING TO ENABLE TRACING
		{"--cpuprofile", "/tmp/cp"},
		{"--trace", "/tmp/tr"},
	}

	// not testing -h/--help because it'd be extra work to avoid printing out the usage
	goods := [][]string{
		{"--number", "1"}, {"-n", "5"},
		{"--fields", "1"}, {"-f", "3,5"},
		{"--grep", "re1"}, {"-g", "re2"},
		{"--vgrep", "re1"}, {"-v", "re2"},
		{"--sed", "foo", "bar"}, {"-s", "z", ""},
		{"--sample"},
		{"--width", "2"}, {"-w", "3"},
		{"--sample", "fname"},
		/* == ENABLE PROFILING ==
		{"--cpuprofile", "/tmp/cp"},
		{"--trace", "/tmp/tr"},
		*/
	}

	for _, bad := range bads {
		var err error
		_, err = Configure(bad)
		if err == nil {
			t.Error("accepted bogus argument: " + bad[0])
		}
	}

	for _, good := range goods {
		var err error
		var c *Config
		c, err = Configure(good)
		if err != nil || c == nil {
			t.Error("rejected good argument: " + good[0])
		}
	}
}

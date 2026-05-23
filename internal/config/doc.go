// Package config provides CLI flag parsing and validation for logslice.
//
// # Usage
//
// Call Parse with a *flag.FlagSet and the argument slice (typically os.Args[1:]):
//
//	cfg, err := config.Parse(flag.CommandLine, os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Flags
//
//	-src        path to the input log file; use "-" for stdin (default: "-")
//	-start      inclusive start of the time window (RFC3339)
//	-end        inclusive end of the time window (RFC3339)
//	-fmt        explicit Go time layout used to parse log timestamps
//	-progress   how often to print a progress report (default: 5s; 0 to disable)
//	-version    print version information and exit
//
// Positional arguments after the flags are treated as output destinations.
// If none are provided, output defaults to stdout ("-").
package config

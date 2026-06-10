package config

import (
	"flag"
)

type Args struct {
	Help    bool
	Version bool
	Debug   bool
}

func ParseArgs() (Args, error) {
	args := Args{}
	flag.BoolVar(&args.Debug, "d", false, "Enable debug logging")
	flag.BoolVar(&args.Help, "h", false, "Show help")
	flag.BoolVar(&args.Version, "v", false, "Show version")

	flag.Parse()

	return args, nil
}

func PrintHelp() {
	flag.Usage()
}

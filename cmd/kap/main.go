package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/kap"
)

var (
	version string
)

func parseArgs() *kap.Options {
	var cli struct {
		kap.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	server := kap.NewServer(options)
	err := server.Run()

	if err != nil {
		log.Fatal(err)
	}
}

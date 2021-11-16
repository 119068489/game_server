package main

import (
	"flag"
	"game_server/backstage"
	"game_server/easygo"
	"os"
)

func main() {
	defer easygo.PanicWriter.Flush()
	defer easygo.RecoverAndLog()

	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	backstage.Entry(flagSet, os.Args[1:])
}

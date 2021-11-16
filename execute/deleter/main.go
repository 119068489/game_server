package main

import (
	"flag"
	"game_server/deleter"
	"game_server/easygo"
	"os"
)

func main() {
	defer easygo.PanicWriter.Flush()
	defer easygo.RecoverAndLog()

	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	deleter.Entry(flagSet, os.Args[1:])
}

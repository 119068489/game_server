package main

import (
	"flag"
	"game_server/easygo"
	"game_server/shop"
	"os"
)

func main() {
	defer easygo.PanicWriter.Flush()
	defer easygo.RecoverAndLog()

	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	shop.Entry(flagSet, os.Args[1:])

}

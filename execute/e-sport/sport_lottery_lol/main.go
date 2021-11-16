package main

import (
	"flag"
	"game_server/e-sports/sport_lottery_lol"
	"game_server/easygo"
	"os"
)

func main() {
	defer easygo.PanicWriter.Flush()
	defer easygo.RecoverAndLog()

	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	sport_lottery_lol.Entry(flagSet, os.Args[1:])

}

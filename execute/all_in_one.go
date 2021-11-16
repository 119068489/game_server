package main

import (
	"flag"
	"fmt"
	"game_server/backstage"
	"game_server/e-sports/sport_api"
	"game_server/e-sports/sport_apply"
	"game_server/e-sports/sport_crawl"
	"game_server/e-sports/sport_lottery_lol"
	"game_server/e-sports/sport_lottery_wzry"
	"game_server/easygo"
	"game_server/hall"
	"game_server/login"
	"game_server/shop"
	"game_server/square"
	"game_server/statistics"
	"game_server/wish"
	"os"
)

func main() {
	defer easygo.PanicWriter.Flush()
	defer easygo.RecoverAndLog()

	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	appName := flagSet.String("app", "", "-app=hall,子类型是 hall,backstage,sub_game,robot 中的其中一种")

	if len(os.Args) < 2 {
		flagSet.Usage()
		return
	}
	flagSet.Parse(os.Args[1:2]) // app_name 必须是第 1 个参数
	if *appName == "" {
		flagSet.Usage()
		return
	}

	switch *appName {
	case "hall":
		hall.Entry(flagSet, os.Args[2:])
	case "backstage":
		backstage.Entry(flagSet, os.Args[2:])
	case "login":
		login.Entry(flagSet, os.Args[2:])
	case "shop":
		shop.Entry(flagSet, os.Args[2:])
	case "square":
		square.Entry(flagSet, os.Args[2:])
	case "statistics":
		statistics.Entry(flagSet, os.Args[2:])
	case "wish":
		wish.Entry(flagSet, os.Args[2:])
	//case "sport_lottery_csgo":
	//	sport_lottery_csgo.Entry(flagSet, os.Args[2:])
	//case "sport_lottery_dota":
	//	sport_lottery_dota.Entry(flagSet, os.Args[2:])
	case "sport_lottery_lol":
		sport_lottery_lol.Entry(flagSet, os.Args[2:])
	case "sport_lottery_wzry":
		sport_lottery_wzry.Entry(flagSet, os.Args[2:])
	case "sport_apply":
		sport_apply.Entry(flagSet, os.Args[2:])
	case "sport_api":
		sport_api.Entry(flagSet, os.Args[2:])
	case "sport_crawl":
		sport_crawl.Entry(flagSet, os.Args[2:])
	default:
		s := fmt.Sprintf("不支持的参数 %s", *appName)
		panic(s)
	}
}

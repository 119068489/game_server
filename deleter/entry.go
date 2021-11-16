// 用于删除运营数据
// 删除数据库 报表/局单/注单
// 删除磁盘 金币日志等
package deleter

import "flag"

func Entry(flagSet *flag.FlagSet, args []string) {
	day := flagSet.Int("day", 0, "要删除几天前的数据")
	flagSet.Parse(args)

	if *day == 0 {
		flagSet.Usage()
	}
}

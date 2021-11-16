package main

import (
	"fmt"
	"os"

	"github.com/Luxurioust/excelize"
)

func main() {
	xlsx, err := excelize.OpenFile("randname.xlsx")
	if err != nil {
		panic(err)
	}

	test := "package for_game\n\nvar ManRandName = []string{\n"
	var i int32 = 2
	for {
		line := fmt.Sprintf("A%d", i)
		s, _ := xlsx.GetCellValue("Sheet2", line)
		if s == "" {
			test += "}\n"
			break
		}
		test += fmt.Sprintf("\"%s\",\n", s)
		i++
	}

	test1 := "var WomanRandName = []string{\n"
	var i1 int32 = 2
	for {
		line := fmt.Sprintf("A%d", i1)
		s, _ := xlsx.GetCellValue("Sheet1", line)
		if s == "" {
			test1 += "}\n"
			break
		}
		test1 += fmt.Sprintf("\"%s\",\n", s)
		i1++
	}

	file, err := os.OpenFile("../for_game/randname.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(1, err)
		return
	}
	defer file.Close()
	file.WriteString(test)
	file.WriteString(test1)

}

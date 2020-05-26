/*
* @Author: scottxiong
* @Date:   2020-05-26 20:55:30
* @Last Modified by:   scottxiong
* @Last Modified time: 2020-05-26 20:58:30
 */
package work

import (
	"fmt"
	"github.com/scott-x/gutils/parse"
)

func Run() {
	tbs := parse.GetTables("temp.sql")
	for _, table := range *tbs {
		fmt.Println("table:" + table.Name)
		for _, field := range table.Fields {
			fmt.Println(field.Name + ":" + field.Type)
		}
		fmt.Println("***********************")
	}
}

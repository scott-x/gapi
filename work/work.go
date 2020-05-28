/*
* @Author: scottxiong
* @Date:   2020-05-26 20:55:30
* @Last Modified by:   scottxiong
* @Last Modified time: 2020-05-28 21:38:28
 */
package work

import (
	"github.com/scott-x/gutils/cmd"
	"github.com/scott-x/gutils/fs"
	"github.com/scott-x/gutils/model"
	"github.com/scott-x/gutils/parse"
	"github.com/scott-x/gutils/str"
	"log"
	"os/exec"
)

type SQL_FILE struct {
	pth  string
	text string
}

type SQL_FILES []SQL_FILE

func Run() {
	tbs := parse.GetTables("temp.sql")
	sql_files := &SQL_FILES{}
	for _, table := range *tbs {
		sql_file := &SQL_FILE{}
		sql_file.pth = table.Name + ".go"
		var ctx string
		ctx += "package db\n\n"
		ctx += "import (\n"
		ctx += fs.Tab(4) + "\"log\"\n"
		ctx += ")\n\n"
		tb_name := str.FirstLetterToUpper(table.Name, 1)
		ctx += "type " + tb_name + " struct {\n"

		for _, field := range table.Fields {
			ctx += fs.Tab(4) + str.FirstLetterToUpper(field.Name, 1) + " " + field.Type + "`json: \"" + field.Name + "\"`\n"
		}
		ctx += "}\n\n"

		ctx += "type " + tb_name + "s []" + tb_name + "\n\n"

		gen_code(table, &ctx)
		sql_file.text = ctx
		*sql_files = append(*sql_files, *sql_file)
	}

	gen_files(sql_files)

}

func gen_files(sql_files *SQL_FILES) {
	//check if directory db exists or not
	fs.CreateDirIfNotExist("./db")
	for _, sql_file := range *sql_files {
		//will create file if not exists
		var target = "./db/" + sql_file.pth
		if !fs.IsExist(target) {
			fs.WriteString(target, sql_file.text)
			err := exec.Command("go", "fmt", "db/"+sql_file.pth).Run()
			if err != nil {
				log.Printf("Error:%s\n", err)
			}
			cmd.Info(target[2:] + " was created")
		}
	}
}

func gen_code(table model.Table, ctx *string) {

	//func getall
	tb_name := str.FirstLetterToUpper(table.Name, 1)

	*ctx += "//get all " + table.Name + "s\n"
	*ctx += "func GetAll" + tb_name + "s() *" + tb_name + "s {\n"
	*ctx += fs.Tab(4) + "//define sql\n"
	*ctx += fs.Tab(4) + "sql := \"select "
	for k, field := range table.Fields {
		*ctx += field.Name
		if k != len(table.Fields)-1 {
			*ctx += ", "
		}
	}
	*ctx += " from " + table.Name + "\"\n\n"

	*ctx += fs.Tab(4) + "//prepare\n"
	*ctx += fs.Tab(4) + "stmt, err := dbCon.Prepare(sql)\n"
	*ctx += fs.Tab(4) + "if err != nil {\n"
	*ctx += fs.Tab(8) + "log.Printf(\"GetAll" + tb_name + "s() Error: %s\", err)\n"
	*ctx += fs.Tab(8) + "return nil\n"
	*ctx += fs.Tab(4) + "}\n"
	*ctx += fs.Tab(4) + "defer stmt.Close()\n\n"

	*ctx += fs.Tab(4) + "//query\n"
	*ctx += fs.Tab(4) + "row, err := stmt.Query()\n"
	*ctx += fs.Tab(4) + "if err != nil {\n"
	*ctx += fs.Tab(8) + "log.Printf(\"GetAll" + tb_name + "s() Error: %s\", err)\n"
	*ctx += fs.Tab(8) + "return nil\n"
	*ctx += fs.Tab(4) + "}\n\n"

	*ctx += fs.Tab(4) + table.Name + "s := &" + tb_name + "s{}\n"
	*ctx += fs.Tab(4) + "for row.Next() {\n"
	*ctx += fs.Tab(8) + table.Name + " := &" + tb_name + "{}\n"
	*ctx += fs.Tab(8) + "row.Scan("
	for k, field := range table.Fields {
		*ctx += "&" + table.Name + "." + str.FirstLetterToUpper(field.Name, 1)
		if k != len(table.Fields)-1 {
			*ctx += ", "
		}
	}
	*ctx += ")\n"

	*ctx += fs.Tab(8) + "*" + table.Name + "s = append(*" + table.Name + "s, *" + table.Name + ")\n"
	*ctx += fs.Tab(4) + "}\n"
	*ctx += fs.Tab(4) + "return " + table.Name + "s\n"
	*ctx += "}\n\n"

	//func limited
	*ctx += "//get " + table.Name + "s in one page\n"
	*ctx += "func Get" + tb_name + "s(page, size int) *" + tb_name + "s {\n"
	*ctx += fs.Tab(4) + "//define sql\n"
	*ctx += fs.Tab(4) + "sql := \"select "
	for k, field := range table.Fields {
		*ctx += field.Name
		if k != len(table.Fields)-1 {
			*ctx += ", "
		}
	}
	*ctx += " from " + table.Name + " order by id limit ?,?\"\n\n"

	*ctx += fs.Tab(4) + "//prepare\n"
	*ctx += fs.Tab(4) + "stmt, err := dbCon.Prepare(sql)\n"
	*ctx += fs.Tab(4) + "if err != nil {\n"
	*ctx += fs.Tab(8) + "log.Printf(\"Get" + tb_name + "s() Error: %s\", err)\n"
	*ctx += fs.Tab(8) + "return nil\n"
	*ctx += fs.Tab(4) + "}\n"
	*ctx += fs.Tab(4) + "defer stmt.Close()\n\n"

	*ctx += fs.Tab(4) + "//query\n"
	*ctx += fs.Tab(4) + "row, err := stmt.Query((page-1)*size, size)\n"
	*ctx += fs.Tab(4) + "if err != nil {\n"
	*ctx += fs.Tab(8) + "log.Printf(\"Get" + tb_name + "s() Error: %s\", err)\n"
	*ctx += fs.Tab(8) + "return nil\n"
	*ctx += fs.Tab(4) + "}\n\n"

	*ctx += fs.Tab(4) + table.Name + "s := &" + tb_name + "s{}\n"
	*ctx += fs.Tab(4) + "for row.Next() {\n"
	*ctx += fs.Tab(8) + table.Name + " := &" + tb_name + "{}\n"
	*ctx += fs.Tab(8) + "row.Scan("
	for k, field := range table.Fields {
		*ctx += "&" + table.Name + "." + str.FirstLetterToUpper(field.Name, 1)
		if k != len(table.Fields)-1 {
			*ctx += ", "
		}
	}
	*ctx += ")\n"

	*ctx += fs.Tab(8) + "*" + table.Name + "s = append(*" + table.Name + "s, *" + table.Name + ")\n"
	*ctx += fs.Tab(4) + "}\n"
	*ctx += fs.Tab(4) + "return " + table.Name + "s\n"
	*ctx += "}\n\n"

	//func get one by id
	var type_id string

	for _, field := range table.Fields {
		if field.Name == "id" {
			type_id = field.Type
		}
	}
	*ctx += "//get one\n"
	*ctx += "func Get" + tb_name + "ById(id " + type_id + ") *" + tb_name + " {\n"
	*ctx += fs.Tab(4) + "//define sql\n"
	*ctx += fs.Tab(4) + "sql := \"select "
	for k, field := range table.Fields {
		*ctx += field.Name
		if k != len(table.Fields)-1 {
			*ctx += ", "
		}
	}
	*ctx += " from " + table.Name + " where id = ?\"\n\n"

	*ctx += fs.Tab(4) + "//prepare\n"
	*ctx += fs.Tab(4) + "stmt, err := dbCon.Prepare(sql)\n"
	*ctx += fs.Tab(4) + "if err != nil {\n"
	*ctx += fs.Tab(8) + "log.Printf(\"Get" + tb_name + "ById() Error: %s\", err)\n"
	*ctx += fs.Tab(8) + "return nil\n"
	*ctx += fs.Tab(4) + "}\n"
	*ctx += fs.Tab(4) + "defer stmt.Close()\n\n"

	*ctx += fs.Tab(4) + table.Name + ":= &" + tb_name + "{}\n"

	*ctx += fs.Tab(4) + "err = stmt.QueryRow(id).Scan("
	for k, field := range table.Fields {
		*ctx += "&" + table.Name + "." + str.FirstLetterToUpper(field.Name, 1)
		if k != len(table.Fields)-1 {
			*ctx += ", "
		}
	}
	*ctx += ")\n"

	*ctx += fs.Tab(4) + "if err != nil {\n"
	*ctx += fs.Tab(8) + "log.Printf(\"Get" + tb_name + "ById() Error: %s\", err)\n"
	*ctx += fs.Tab(8) + "return nil\n"
	*ctx += fs.Tab(4) + "}\n"

	*ctx += fs.Tab(4) + "return " + table.Name + "\n"
	*ctx += "}"
}

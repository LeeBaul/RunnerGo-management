package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "test_read_only:db#01^st$Post@tcp(rm-2zem14s80lyu5c4z7.mysql.rds.aliyuncs.com:3306)/kunpeng?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/pkg/dal/query",
	})

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("target"),
		g.GenerateModel("api"),
		g.GenerateModel("folder"),
		g.GenerateModel("group"),
		g.GenerateModel("operation"),
		g.GenerateModel("user"),
		g.GenerateModel("team"),
		g.GenerateModel("user_team"),
	)

	g.Execute()
}

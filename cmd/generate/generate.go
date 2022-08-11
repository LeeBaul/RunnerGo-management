package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:jJWp01F73zx2F5tqcZt6@tcp(127.0.0.1:23306)/apipost?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/pkg/dal/query",
	})

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("machine"),
	)

	g.Execute()
}

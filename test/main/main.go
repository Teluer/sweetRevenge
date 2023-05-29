package main

import (
	"fmt"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/websites/target"
)

func main() {
	dao.Dao.OpenDatabaseConnection("goblin:password1!@tcp(localhost:3306)/sweet?charset=utf8mb4&parseTime=True&loc=Local")
	fmt.Println(target.CreateRandomCustomer())
}

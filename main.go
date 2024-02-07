package main

import (
	"homework/api/router"
	"homework/common"
)

func main() {
	router.InitRouter()
	common.InitDatabase()
}

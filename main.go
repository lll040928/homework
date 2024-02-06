package main

import (
	"homework/api/router"
	"homework/common"
)

func main() {
	common.SetupViper()
	common.InitDatabase()
	router.InitRouter()
}

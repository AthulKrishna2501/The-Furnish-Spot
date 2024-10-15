package main

import (
	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/route"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDatabase()
	router := gin.Default()
	route.RegisterURL(router)
	router.Run(":3000")

}

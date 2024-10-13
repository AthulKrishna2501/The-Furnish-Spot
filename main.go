package main

import (
	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDatabase()
	router := gin.Default()
	router.Run(":3000")

}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	InitRoutes(app)

	fmt.Println(app.Run(":4848"))
}

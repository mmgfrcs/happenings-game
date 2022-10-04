package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, World!")
	app := gin.Default()

	InitRoutes(app)

	fmt.Println(app.Run(":4848"))
}

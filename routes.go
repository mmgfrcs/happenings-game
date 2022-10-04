package main

import (
  "github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
  app.LoadHTMLGlob("res/views/**/*.html")
  app.Static("/assets", "./res/public")
  app.GET("/", func(ctx *gin.Context) {
    ctx.HTML(200, "base.html", map[string]interface{}{})
  })
}


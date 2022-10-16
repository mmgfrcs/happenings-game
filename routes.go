package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"main/entities"
	"main/util"

	"github.com/gin-gonic/gin"
)

//go:embed res/data/***/*
var resF embed.FS

func InitRoutes(app *gin.Engine) {
  pData := util.HomeData{Characters: []entities.Character{}}
  
  fs.WalkDir(resF, "res/data/characters", func(path string, d fs.DirEntry, err error) error {
    if !d.IsDir() {
      f, err := resF.ReadFile("res/data/characters/ellie.json")
      if err != nil {
        return err
      }
      var data entities.Character
      err = json.Unmarshal(f, &data)
      if err != nil {
        return err
      }
      fmt.Print(data["name"])
    }
    
    return nil
  })
  
  app.LoadHTMLGlob("res/views/**/*.html")
  app.Static("/assets", "./res/public")
  app.GET("/", func(ctx *gin.Context) {
    ctx.HTML(200, "base.html", data)
  })
}


package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmgfrcs/happenings-game/game"
	"github.com/mmgfrcs/happenings-game/util"
)

//go:embed src/**/*
var embedFS embed.FS

func InitRoutes(app *gin.Engine) {
	manager := game.InitGame()

	pData := util.HomeData{Characters: manager.GetCharactersMap()}

	app.Use(func(ctx *gin.Context) {
		ctx.Next()
		for _, err := range ctx.Errors {
			fmt.Println("Error:", err.Err)
		}
		if len(ctx.Errors) > 0 {
			ctx.HTML(-1, "error.html", map[string]interface{}{"Code": 500})
		}
	})
	// app.LoadHTMLGlob("src/views/**/*.html")
	// app.Static("/assets", "./src/public")
	app.SetHTMLTemplate(template.Must(template.ParseFS(embedFS, "src/views/**/*.html")))
	subFs, err := fs.Sub(embedFS, "src/public")
	if err != nil {
		panic(err)
	}
	app.StaticFS("/assets", http.FS(subFs))

	app.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", pData)
	})

	app.GET("/start/:char", func(ctx *gin.Context) {
		character := ctx.Param("char")
		actions := ctx.Query("action")
		fmt.Println(actions)

		if char, ok := pData.Characters[character]; ok {
			ctx.SetCookie("happ-game-id", "1234", 3600, "/", "", true, true)
			ctx.HTML(200, "game.html", map[string]interface{}{"Character": char})
		} else {
			ctx.HTML(404, "error.html", map[string]interface{}{"Code": 404})
			ctx.AbortWithError(-1, errors.New("Character not found"))
		}
	})
}

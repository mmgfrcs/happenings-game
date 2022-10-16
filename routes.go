package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"main/data"
	"main/entities"
	"main/util"

	"github.com/gin-gonic/gin"
)

//go:embed res/data/***/*
var resF embed.FS

func InitRoutes(app *gin.Engine) {
	pData := util.HomeData{Characters: []data.Character{}}
	attrSlice := []entities.Attribute{}

	fs.WalkDir(resF, "res/data/attributes", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			f, err := resF.ReadFile(path)
			if err != nil {
				return err
			}
			var d entities.Attribute
			err = json.Unmarshal(f, &d)
			if err != nil {
				return err
			}
			attrSlice = append(attrSlice, d)
		}

		return nil
	})

	fs.WalkDir(resF, "res/data/characters", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			f, err := resF.ReadFile(path)
			if err != nil {
				return err
			}
			var d entities.Character
			err = json.Unmarshal(f, &d)
			if err != nil {
				return err
			}

			charData := data.Character{ID: d.ID, Name: d.Name, Description: d.Description, Attributes: make([]data.Attribute, 0)}
			for _, v := range attrSlice {
				charData.Attributes = append(charData.Attributes, data.Attribute{
					ShortName:   v.ShortName,
					Name:        v.Name,
					Description: v.Description,
					Range:       v.Range,
					Critical:    v.Critical,
					Value:       d.Attributes[v.ShortName],
				})
			}

			pData.Characters = append(pData.Characters, charData)
		}

		return nil
	})

	app.LoadHTMLGlob("res/views/**/*.html")
	app.Static("/assets", "./res/public")
	app.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "base.html", pData)
	})
}

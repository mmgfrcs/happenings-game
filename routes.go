package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"main/data"
	"main/entities"
	"main/util"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

//go:embed res/data/**/*
var resF embed.FS

func InitRoutes(app *gin.Engine) {
	pData := util.HomeData{Characters: map[string]data.Character{}}
	actions := []data.Action{}
	attrSlice := []entities.Attribute{}
	attrMap := map[string]entities.Attribute{}

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
			attrMap[d.ShortName] = d
		}

		return nil
	})

	err := fs.WalkDir(resF, "res/data/characters", func(path string, d fs.DirEntry, err error) error {
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
				attrVal, ok := d.Attributes[v.ShortName]
				if !ok {
					continue
				}
				charData.Attributes = append(charData.Attributes, data.Attribute{
					ShortName:   v.ShortName,
					Name:        v.Name,
					Description: fmt.Sprintf(v.Description, d.Name),
					Range:       v.Range,
					Critical:    v.Critical,
					Value:       attrVal,
				})
			}

			pData.Characters[charData.ID] = charData
		}

		return nil
	})

	err = fs.WalkDir(resF, "res/data/actions", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			f, err := resF.ReadFile(path)
			if err != nil {
				return err
			}
			var acts []entities.Action
			err = json.Unmarshal(f, &acts)
			if err != nil {
				return err
			}
      
			for _, act := range acts {
        actData := data.Action{
          ID: act.ID,
          Name: act.Name,
          Description: act.Description,
          Tags: act.Tags,
          Mods: []data.ActionMod{},
        }
				for _, mod := range act.Mods {
          reflect.New()
          mapstructure.NewDecoder(&mapstructure.DecoderConfig{
            Squash: true,
            Result: 
          })
				  if _, ok := attrMap[mod.AttributeID]; !ok {
				     panic(fmt.Sprintf("Attribute %s not found for action %s", mod.AttributeID, act.Name))
				  }
				}
				actions = append(actions, actData)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded", len(attrSlice), "Attributes,", len(pData.Characters), "Characters, and", len(actions), "Actions")

	app.LoadHTMLGlob("res/views/**/*.html")

	app.Static("/assets", "./res/public")
	app.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", pData)
	})

	app.GET("/start/:char", func(ctx *gin.Context) {
		character := ctx.Param("char")
		actions := ctx.Query("action")
		fmt.Println(actions)

		if char, ok := pData.Characters[character]; ok {
			ctx.HTML(200, "game.html", map[string]interface{}{"Character": char})
		} else {
			ctx.HTML(404, "error.html", map[string]interface{}{"Code": 404})
		}

	})
}

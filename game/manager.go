package game

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/mmgfrcs/happenings-game/game/data"
	"github.com/mmgfrcs/happenings-game/game/entities"
)

//go:embed res/**/*
var resF embed.FS

type ActionTypeError struct {
	Action string
	Mod    int
}

func (m *ActionTypeError) Error() string {
	return fmt.Sprintf("Action %s mod %d has incorrect/invalid type", m.Action, m.Mod)
}

type GameManager struct {
	characters []entities.Character
	attributes []entities.Attribute
	actions    []entities.Action

	actionTypes map[string]interface{}
}

func InitGameEmpty() GameManager {
	manager := GameManager{
		characters: []entities.Character{},
		attributes: []entities.Attribute{},
		actions:    []entities.Action{},
		actionTypes: map[string]interface{}{
			"mod":  entities.ActionModAddSub{},
			"mult": entities.ActionModMultDiv{},
			"comp": entities.ActionModCompare{},
		},
	}
	return manager
}

func InitGame() GameManager {
	manager := InitGameEmpty()
	attrMap := map[string]entities.Attribute{}

	fs.WalkDir(resF, "res/attributes", func(path string, d fs.DirEntry, err error) error {
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
			manager.attributes = append(manager.attributes, d)
			attrMap[d.ShortName] = d
		}

		return nil
	})

	err := fs.WalkDir(resF, "res/characters", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			f, err := resF.ReadFile(path)
			if err != nil {
				return err
			}
			var d data.Character
			err = json.Unmarshal(f, &d)
			if err != nil {
				return err
			}

			charData := entities.Character{ID: d.ID, Name: d.Name, Description: d.Description, Attributes: make([]entities.Attribute, 0)}
			for _, v := range manager.attributes {
				attrVal, ok := d.Attributes[v.ShortName]
				if !ok {
					continue
				}
				charData.Attributes = append(charData.Attributes, entities.Attribute{
					ShortName:   v.ShortName,
					Name:        v.Name,
					Description: fmt.Sprintf(v.Description, d.Name),
					Range:       v.Range,
					Critical:    v.Critical,
					Value:       attrVal,
				})
			}
			manager.characters = append(manager.characters, charData)
		}

		return nil
	})

	err = fs.WalkDir(resF, "res/actions", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			f, err := resF.ReadFile(path)
			if err != nil {
				return err
			}
			err = manager.LoadActionsByte(f)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded", len(manager.attributes), "Attributes,", len(manager.characters), "Characters, and", len(manager.actions), "Actions")

	return manager
}

func (gm *GameManager) GetCharactersMap() (charMap map[string]entities.Character) {
	charMap = map[string]entities.Character{}
	for _, char := range gm.characters {
		charMap[char.ID] = char
	}
	return
}

func (gm *GameManager) LoadActionsFile(path string) error {
	byt, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return gm.LoadActionsByte(byt)
}

func (gm *GameManager) LoadActionsByte(b []byte) error {
	var acts []data.Action
	err := json.Unmarshal(b, &acts)
	if err != nil {
		return err
	}
	for _, act := range acts {
		actData := entities.Action{
			ID:          act.ID,
			Name:        act.Name,
			Description: act.Description,
			Tags:        act.Tags,
			Mods:        []entities.ActionMod{},
		}
		for i, mod := range act.Mods {
			actTypeStr, ok := mod["type"].(string)
			if !ok {
				return &ActionTypeError{Action: act.Name, Mod: i}
			}
			actImpl, ok := gm.actionTypes[actTypeStr]
			if !ok {
				return &ActionTypeError{Action: act.Name, Mod: i}
			}
			actType := reflect.TypeOf(actImpl)
			actInstance := reflect.New(actType).Elem().Interface().(entities.ActionMod)
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Squash:  true,
				TagName: "json",
				Result:  &actInstance,
			})
			if err != nil {
				return err
			}

			err = decoder.Decode(mod)
			if err != nil {
				return err
			}
			actData.Mods = append(actData.Mods, actInstance)
		}
		gm.actions = append(gm.actions, actData)
	}

	return nil
}

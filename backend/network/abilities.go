package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SpecWrapper struct {
	Type    string                  `json:"type"`
	Format  string                  `json:"format"`
	Version string                  `json:"version"`
	Data    map[string]SpecChampion `json:"data"`
}

type SpecChampion struct {
	Spells []Ability `json:"spells"`
}

type Ability struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetAbilities(id string) []string {
	res, err := http.Get("https://ddragon.leagueoflegends.com/cdn/14.2.1/data/en_US/champion/" + id + ".json")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil
	}

	var wrapper SpecWrapper
	abilities := []string{}
	json.Unmarshal(body, &wrapper)
	for _, champ := range wrapper.Data {
		for _, ability := range champ.Spells {
			if strings.Contains(ability.Name, "/") {
				abilities = append(abilities, strings.Split(ability.Name, " / ")...)
			}
			abilities = append(abilities, ability.Name)
		}
	}

	return abilities
}

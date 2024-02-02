package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Champion struct {
	Version string `json:"version"`
	Id      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Title   string `json:"title"`
}

type Wrapper struct {
	Type    string              `json:"type"`
	Format  string              `json:"format"`
	Version string              `json:"version"`
	Data    map[string]Champion `json:"data"`
}

func GetChampions() *Wrapper {
	data, err := http.Get("https://ddragon.leagueoflegends.com/cdn/14.2.1/data/en_US/champion.json")
	if err != nil {
		fmt.Printf("Error %s", err)
		return nil
	}

	body, err := ioutil.ReadAll(data.Body)

	if err != nil {
		fmt.Printf("Error %s", err)
		return nil
	}

	var wrapper Wrapper
	json.Unmarshal(body, &wrapper)
	return &wrapper
}

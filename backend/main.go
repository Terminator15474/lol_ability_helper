package main

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Terminator15474/lol_ability_helper/backend/network"
	"github.com/Terminator15474/lol_ability_helper/backend/templates"
)

func main() {
	data := network.GetChampions()
	abilities := []string{}
	for _, champ := range data.Data {
		abilities = append(abilities, network.GetAbilities(champ.Id)...)
	}
	// `start[a-zA-Z ]*?end$`gmis
	fmt.Println(data.Data["Aatrox"])
	http.HandleFunc("/abilities", func(res http.ResponseWriter, req *http.Request) {

		start := req.URL.Query().Get("start")
		end := req.URL.Query().Get("end")
		contains := req.URL.Query().Get("contains")
		not := req.URL.Query().Get("not")
		len_str := req.URL.Query().Get("length")
		var length int64
		if len(len_str) > 0 {
			temp, err := strconv.ParseInt(len_str, 10, 64)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			length = temp
		} else {
			length = -1
		}

		var filtered = []string{}
		for _, elem := range abilities {
			replaced := strings.NewReplacer("!", "", "?", "").Replace(strings.ToLower(elem))
			if length != -1 && len(replaced) != int(length) {
				continue
			}

			if filter(replaced, start, contains, end, not) {
				filtered = append(filtered, elem)
			}
		}
		slices.Sort(filtered)
		templates.Abilities(filtered).Render(req.Context(), res)
	})
	fmt.Println("Listening on :8080")
	http.ListenAndServe("localhost:8080", nil)
}

func filter(s string, start string, contains string, end string, not string) bool {
	if len(not) > 0 {
		return strings.HasPrefix(s, start) && strings.Contains(s, contains) && strings.HasSuffix(s, end) && !strings.Contains(s, not)
	}
	return strings.HasPrefix(s, start) && strings.Contains(s, contains) && strings.HasSuffix(s, end)
}

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
	champions := network.GetChampions()
	championNames := []string{}
	for _, name := range champions.Data {
		championNames = append(championNames, name.Name)
	}

	abilityNames := []string{}
	for _, champ := range champions.Data {
		abilityNames = append(abilityNames, network.GetAbilities(champ.Id)...)
	}

	http.HandleFunc("/abilities", func(res http.ResponseWriter, req *http.Request) {
		setHeaders(res)

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
		for _, elem := range abilityNames {
			replaced := strings.NewReplacer("!", "", "?", "", " ", "").Replace(strings.ToLower(elem))
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

	http.HandleFunc("/champions", func(res http.ResponseWriter, req *http.Request) {
		setHeaders(res)

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
		for _, elem := range championNames {
			replaced := strings.NewReplacer("!", "", "?", "", " ", "").Replace(strings.ToLower(elem))
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
	ret := strings.HasPrefix(s, start) && strings.HasSuffix(s, end)

	if len(contains) > 0 {
		ret = ret && strings.ContainsAny(s, contains)
	}

	if len(not) > 0 {
		ret = ret && !strings.ContainsAny(s, not)
	}
	return ret
}

func setHeaders(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "text/html")
	res.Header().Set("Access-Control-Allow-Methods", "*")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
	res.Header().Set("Access-Control-Allow-Headers", "*")
}

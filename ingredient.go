package main

import (
	"encoding/json"
	"log"
	"os"
)

type Ingredient struct {
	ServingSize int
	Calories    int
	Protein     int
	Fiber       int
	Fat         int
}

func readJsonIngredients(filename string) map[string]Ingredient {
	ings := make(map[string]Ingredient)
	contents, err := os.ReadFile(filename)
	//TODO: Make this error handling a little nicer, but don't worry too much about it,
	//since the JSON will be switched to SQLite at some point
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			return ings
		}
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &ings)
	if err != nil {
		log.Fatal(err)
	}
	return ings
}

func writeJsonIngredients(filename string, ingredients map[string]Ingredient) {
	ingJson, err := json.Marshal(ingredients)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Rename(filename, filename+".bak")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filename, ingJson, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

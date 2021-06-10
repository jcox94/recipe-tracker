package main

import (
	"encoding/json"
	"log"
	"os"
)

type IngredientStub struct {
	Name   string
	Amount int //Start with just grams for measurement
}

type Recipe struct {
	Name  string
	Stubs []IngredientStub
}

func readJsonRecipes(filename string) []Recipe {
	r := []Recipe{}
	contents, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			return r
		}
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func writeJsonRecipes(filename string, recipes []Recipe) {
	ingJson, err := json.Marshal(recipes)
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

type RecipeTotals struct {
	Calories int
	Protein  int
	Fiber    int
	Fat      int
}

func totalNutrition(ings map[string]Ingredient, recipe Recipe) RecipeTotals {
	total := RecipeTotals{}
	for _, stub := range recipe.Stubs {
		ing, ok := ings[stub.Name]
		//TODO: Instead of crashing here, prompt the user the enter in ingredient information
		if !ok {
			log.Fatal("Could not find ingredient from stub")
		}
		//Some stuff has no nutrition info, like seasonings, so
		//this skips those. TODO: have a type for seasonings
		//(and other things where nutrition info is basically 0 for everyting)
		if ing.ServingSize == 0 {
			continue
		}
		grams := stub.Amount / ing.ServingSize
		total.Calories += grams * ing.Calories
		total.Protein += grams * ing.Protein
		total.Fiber += grams * ing.Fiber
		total.Fat += grams * ing.Fat
	}
	return total
}

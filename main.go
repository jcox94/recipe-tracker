package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/style"
)

func main() {
	state := initState()
	wnd := nucular.NewMasterWindow(0, "Recipe Tracker", state.updatefn)
	wnd.SetStyle(style.FromTheme(style.WhiteTheme, 2.0))
	wnd.Main()
}

func initState() UIState {
	state := UIState{
		ingredients: readJsonIngredients("ingredients.json"),
		allRecipes:  readJsonRecipes("recipes.json"),
	}
	state.focusedRecipes = make([]Recipe, len(state.allRecipes))
	state.searchBox.Flags = nucular.EditSimple
	state.nameBox.Flags = nucular.EditSimple
	state.recipeNameBox.Flags = nucular.EditSimple
	state.gramsBox.Flags = nucular.EditSimple
	state.servingSizeBox.Flags = nucular.EditSimple
	state.caloriesBox.Flags = nucular.EditSimple
	state.proteinBox.Flags = nucular.EditSimple
	state.fiberBox.Flags = nucular.EditSimple
	state.fatBox.Flags = nucular.EditSimple
	return state
}

//Any data that needs to persist across frames should go here, as Nucular does not
//maintain any state when it redraws the UI. One big struct seems fine for now,
//but it might make sense to break it up in the future if it grows.
type UIState struct {
	ingredients      map[string]Ingredient
	allRecipes       []Recipe
	focusedRecipes   []Recipe
	selectedRecipe   Recipe
	searchBox        nucular.TextEditor
	nameBox          nucular.TextEditor
	recipeNameBox    nucular.TextEditor
	directionsBox    nucular.TextEditor
	Box              nucular.TextEditor
	gramsBox         nucular.TextEditor
	servingSizeBox   nucular.TextEditor
	caloriesBox      nucular.TextEditor
	proteinBox       nucular.TextEditor
	fiberBox         nucular.TextEditor
	fatBox           nucular.TextEditor
	addingRecipe     bool
	addingIngredient bool
	stubs            []IngredientStub
	recipeName       string
}

func (stub IngredientStub) String() string {
	return fmt.Sprintf("%d grams of %s", stub.Amount, stub.Name)
}

func (total RecipeTotals) String() string {
	return fmt.Sprintf("Calories: %d\nProtein: %d\nFiber: %d\nFat: %d", total.Calories, total.Protein, total.Fiber, total.Fat)
}

func (state *UIState) updatefn(w *nucular.Window) {
	w.Row(25).Dynamic(4)
	state.searchBox.Edit(w)
	if w.Button(label.T("Add New Recipe"), false) {
		state.addingRecipe = true
	}

	if state.addingRecipe {
		w.Label("Recipe Name:", "LC")
		state.recipeNameBox.Edit(w)
	}

	state.focusedRecipes = state.focusedRecipes[:0]
	for _, r := range state.allRecipes {
		if strings.Contains(r.Name, string(state.searchBox.Buffer)) {
			state.focusedRecipes = append(state.focusedRecipes, r)
		}
	}

	w.Row(0).Dynamic(2)
	if gl, w := nucular.GroupListStart(w, len(state.focusedRecipes), "recipes", nucular.WindowDefaultFlags); w != nil {
		w.Row(20).Dynamic(1)
		for gl.Next() {
			selected := state.focusedRecipes[gl.Index()].Name == state.selectedRecipe.Name
			if w.SelectableLabel(state.focusedRecipes[gl.Index()].Name, "LC", &selected) {
				state.selectedRecipe = state.focusedRecipes[gl.Index()]
			}
		}
	}
	if state.addingRecipe {
		state.addRecipe(w)
	} else {
		if group := w.GroupBegin("Display Recipe", nucular.WindowNoScrollbar); group != nil {
			if state.selectedRecipe.Name != "" {
				group.Row(40).Dynamic(1)
				group.Label(state.selectedRecipe.Name, "LT")
				for _, stub := range state.selectedRecipe.Stubs {
					group.Row(15).Dynamic(1)
					group.Label(stub.String(), "LC")
				}
				total := totalNutrition(state.ingredients, state.selectedRecipe)
				group.Row(80).Dynamic(1)
				group.Label(total.String(), "LC")
			}
			group.GroupEnd()
		}
	}
}

func validInputBoxes(boxes ...*nucular.TextEditor) bool {
	for _, box := range boxes {
		if _, err := strconv.Atoi(string(box.Buffer)); err != nil {
			return false
		}
	}
	return true
}

func (state *UIState) addRecipe(w *nucular.Window) {
	if group := w.GroupBegin("Add Recipe", nucular.WindowNoScrollbar); group != nil {
		if state.addingIngredient {
			state.addIngredient(group)
		} else {
			group.Row(25).Dynamic(2)
			if group.Button(label.T("Cancel"), false) {
				state.addingRecipe = false
			}
			if string(state.recipeNameBox.Buffer) == "" {
				group.Label("Recipe Name Required", "CC")
			} else if group.Button(label.T("Submit Recipe"), false) {
				state.nameBox.Delete(0, len(state.nameBox.Buffer))
				state.gramsBox.Delete(0, len(state.gramsBox.Buffer))
				state.addingRecipe = false
				recipe := Recipe{
					Name:  string(state.recipeNameBox.Buffer),
					Stubs: state.stubs,
				}
				state.allRecipes = append(state.allRecipes, recipe)
				writeJsonRecipes("recipes.json", state.allRecipes)
				state.stubs = make([]IngredientStub, 0)
			}
			group.Row(25).Dynamic(2)
			group.Label("Ingedient Name:", "LC")
			state.nameBox.Edit(group)
			group.Row(25).Dynamic(2)
			group.Label("Grams:", "LC")
			state.gramsBox.Edit(group)
			group.Row(25).Dynamic(1)
			if string(state.nameBox.Buffer) == "" || !validInputBoxes(&state.gramsBox) {
				group.Label("Enter valid name and grams", "CC")
			} else if group.Button(label.T("Add Ingredient"), false) {
				//Error cannot occur here as input must be valid for "Add Ingredient" button to appear
				amount, _ := strconv.Atoi(string(state.gramsBox.Buffer))
				ingName := string(state.nameBox.Buffer)
				ingName = strings.ToLower(ingName)
				stub := IngredientStub{
					Name:   ingName,
					Amount: amount,
				}
				_, ok := state.ingredients[stub.Name]
				if !ok {
					state.addingIngredient = true
				} else {
					state.stubs = append(state.stubs, stub)
					state.nameBox.Delete(0, len(state.nameBox.Buffer))
					state.gramsBox.Delete(0, len(state.gramsBox.Buffer))
				}
			}
		}
		group.GroupEnd()
	}
}

func (state *UIState) addIngredient(group *nucular.Window) {
	group.Row(25).Dynamic(2)
	group.Label("Serving Size:", "LC")
	state.servingSizeBox.Edit(group)
	group.Row(25).Dynamic(2)
	group.Label("Calories:", "LC")
	state.caloriesBox.Edit(group)
	group.Row(25).Dynamic(2)
	group.Label("Protein:", "LC")
	state.proteinBox.Edit(group)
	group.Row(25).Dynamic(2)
	group.Label("Fiber", "LC")
	state.fiberBox.Edit(group)
	group.Row(25).Dynamic(2)
	group.Label("Fat", "LC")
	state.fatBox.Edit(group)
	group.Row(25).Dynamic(2)
	if group.Button(label.T("Cancel"), false) {
		state.addingIngredient = false
	}
	valid := validInputBoxes(&state.servingSizeBox, &state.caloriesBox, &state.proteinBox, &state.fiberBox, &state.fatBox)
	if !valid {
		group.Label("Info not valid", "LC")
	} else if group.Button(label.T("Submit"), false) {
		name := string(state.nameBox.Buffer)
		name = strings.ToLower(name)
		//Errors don't need to be checked here, as all inputs have to be valid before the Submit button can be clicked
		servingSize, _ := strconv.Atoi(string(state.servingSizeBox.Buffer))
		calories, _ := strconv.Atoi(string(state.caloriesBox.Buffer))
		protein, _ := strconv.Atoi(string(state.proteinBox.Buffer))
		fiber, _ := strconv.Atoi(string(state.fiberBox.Buffer))
		fat, _ := strconv.Atoi(string(state.fatBox.Buffer))
		ing := Ingredient{
			ServingSize: servingSize,
			Calories:    calories,
			Protein:     protein,
			Fiber:       fiber,
			Fat:         fat,
		}
		state.ingredients[name] = ing
		state.addingIngredient = false
		writeJsonIngredients("ingredients.json", state.ingredients)
	}
}

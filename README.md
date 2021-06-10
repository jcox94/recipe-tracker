# Recipe Tracker (I should come up with a cooler name)

A (work in progress!) gui-based program for keeping track of recipes and their nutritional info.

## Overview

The bulk of the data consists of ingredients, which hold their own nutritional info, and recipes, which
hold a list of ingredients and how much of each ingredient to use. When adding a new recipe, if 
an ingredient which has not been added before is referenced, the user is prompted to enter the
information for that ingredient. Later, when the user selects a recipe to view, the total nutritional value
for that recipe is calculated by adding up the values for each ingredient.

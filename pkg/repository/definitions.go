package repository

import (
	"github.com/jinzhu/gorm"
)

/*
	This is the table definitions
	#################################
*/

// RecipeTable is the base table for a recipe
type recipe struct {
	gorm.Model
	Name        string
	Description string
	Category    string
	Cuisine     string
	Yield       int
	FoundOn     string
}

// CategoryTable is self-explanatory
type category struct {
	gorm.Model
	category string
}

// CuisineTable is self-explanatory
type cuisine struct {
	gorm.Model
	name string
}

// RecipeCategoryTable is a connection table
type recipeCategory struct {
	RecipeID   uint
	CategoryID uint
	Recipe     recipe
	Category   category
}

// RecipeCuisineTable is a connection table
type recipeCuisine struct {
	RecipeID  uint
	CuisineID uint
	Recipe    recipe
	Cuisine   cuisine
}

// IngredientTable contains ingredient with a reference to a certain recipe
type recipeIngredient struct {
	gorm.Model
	// TODO: change this column name
	Ingredient string
	RecipeID   uint
	Recipe     recipe
}

// KeywordTable contains keywords with a reference to a certain recipe
type recipeKeyword struct {
	gorm.Model
	Keyword  string
	RecipeID uint
	Recipe   recipe
}

// RecipeImageTable is images of a certain recipe
type recipeImage struct {
	gorm.Model
	Image    string
	RecipeID uint
	Recipe   recipe
}

// RecipeInstructionTable is instructions for a certain recipe
type recipeInstruction struct {
	gorm.Model
	Step     int
	Text     string
	RecipeID uint
	Recipe   recipe
}

// RatingTable is the score and number of votes a recipe has
type rating struct {
	gorm.Model
	Votes    string
	Score    string
	RecipeID uint
	Recipe   recipe
}

// UserTable ...
type user struct {
	gorm.Model
	Email string `gorm:"not null;unique_index"`
	Name  string
}

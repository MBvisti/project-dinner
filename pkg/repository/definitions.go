package repository

import "github.com/jinzhu/gorm"

/*
	This is the table definitions
	#################################
*/

// RecipeTable is the base table for a recipe
type RecipeTable struct {
	gorm.Model
	Name        string
	Type        string
	Description string
	Category    string
	Cuisine     string
	Yield       string
}

// CategoryTable is self-explanatory
type CategoryTable struct {
	gorm.Model
	category string
}

// CuisineTable is self-explanatory
type CuisineTable struct {
	gorm.Model
	name string
}

// RecipeCategoryTable is a connection table
type RecipeCategoryTable struct {
	RecipeID   uint
	CategoryID uint
	Recipe     RecipeTable
	Category   CategoryTable
}

// RecipeCuisineTable is a connection table
type RecipeCuisineTable struct {
	RecipeID  uint
	CuisineID uint
	Recipe    RecipeTable
	Cuisine   CuisineTable
}

// IngredientTable contains ingredient with a reference to a certain recipe
type IngredientTable struct {
	gorm.Model
	// TODO: change this column name
	Ingredient string
	RecipeID   uint
	Recipe     RecipeTable
}

// KeywordTable contains keywords with a reference to a certain recipe
type KeywordTable struct {
	gorm.Model
	Keyword  string
	RecipeID uint
	Recipe   RecipeTable
}

// RecipeImageTable is images of a certain recipe
type RecipeImageTable struct {
	gorm.Model
	Image    string
	RecipeID uint
	Recipe   RecipeTable
}

// RecipeInstructionTable is instructions for a certain recipe
type RecipeInstructionTable struct {
	gorm.Model
	Step     int
	Text     string
	RecipeID uint
	Recipe   RecipeTable
}

// RatingTable is the score and number of votes a recipe has
type RatingTable struct {
	gorm.Model
	Votes    string
	Score    string
	RecipeID uint
	Recipe   RecipeTable
}

// UserTable ...
type UserTable struct {
	gorm.Model
	Email string `gorm:"not null;unique_index"`
	Name  string
}

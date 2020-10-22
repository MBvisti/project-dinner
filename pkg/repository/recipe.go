package repository

import (
	"html/template"

	"github.com/jinzhu/gorm"
)

// EmailRecipe is a shorter version of a recipe only intended for emails
type EmailRecipe struct {
	Name        string
	Description string
	Category    string
	Cuisine     string
	ThumbNail   string
	FoundOn     string
}

// Recipe is for a whole recipe
type Recipe struct {
	Name         string
	Description  string
	Category     string
	Cuisine      string
	Yield        int
	Ingredients  []string
	Keywords     []string
	Images       []string
	Instructions []Instruction
	Score        Rating
	FoundOn      string
}

// Instruction ...
type Instruction struct {
	Step int
	Text string
}

// Rating ...
type Rating struct {
	Votes string
	Score string
}

// DailyRecipes ... TODO: update this table
type DailyRecipes struct {
	gorm.Model
	Name         string        `json:"title"`
	Image        string        `json:"image"`
	Description  template.HTML `json:"summary"`
	Source       string        `json:"sourceUrl"`
	Instructions string        `json:"instructions"`
}

package repository

import (
	"html/template"

	"github.com/jinzhu/gorm"
)

// Recipe ...
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

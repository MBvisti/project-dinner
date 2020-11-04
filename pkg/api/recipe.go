package api

type FeaturedRecipe struct {
	Image string `json:"image"`
	Name  string `json:"name"`
}

type EmailRecipe struct {
	Name        string
	Description string
	Category    string
	Cuisine     string
	ThumbNail   string
	FoundOn     string
}

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

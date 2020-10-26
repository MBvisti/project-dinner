package api

// BaseSection is the basic recipe structure in ld+json format
type BaseSection struct {
	Context string          `json:"@context"`
	Graph   []RecipeSection `json:"@graph"`
}

type RecipeSection struct {
	Type               string               `json:"@type"`
	Name               string               `json:"name"`
	ID                 string               `json:"@id"`
	AggregatedRating   RatingSection        `json:"aggregateRating"`
	Description        string               `json:"description"`
	RecipeIngredients  []string             `json:"recipeIngredient"`
	Image              []string             `json:"image"`
	Nutrition          map[string]string    `json:"nutrition"`
	Keywords           string               `json:"keywords"`
	RecipeCategory     []string             `json:"recipeCategory"`
	RecipeCuisine      []string             `json:"recipeCuisine"`
	RecipeInstructions []InstructionSection `json:"recipeInstructions"`
	RecipeYield        []string             `json:"recipeYield"`
	FoundOn            string
}

// InstructionSection is a scraped recipe's instructions defined in ld+json format
type InstructionSection struct {
	Type            string            `json:"@type"`
	Name            string            `json:"name,omitempty"`
	URL             string            `json:"url,omitempty"`
	Text            string            `json:"text,omitempty"`
	ItemListElement []ItemListElement `json:"itemListElement,omitempty"`
}

// ItemListElement is used when present and is a way to list instructions for a recipe
type ItemListElement struct {
	Type string `json:"@type"`
	Text string `json:"text"`
}

// RatingSection are the rating a given recipe have
type RatingSection struct {
	RatingCount string `json:"ratingCount"`
	RatingValue string `json:"ratingValue"`
}

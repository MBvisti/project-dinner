package api

// User ....
type User struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	RecipeType  string `json:"recipe_type"`
	DietaryType string `json:"dietary_type"`
}

//type Recipe struct {
//	Name string
//}

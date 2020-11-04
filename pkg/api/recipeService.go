package api

// RecipeService is the recipe service's interface
type RecipeService interface {
	GetFeaturedRecipes() []FeaturedRecipe
}

// RecipeRepository ...
type RecipeRepository interface {
	GetFeaturedRecipes() []FeaturedRecipe
}

type recipeService struct {
	storage RecipeRepository
}

// NewRecipeService ...
func NewRecipeService(r RecipeRepository) RecipeService {
	return &recipeService{
		r,
	}
}

func (s *recipeService) GetFeaturedRecipes() []FeaturedRecipe {
	fR := s.storage.GetFeaturedRecipes()

	return fR
}

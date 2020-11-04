package repository

import (
	"project-dinner/pkg/api"
	"regexp"
)

// GetEmailList returns the email list - TODO: maybe separate the user list and emailing list into two different tables?
func (r *repoService) GetEmailList() ([]api.User, error) {
	var emailList []api.User
	err := r.db.Table("users").Select("email, name").Scan(&emailList).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}

	return emailList, nil
}

var isMailValid = regexp.MustCompile(
	`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)

// CreateUser ...
func (r *repoService) CreateUser(usr api.User) error {
	if usr.Email == "" {
		return ErrEmailRequired
	}

	if !isMailValid.MatchString(usr.Email) {
		return ErrEmailInvalid
	}

	var recipeTypeID recipeType
	err := r.db.Table("recipe_types").Where("type = ?", usr.RecipeType).Select("id").Scan(&recipeTypeID).Error
	if err != nil {
		return ErrNoResourceFound
	}

	var dietaryTypeID dietaryType
	err = r.db.Table("dietary_types").Where("type = ?", usr.DietaryType).Select("id").Scan(&dietaryTypeID).Error
	if err != nil {
		return ErrNoResourceFound
	}

	nU := user{
		Email:         usr.Email,
		Name:          usr.Name,
		RecipeTypeID:  recipeTypeID.ID,
		DietaryTypeID: dietaryTypeID.ID,
	}

	err = r.db.Create(&nU).Error

	if err != nil {
		return ErrNoCreate
	}

	return nil
}

package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repository struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Email string `gorm:"not null;unique_index"`
	Name  string
}

type Recipe struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Image       string
	Link        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Source      string
	Rating      string
	Review      string
}

type EmailList struct {
	Email string
	Name  string
}

func NewRepository(db *gorm.DB) *Repository {

	return &Repository{
		db: db,
	}
}

func (r *Repository) DestructiveReset() error {
	err := r.db.DropTableIfExists(&User{}, &Recipe{}).Error
	if err != nil {
		return err
	}

	// TODO: this goes here for now until a better option is needed
	var recipes = []Recipe{
		{
			Name:        "Zucchini Noodles with Avocado Pesto & Shrimp",
			Image:       "https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fimages.media-allrecipes.com%2Fuserphotos%2F4465921.jpg",
			Link:        "http://www.eatingwell.com/recipe/257004/zucchini-noodles-with-avocado-pesto-shrimp/",
			Description: "Cut some carbs and use spiralized zucchini in place of noodles in this zesty pesto pasta dish recipe. Top with Cajun-seasoned shrimp to complete this quick and easy dinner.",
			Source:      "EatingWell.com",
			Rating:      "4/5",
			Review:      "Delicious. I took a shortcut and used 4 oz of prepackaged pesto. Followed the recipes otherwise. Only took minutes to throw together....awesome!",
		},
		{
			Name:        "World's Best Lasagna",
			Image:       "https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fimages.media-allrecipes.com%2Fuserphotos%2F3359675.jpg",
			Link:        "https://www.allrecipes.com/recipe/23600/worlds-best-lasagna/",
			Description: "It takes a little work, but it is worth it.",
			Source:      "https://www.allrecipes.com/",
			Rating:      "5/5",
			Review:      "I made this recipe using some tips and tricks from previous reviewers and from my experience having worked in an italian restaurant.  This recipe does not have to be as time intensive as is recommended in the instructions.  As another reviewer noted, it is not necessary to cook the noodles--just soak them in hot water while you are cooking the rest of the ingredients.  You also do not need to cook the sauce.  Sounds bizarre, but it's true.  If you make the lasagna ahead and leave it in the fridge overnight, the flavors blend perfectly, no cooking required.  For a chunkier tomato sauce, I subbed 1/2 of the crushed tomatoes with diced italian-flavored tomatoes.  I layered everything in this order: sauce-noodles-ricotta mixture-sauce-meat-shredded mozzerella...i did 3 layers like this and then finished with a layer of noodles, sauce and shredded mozzerella and a generous sprinkling of parmesean on top.  Following this method, I had the whole thing put together and in the fridge in less than an hour.  It sat overnight and my husband had it in the oven when I got home from work.  Fantastic recipe!",
		},
		{
			Name:        "Oyakodon (Japanese Chicken and Egg Rice Bowl)",
			Image:       "https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fimages.media-allrecipes.com%2Fuserphotos%2F3963616.jpg",
			Link:        "https://www.allrecipes.com/recipe/128589/oyakodon-japanese-chicken-and-egg-rice-bowl/",
			Description: "This is a delicious traditional Japanese meal consisting of chicken sauteed and then cooked in a Japanese broth, and then finished with egg and served over rice.  It's really easy, filling and delicious.",
			Source:      "https://www.allrecipes.com/",
			Rating:      "4/5",
			Review:      "This is the ORIGINAL Oyakodon recipe! If you can't find dashi, try putting fish stock instead of chicken. Dashi are little dryed fish turned into powder...!  By the way, did you know that \"Oyako\" means \"Parents and Children\" (in this case, chicken and eggs); don is just the word used for rice bowl ;)",
		},
		{
			Name:        "Low-Carb Tacos",
			Image:       "https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fimages.media-allrecipes.com%2Fuserphotos%2F3884031.jpg",
			Link:        "https://www.allrecipes.com/recipe/239831/low-carb-tacos/",
			Description: "This is a great low-carb alternative to your standard homemade tacos. I love Mexican food and wasn't willing to part with tacos after starting my low-carb diet. This always satisfies my craving.",
			Source:      "https://www.allrecipes.com/",
			Rating:      "4/5",
			Review:      "Delicious! I used a taco seasoning recipe from this site and did everything else as instructed except I used fresh jalapeños because I didn't have canned. I was wondering if I was supposed to drain the meat and decided not to since the recipe didn't say to and it was so good. Thank you I will be making this again for sure.",
		},
	}

	err = r.AutoMigrate()
	if err != nil {
		return err
	}

	for _, recipe := range recipes {
		err = r.db.Create(&recipe).Error

		if err != nil {
			return err
		}
	}

	user := User{
		Email: "mbv1406@gmail.com",
		Name:  "Morten",
	}

	err = r.db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateRecipe(recipe *Recipe) error {
	err := r.db.Create(&recipe).Error

	return err
}

func (r *Repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(&User{}, &Recipe{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetEmailList() ([]EmailList, error) {
	var emailList []EmailList
	err := r.db.Table("users").Select("email, name").Scan(&emailList).Error

	if err != nil {
		return nil, err
	}

	return emailList, nil
}

//func randInt(min int, max int) int {
//	return min + rand.Intn(max-min)
//}

func (r *Repository) TodaysRecipes() ([]Recipe, error) {
	var recipes []Recipe
	var count int
	err := r.db.Model(&recipes).Count(&count).Error
	if err != nil {
		return nil, err
	}

	var selectedRecipes []Recipe
	err = r.db.Find(&selectedRecipes, []int{1, 2, 3, 4}).Error

	if err != nil {
		return nil, err
	}

	return selectedRecipes, nil
}

package api_test

import (
	"gopkg.in/gomail.v2"
	"project-dinner/pkg/api"
	"testing"
)

type mockEmailRepo struct{}

func (mE *mockEmailRepo) GetDailyRecipes() ([]api.EmailRecipe, error) {
	return nil, nil
}
func (mE *mockEmailRepo) GetEmailList() ([]api.User, error) {
	return nil, nil
}

func TestCreateWelcomeMail(t *testing.T) {
	mockEmailR := mockEmailRepo{}
	t.Run("should create and return a valid welcome mail", func(t *testing.T) {
		mockMailProvider := gomail.NewDialer("example.smtp.com", 25, "username", "password")
		mockEmailService := api.NewEmailService(mockMailProvider, &mockEmailR)

		u := api.User{
			Name:  "Testesen",
			Email: "testesen@live.dk",
		}
		msg, err := mockEmailService.CreateWelcomeMail(u)

		if err != nil {
			t.Errorf("Create welcome mail test failed. Got err: %s", err.Error())
		}

		expectedFrom := []string{"\"Morten's recipe service\" <noreply@mbvistisen.dk>"}
		expectedTo := []string{u.Email}
		expectedSubject := []string{"Thanks for signing up!"}

		from := msg.GetHeader("From")
		to := msg.GetHeader("To")
		subject := msg.GetHeader("Subject")

		if from[0] != expectedFrom[0] {
			t.Errorf("Create welcome mail test failed. expected: %s , got: %s", expectedFrom[0], from[0])
		}
		if to[0] != expectedTo[0] {
			t.Errorf("Create welcome mail test failed. expected: %s , got: %s", expectedTo[0], to[0])
		}
		if subject[0]  != expectedSubject[0] {
			t.Errorf("Create welcome mail test failed. expected: %s , got: %s", expectedSubject[0], subject[0])
		}
	})
}

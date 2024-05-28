package auth

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"mpt_data/test/vars"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("succesfull", func(t *testing.T) {
		// Prepare
		// Act
		_, err := Login(vars.UserAPI)
		// Assert
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		// Prepare
		// Act
		_, err := Login(apiModel.UserLogin{})
		// Assert
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}
func TestHash(t *testing.T) {
	t.Run("succesfull", func(t *testing.T) {
		// Prepare
		pw := "mySupersafePassword"
		// Act
		hash, err := hash(pw)
		// Assert
		if err != nil {
			t.Errorf("expexted no error, got %s", err)
		}

		if hash == nil {
			t.Errorf("expexted data, got empty slice")
		}
	})
}

func TestValidateUser(t *testing.T) {
	// Prepare
	user := apiModel.UserLogin{Username: "tester", Password: "test"}
	if err := CreateUser(user); err != nil {
		t.Skip("error preparing test")
	}
	t.Cleanup(func() { database.DB.Unscoped().Delete(&dbModel.User{}, "username = ?", user.Username) })
	var testcases = []struct {
		name string
		user apiModel.UserLogin
		err  error
	}{
		{"succesful", user, nil},
		{"no user data", apiModel.UserLogin{}, errors.ErrInvalidAuth},
		{"invalid user data", apiModel.UserLogin{Username: "bla", Password: "bla"}, errors.ErrInvalidAuth},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			_, err := validateUser(testcase.user)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

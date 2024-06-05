package auth

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/helper/config"
	"mpt_data/helper/errors"
	apiModel "mpt_data/models/apimodel"
	dbModel "mpt_data/models/dbmodel"
	"mpt_data/test/vars"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := database.Connect(vars.GetDbPAth()); err != nil {
		panic(err)
	}
	if err := os.Chdir("../.."); err != nil {
		fmt.Println(err)
	}
	// Load the config
	config.LoadConfig()
	m.Run()
}

func TestCreateUser(t *testing.T) {
	user := apiModel.UserLogin{Username: "M_Maier", Password: "TestingPW"}
	t.Cleanup(func() {
		username, err := (&(dbModel.User{Username: user.Username})).EncryptedUsername()
		if err != nil {
			t.Skipf("test cleanup failed: %v", err)
		}
		database.DB.Unscoped().Delete(&dbModel.User{}, "username = ?", username)
	})
	var testcases = []struct {
		name string
		user apiModel.UserLogin
		err  error
	}{
		{"succesfull", user, nil},
		{"duplicate user", user, errors.ErrUserAlreadyExists},
		{"not enough values", apiModel.UserLogin{}, errors.ErrUserNotComplete},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			// Act
			err := CreateUser(testcase.user)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}
func TestAddUser(t *testing.T) {
	user := dbModel.User{Username: "M_Maier", Hash: "Test"}
	t.Cleanup(func() {
		username, err := user.EncryptedUsername()
		if err != nil {
			t.Skipf("test cleanup failed: %v", err)
		}
		database.DB.Unscoped().Delete(&dbModel.User{}, "username = ?", username)
	})

	var testcases = []struct {
		name string
		user dbModel.User
		err  error
	}{
		{"succesfull", user, nil},
		{"duplicate user", user, errors.ErrUserAlreadyExists},
		{"not enough values", dbModel.User{}, errors.ErrUserdataNotComplete},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			// Act
			err := addUser(testcase.user)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

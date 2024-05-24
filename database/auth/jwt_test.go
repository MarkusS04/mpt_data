package auth

import (
	"mpt_data/database"
	dbModel "mpt_data/models/dbmodel"
	"mpt_data/test/vars"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// Prepare
		var user dbModel.User
		database.DB.First(&user, "username = ?", vars.UserAPI.Username)
		// Act
		_, err := generateJWT(user)
		// Assert
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	})
}

func TestValidateJwt(t *testing.T) {
	t.Run("succesfull", func(t *testing.T) {
		// Prepare
		var user dbModel.User
		database.DB.First(&user, "username = ?", vars.UserAPI.Username)
		jwt, _ := generateJWT(user)
		// Act
		_, err := ValidateJWT(jwt)
		// Assert
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	})
}

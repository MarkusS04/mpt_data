package person

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	"mpt_data/models/dbmodel"
	dbModel "mpt_data/models/dbmodel"
	"mpt_data/test/vars"
	"testing"
)

func TestMain(m *testing.M) {
	vars.PrepareConfig()
	m.Run()
}

func TestAddPerson(t *testing.T) {
	var testcases = []struct {
		name   string
		person *dbModel.Person
		err    error
	}{
		{"succesfull", &dbModel.Person{GivenName: "Max", LastName: "Maier"}, nil},
		{"no first name", &dbModel.Person{GivenName: "", LastName: "Maier"}, errors.ErrPersonMissingName},
		{"no last name", &dbModel.Person{GivenName: "Max", LastName: ""}, errors.ErrPersonMissingName},
		{"no name", &dbModel.Person{GivenName: "", LastName: ""}, errors.ErrPersonMissingName},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			tx := database.DB.Begin()
			defer tx.Rollback()
			// Act
			err := AddPerson(tx, testcase.person)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	var testcases = []struct {
		name         string
		personOrig   *dbModel.Person
		personUpdate *dbModel.Person
		err          error
	}{
		{"succesfull", &dbModel.Person{GivenName: "Max", LastName: "Maier"},
			&dbModel.Person{GivenName: "Moritz", LastName: "Muster"}, nil},
		{"no given Name", &dbModel.Person{GivenName: "Max", LastName: "Maier"},
			&dbModel.Person{GivenName: "", LastName: "Muster"}, errors.ErrPersonMissingName},
		{"no last Name", &dbModel.Person{GivenName: "Max", LastName: "Maier"},
			&dbModel.Person{GivenName: "Moritz", LastName: ""}, errors.ErrPersonMissingName},
		{"no Name", &dbModel.Person{GivenName: "Max", LastName: "Maier"},
			&dbModel.Person{GivenName: "", LastName: ""}, errors.ErrPersonMissingName},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()

			AddPerson(tx, testcase.personOrig)
			testcase.personUpdate.ID = testcase.personOrig.ID
			// Act
			err := UpdatePerson(tx, testcase.personUpdate)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if err == nil {
				var personTest dbModel.Person
				database.DB.First(&personTest, "id = ?", testcase.personUpdate.ID)
				if personTest.GivenName == testcase.personUpdate.GivenName && personTest.LastName == testcase.personUpdate.LastName {
					t.Errorf("expected givenName to be changed but didn't")
				}
			}
		})
	}

	t.Run("no id set", func(t *testing.T) {
		// Prepare
		tx := database.DB.Begin()
		defer tx.Rollback()
		person := &dbmodel.Person{
			GivenName: "Test",
			LastName:  "Tester",
		}

		AddPerson(tx, person)
		// Act
		person.ID = 0
		err := UpdatePerson(tx, person)
		// Assert
		if err != errors.ErrIDNotSet {
			t.Errorf("expected %s, got %s", errors.ErrIDNotSet, err)
		}
	})
}

func TestGetPerson(t *testing.T) {
	t.Run("succesfull", func(t *testing.T) {
		// Prepare
		tx := database.DB.Begin()
		defer tx.Rollback()

		personTest := &dbModel.Person{GivenName: "Max", LastName: "Maier"}
		AddPerson(tx, personTest)
		// Act
		persons, err := GetPerson(tx)
		// Assert
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
		personFound := false
		for _, person := range persons {
			if person.GivenName == personTest.GivenName && person.LastName == personTest.LastName {
				personFound = true
				break
			}
		}
		if !personFound {
			t.Errorf("expected %v as part of %v", personTest, persons)
		}
	})
}

func TestDeletePerson(t *testing.T) {
	var testcases = []struct {
		name   string
		person *dbModel.Person
		setID  bool
		err    error
	}{
		{"succesfull", &dbModel.Person{GivenName: "Max", LastName: "Maier"}, true, nil},
		{"no id set", &dbModel.Person{GivenName: "Max", LastName: "Maier"}, false, errors.ErrIDNotSet},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Prepare
			tx := database.DB.Begin()
			defer tx.Rollback()
			AddPerson(tx, testcase.person)
			// Act
			if !testcase.setID {
				testcase.person.ID = 0
			}
			err := DeletePerson(tx, *testcase.person)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

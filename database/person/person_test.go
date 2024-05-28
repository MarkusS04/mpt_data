package person

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/helper/config"
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

var person = &dbModel.Person{GivenName: "Max", LastName: "Maier"}

func TestAddPerson(t *testing.T) {
	var testcases = []struct {
		name   string
		person *dbModel.Person
		err    error
	}{
		{"succesfull", person, nil},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := AddPerson(testcase.person)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	var testcases = []struct {
		name   string
		person *dbModel.Person
		err    error
	}{
		{"succesfull", person, nil},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			testcase.person.GivenName = "Moritz"
			err := UpdatePerson(*testcase.person)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			var personTest dbModel.Person
			database.DB.First(&personTest, "id = ?", person.ID)
			if personTest.GivenName != testcase.person.GivenName || personTest.LastName != testcase.person.LastName {
				t.Errorf("expected givenName to be changed but didn't")
			}
		})
	}
}

func TestGetPerson(t *testing.T) {
	var testcases = []struct {
		name   string
		person dbModel.Person
		err    error
	}{
		{"succesfull", *person, nil},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			persons, err := GetPerson()
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			personFound := false
			for _, person := range persons {
				if person.GivenName == testcase.person.GivenName || person.LastName == testcase.person.LastName {
					personFound = true
					break
				}
			}
			if !personFound {
				t.Errorf("expected to get same person")
				t.Log(testcase.person)
			}
		})
	}
}

func TestDeletePerson(t *testing.T) {
	var testcases = []struct {
		name   string
		person dbModel.Person
		err    error
	}{
		{"succesfull", *person, nil},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			err := DeletePerson(testcase.person)
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
		})
	}
}

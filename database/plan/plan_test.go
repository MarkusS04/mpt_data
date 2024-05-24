package plan

import (
	"fmt"
	"mpt_data/database"
	"mpt_data/helper"
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
	helper.LoadConfig()
	m.Run()
}

func TestCreatePlan(t *testing.T) {
	var testcases = []struct {
		name string
		err  error
	}{}
	t.Cleanup(func() {})
	for _, testcase := range testcases {
		t.Run(testcase.name, func(_ *testing.T) {
			// Act
			// Assert
		})
	}
}

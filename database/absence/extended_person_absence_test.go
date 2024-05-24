package absence

import (
	"mpt_data/database"
	dbModel "mpt_data/models/dbmodel"
	database_test "mpt_data/test/database"
	"testing"
)

var (
	recurring = []*dbModel.PersonRecurringAbsence{
		{
			PersonID: 1,
			Weekday:  0,
		},
		{
			PersonID: 1,
			Weekday:  2,
		},
	}
)

func TestAddRecurringAbsence(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name    string
		absence []*dbModel.PersonRecurringAbsence
		err     error
		nums    int64
	}{
		{"success", recurring, nil, int64(len(recurring))},
	}
	for _, testcase := range testcases {
		db := database.DB
		t.Run(testcase.name, func(t *testing.T) {
			countBefore := database_test.CountEntriesDB(db, &dbModel.PersonRecurringAbsence{})
			// Act
			err := AddRecurringAbsence(testcase.absence, db)
			countAfter := database_test.CountEntriesDB(db, &dbModel.PersonRecurringAbsence{})
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if countAfter-countBefore != testcase.nums {
				t.Errorf("expected %d entries created, got %d", testcase.nums, countAfter-countBefore)
			}
		})
	}
}

func TestGetRecurringAbsence(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name    string
		absence []*dbModel.PersonRecurringAbsence
		err     error
		nums    int64
	}{
		{"success", recurring, nil, int64(len(recurring))},
	}
	for _, testcase := range testcases {
		db := database.DB.Begin()
		t.Run(testcase.name, func(t *testing.T) {
			// Act
			_, err := GetRecurringAbsence(testcase.absence[0].PersonID, db)
			count := database_test.CountEntriesDB(db, &dbModel.PersonRecurringAbsence{})
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if count != testcase.nums {
				t.Errorf("expected %d entries created, got %d", testcase.nums, count)
			}

		})
		db.Rollback()
	}
}

func TestDeleteRecurringAbsence(t *testing.T) {
	// Prepare
	var testcases = []struct {
		name    string
		absence []dbModel.PersonRecurringAbsence
		err     error
	}{
		{"success", []dbModel.PersonRecurringAbsence{*recurring[0], *recurring[1]}, nil},
	}
	for _, testcase := range testcases {
		db := database.DB.Begin()
		t.Run(testcase.name, func(t *testing.T) {
			countBefore := database_test.CountEntriesDB(db, &dbModel.PersonRecurringAbsence{})
			// Act
			err := DeleteRecurringAbsence(testcase.absence, db)
			count := database_test.CountEntriesDB(db, &dbModel.PersonRecurringAbsence{})
			// Assert
			if err != testcase.err {
				t.Errorf("expected %s, got %s", testcase.err, err)
			}
			if countBefore-count != 2 {
				t.Errorf("expected %d entries deleted, got %d", 2, countBefore-count)
			}

		})
		db.Rollback()
	}
}

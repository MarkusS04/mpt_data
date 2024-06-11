// Package plan provides functions to create, retrive, update plans and select availabe and absent people for the plan
package plan

import (
	"mpt_data/database"
	"mpt_data/helper/errors"
	"mpt_data/models/apimodel"
	"mpt_data/models/dbmodel"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const packageName = "database.plan"

// GetPlan loads all plan items in the specified Period.
// Ordered by the date of the meeting
func GetPlan(period generalmodel.Period) ([]dbModel.Plan, error) {
	var plan []dbModel.Plan
	if err :=
		database.DB.Preload("Person").
			Preload("Meeting.Tag").
			Preload("Meeting").
			Preload("TaskDetail.Task").
			Preload("TaskDetail").
			Joins("JOIN meetings m on m.id = meeting_id").
			Where("meeting_id IN (?)", database.DB.Table("meetings").Where("date between ? and ?", period.StartDate, period.EndDate).Select("id")).
			Order("m.date asc").
			Find(&plan).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

// GetPlanWithID loads the data for a specific plan item
func GetPlanWithID(planID uint) (plan dbModel.Plan, err error) {
	if err :=
		database.DB.Preload("Person").
			Preload("Meeting.Tag").
			Preload("Meeting").
			Preload("TaskDetail.Task").
			Preload("TaskDetail").
			Where("id = ?", planID).
			Find(&plan).Error; err != nil {
		return dbModel.Plan{}, err
	}
	return plan, nil
}

// CreatePlanData creates all entries in table plans for the specified period and if people are available they will be automatically assigned
func CreatePlanData(db *gorm.DB, period generalmodel.Period) ([]dbModel.Plan, error) {
	const funcName = packageName + ".CreatePlanData"
	err :=
		db.Model(&dbModel.PDF{}).
			Where("start_date between ? and ?", period.StartDate, period.EndDate).
			Or("end_date between ? and ?", period.StartDate, period.EndDate).
			Update("data_changed", true).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error(generalmodel.PlanCreationFailed, zap.Error(err), zap.String(generalmodel.AdditionalInfo, "PDF loading failed"))
		return nil, err
	}

	var meetings []dbModel.Meeting
	if err :=
		db.Preload("Tag").Where("date between ? and ?", period.StartDate, period.EndDate).
			Order("date").Find(&meetings).Error; err != nil {
		zap.L().Error(generalmodel.PlanCreationFailed, zap.Error(err), zap.String(generalmodel.AdditionalInfo, "Failed to load meetings"))
		return nil, err
	}

	var tasks []dbModel.TaskDetail
	if err := db.Find(&tasks).Error; err != nil {
		zap.L().Error(generalmodel.PlanCreationFailed, zap.Error(err), zap.String(generalmodel.AdditionalInfo, "Failed to load tasks"))
		return nil, err
	}

	var planIDs []uint
	for _, meeting := range meetings {
		if meeting.Tag.ID != 0 {
			var ids []uint
			if db.Table("plans").Where("meeting_id = ?", meeting.ID).Select("id").Find(&ids); len(ids) != 0 {
				continue
			}
			if err := db.Create(&dbmodel.Plan{MeetingID: meeting.ID}).Error; err != nil {
				zap.L().Error(generalmodel.PlanCreationFailed, zap.Error(err), zap.String(generalmodel.AdditionalInfo, "Failed to create plan element"))
				db.RollbackTo("beforePlanCreation")
			}
			continue
		}

		for _, task := range tasks {
			var ids []uint
			if db.Table("plans").Where("meeting_id = ?", meeting.ID).Where("task_detail_id = ?", task.ID).Select("id").Find(&ids); len(ids) != 0 {
				continue
			}

			person, err := getFirstPersonAvailable(meeting, task, period, db)
			if err != nil {
				zap.L().Info(generalmodel.PlanCreationError, zap.Error(err), zap.String(generalmodel.AdditionalInfo, "person loading error"))
			}
			if person == nil {
				person = &dbModel.Person{}
			}
			plan := dbModel.Plan{PersonID: person.ID, MeetingID: meeting.ID, TaskDetailID: task.ID}
			db.SavePoint("beforePlanCreation")
			if err := db.Create(&plan).Error; err != nil {
				db.RollbackTo("beforePlanCreation")
				zap.L().Error(generalmodel.PlanCreationFailed, zap.Error(err), zap.String(generalmodel.AdditionalInfo, "Failed to create plan element"))
			} else {
				planIDs = append(planIDs, plan.ID)
			}
		}
	}

	var plan []dbModel.Plan
	if err :=
		db.Preload("Person").
			Preload("Meeting").
			Preload("TaskDetail.Task").
			Preload("TaskDetail").
			Where("id IN (?)", planIDs).Find(&plan).Error; err != nil {
		return nil, err
	}

	return plan, nil
}

// UpdatePlanElement updates personId to the parameter, parameter also holds id for update
func UpdatePlanElement(element dbModel.Plan) error {
	db := database.DB.Begin()
	defer db.Commit()

	// is person allowed to be asigned to task
	var p dbModel.PersonTask
	if err :=
		db.Table("person_tasks").
			Where("person_id = ?", element.PersonID).
			Where("task_detail_id = ?", element.TaskDetailID).
			First(&p).Error; err != nil {
		return errors.ErrTaskForPersonNotAllowed
	}

	if err := db.Table("plans").
		Where("id = ?", element.ID).
		Update("person_id", element.PersonID).Error; err != nil {
		db.Rollback()
		return err
	}

	if err :=
		db.Table("pdfs").
			Where("(?) between start_date and end_date", db.Table("meetings").
				Where("id = (?)",
					db.Table("plans").Where("id = ?", element.ID).Select("meeting_id")).Select("date")).
			Update("data_changed", true).Error; err != nil {
		db.Rollback()
		return err
	}

	return nil
}

// GetAllPersonAvailable loads all available people for a meeting with the specified task
func GetAllPersonAvailable(db *gorm.DB, plan dbModel.Plan) (person apimodel.People, err error) {
	person.Available, err = getAvailablePeople(plan, generalmodel.Period{
		StartDate: time.Date(plan.Meeting.Date.Year(), plan.Meeting.Date.Month(), 1, 0, 0, 0, 0, time.Local),
		EndDate:   time.Date(plan.Meeting.Date.Year(), plan.Meeting.Date.Month()+1, 0, 0, 0, 0, 0, time.Local),
	}, db, false)

	if err != nil {
		return apimodel.People{}, err
	}
	ids := []uint{plan.PersonID}
	for _, person := range person.Available {
		ids = append(ids, person.ID)
	}

	err =
		db.Not("id in (?)", ids).
			Where("id in (select person_id from person_tasks where task_detail_id = ?)", plan.TaskDetailID).
			Find(&person.Absent).Error
	if err != nil {
		return apimodel.People{}, err
	}

	if plan.PersonID != 0 {
		if err =
			db.Where("id = ?", plan.PersonID).
				First(&person.Assigned).Error; err != nil {
			return apimodel.People{}, err
		}
	}

	return person, nil
}

// getFirstPersonAvailable loads the first available Person for a meeting with the specified task in the specified period
func getFirstPersonAvailable(meeting dbModel.Meeting, taskDetail dbModel.TaskDetail, period generalmodel.Period, db *gorm.DB) (person *dbModel.Person, err error) {
	people, err := getAvailablePeople(dbModel.Plan{TaskDetailID: taskDetail.ID, MeetingID: meeting.ID, Meeting: meeting}, period, db, true)
	if err != nil || len(people) == 0 {
		return nil, err
	}
	return &people[0], err
}

func getAvailablePeople(plan dbModel.Plan, period generalmodel.Period, db *gorm.DB, order bool) (person []dbModel.Person, err error) {
	timesInPeriod := `LEFT JOIN (
		SELECT person_id, COUNT(*) as all_entries
		FROM plans
			WHERE meeting_id in (
				SELECT id FROM meetings
				WHERE date between ? and ?
			) GROUP BY person_id
		) plan_count
		ON p.id = plan_count.person_id`
	tasksInPeriod := `LEFT JOIN (
		SELECT person_id, count(*) as task_count
		FROM plans
			WHERE meeting_id in (
				SELECT id FROM meetings
				WHERE date between ? and ?
			)
				AND task_detail_id = ?
			GROUP BY person_id
		) t_count
		ON p.id = t_count.person_id`
	peopleAssigned := db.Table("plans").
		Select("COALESCE(person_id, -1)").
		Where("meeting_id = ? AND person_id IS NOT NULL", plan.MeetingID)
	peopleAbsent := db.Table("person_absences").
		Select("COALESCE(person_id, -1)").
		Where("meeting_id = ?", plan.MeetingID)
	peopleRecuringAbsent := db.Table("person_recurring_absences").
		Select("COALESCE(person_id, -1)").
		Where("weekday = ?", plan.Meeting.Date.Weekday())

	query := db.Table("people p").
		// load task of person
		Joins("JOIN person_tasks pt ON p.id = pt.person_id").
		// load allowed tasks
		Joins("JOIN task_details td ON td.id = pt.task_detail_id").
		Joins(timesInPeriod, period.StartDate, period.EndDate).
		Joins(tasksInPeriod, period.StartDate, period.EndDate, plan.TaskDetailID).
		// filter task
		Where("td.id = ?", plan.TaskDetailID).
		Not("p.id IN (?)", peopleAssigned).
		Not("p.id IN (?)", peopleAbsent).
		Not("p.id IN (?)", peopleRecuringAbsent)
	if order {
		query = query.
			// Least entries in period first
			Order("plan_count.all_entries ASC NULLS FIRST").
			Order("t_count.task_count ASC NULLS FIRST")
	}

	if err := query.Find(&person).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return person, nil
}

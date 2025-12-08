package services

import (
	"lexxi/models"

	"github.com/kamva/mgm/v3"
)

func SaveDailyTasks(userID, planID string, plan models.StudyPlan) error {

	for _, day := range plan.DailyTable {
		for _, t := range day.Tasks {

			task := &models.TodoTask{
				UserID: userID,
				PlanID: planID,
				Date:   day.Date,
				Task:   t,
				Done:   false,
			}

			if err := mgm.Coll(task).Create(task); err != nil {
				return err
			}
		}
	}

	return nil
}

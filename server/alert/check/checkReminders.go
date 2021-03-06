package check

import (
	"fmt"

	"github.com/cjburchell/reefstatus/server/alert/slack"
	"github.com/cjburchell/reefstatus/server/alert/state"
	"github.com/cjburchell/reefstatus/server/data/repo"
)

// Reminders check
func Reminders(controller repo.Controller, state state.StateData, slackDestination string) error {
	info, err := controller.GetInfo()
	if err != nil {
		return err
	}

	for _, reminder := range info.Reminders {
		if reminder.IsOverdue {
			updated, err := state.UpdateReminderSent(reminder.Index, true)
			if err != nil {
				return err
			}

			if !updated {
				continue
			}

			var body = fmt.Sprintf("Reminder \"%s\" is now overdue", reminder.Text)
			err = slack.PrintMessage(body, slackDestination)
			if err != nil {
				return err
			}

		} else {
			_, err = state.UpdateReminderSent(reminder.Index, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

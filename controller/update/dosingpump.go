package update

import (
	"github.com/cjburchell/reefstatus/common/models"
	"github.com/cjburchell/reefstatus/common/profilux/types"
	"github.com/cjburchell/reefstatus/controller/service"
	"github.com/cjburchell/reefstatus/controller/profilux"
)

func dosingPumps(profiluxController *profilux.Controller, repo service.Controller) error {

	count, err := profiluxController.GetTimerCount()
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		pump, err := repo.GetDosingPump(models.GetID(models.DosingPumpType, i))

		if err != nil && err != service.ErrNotFound {
			return err
		}

		found := err != service.ErrNotFound
		settings, err := profiluxController.GetTimerSettings(i)
		if err != nil {
			return err
		}

		if settings.Mode == types.TimerModeAutoDosing {
			if !found {
				pump = models.NewDosingPump(i)
			}

			pump.Settings, err = profiluxController.GetTimerSettings(pump.Channel)
			if err != nil {
				return err
			}

			pump.Rate, err = profiluxController.GetDosingRate(pump.Channel)
			if err != nil {
				return err
			}

			pump.PerDay = pump.Settings.SwitchingCount

			err = repo.SetDosingPump(pump, !found)
			if err != nil {
				return err
			}

		} else {
			if found {
				err = repo.DeleteDosingPump(pump)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

package main

import (
	"strconv"
	"time"

	"github.com/cjburchell/reefstatus/controller/settings"

	logger "github.com/cjburchell/uatu-go"

	"github.com/cjburchell/reefstatus/controller/commands"

	"github.com/cjburchell/reefstatus/controller/update"

	"github.com/cjburchell/reefstatus/common/communication"
	"github.com/cjburchell/reefstatus/controller/service"
)

const logRate = time.Second * 30

func Update(session communication.PublishSession, isInitial bool) error {
	return session.Publish(communication.UpdateMessage, strconv.FormatBool(isInitial))
}

func main() {
	log := logger.Create()

	controller, err := service.NewController(settings.DataServiceAddress, settings.DataServiceToken)
	if err != nil {
		log.Fatal(err, "Unable to Connect to data database:")
	}

	session, err := communication.NewSession(settings.PubSubAddress, settings.PubSubToken, log)
	if err != nil {
		log.Fatal(err, "Unable to Connect to pub sub")
	}

	defer session.Close()
	go commands.Handle(session, controller, log)

	for {
		err = update.All(controller, log)
		if err == nil {
			err = Update(session, true)
			if err != nil {
				log.Errorf(err, "Unable to send first update")
			}
			break
		}

		log.Error(err, "Unable to do first update")
		log.Debugf("RefreshSettings Sleeping for %s", logRate.String())
		<-time.After(logRate)
		continue
	}

	updateCount := 0
	for {
		log.Debugf("RefreshSettings Sleeping for %s", logRate.String())
		<-time.After(logRate)
		if updateCount%20 == 19 {
			err = update.All(controller, log)
			if err != nil {
				log.Errorf(err, "Unable to update")
			}

		} else {
			err = update.State(controller, log)
			if err != nil {
				log.Errorf(err, "Unable to update state")
			}
		}

		err = Update(session, false)
		if err != nil {
			log.Errorf(err, "Unable to send update")
		}
		updateCount++
	}
}

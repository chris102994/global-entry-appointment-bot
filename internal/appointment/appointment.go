package appointment

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/chris102994/global-entry-appointment-bot/internal/run"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Appointment struct {
	LocationID     int        `json:"locationId"`
	StartTimestamp CustomTime `json:"startTimestamp"`
	EndTimestamp   CustomTime `json:"endTimestamp"`
	Active         bool       `json:"active"`
	Duration       int        `json:"duration"`
	RemoteInd      bool       `json:"remoteInd"`
}

type CustomTime struct {
	time.Time
}

func (c *CustomTime) DateFmt() string {
	return c.Format("2006-01-02")
}

func (c *CustomTime) TimeFmt() string {
	return c.Format("03:04 PM")
}

func (c *CustomTime) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02T15:04"

	parsedTime, err := time.Parse(`"`+layout+`"`, string(b))
	if err != nil {
		return err
	}

	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "UTC"
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("could not load location %v: %w", tz, err)
	}

	c.Time = parsedTime.In(location)

	return nil
}

func GetAppointments(runConfig run.Run) ([]*Appointment, error) {
	log.WithFields(log.Fields{
		"runConfig": runConfig,
	}).Trace("Running with Config")

	var availableAppointmentSlots []*Appointment

	baseUrl := "https://ttp.cbp.dhs.gov/schedulerapi/slots"

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: customTransport,
		Timeout:   10 * time.Second,
	}

	if runConfig.LocationIDs == nil {
		log.WithFields(log.Fields{}).Fatal("Location IDs needed to run")
		return availableAppointmentSlots, fmt.Errorf("no Location IDs specified")
	}

	for _, locationId := range runConfig.LocationIDs {
		log.WithFields(log.Fields{
			"LocationID": locationId,
		}).Trace("Running with")

		var thisAvailableAppointmentSlots []*Appointment
		requestParams := url.Values{
			"orderBy":    {runConfig.OrderBy},
			"limit":      {strconv.Itoa(runConfig.Limit)},
			"locationId": {locationId},
			"minimum":    {strconv.Itoa(runConfig.Minimum)},
		}

		requestUrl := baseUrl + "?" + requestParams.Encode()

		log.WithFields(log.Fields{
			"requestUrl": requestUrl,
		}).Trace("Scraping endpoint")

		resp, err := client.Get(requestUrl)
		if err != nil {
			return availableAppointmentSlots, err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return availableAppointmentSlots, err
		}

		err = json.Unmarshal(bodyBytes, &thisAvailableAppointmentSlots)
		if err != nil {
			return availableAppointmentSlots, err
		}

		availableAppointmentSlots = append(availableAppointmentSlots, thisAvailableAppointmentSlots...)
	}

	return availableAppointmentSlots, nil
}

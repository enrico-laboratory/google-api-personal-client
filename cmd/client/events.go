package client

import (
	"errors"
	"github.com/enrico-laboratory/go-validator"
	"google.golang.org/api/calendar/v3"
	"log"
	"time"
)

type GEvent struct {
	Description   string
	EndDateTime   GEventDateTime
	Location      string
	StartDateTime GEventDateTime
	Summary       string
}

type GEventDateTime struct {
	Date     time.Time
	DateTime time.Time
}

func (c *GClient) Insert(calendarID string, event *GEvent) (string, error) {
	v := validator.New()

	if validateDates(v, event); !v.Valid() {
		return "", errors.New(v.ErrorsToString())
	}

	var endDateTime calendar.EventDateTime
	if event.EndDateTime.Date.IsZero() && !event.EndDateTime.DateTime.IsZero() {
		endDateTime.DateTime = event.EndDateTime.DateTime.Format(c.config.formatTime)
	} else if !event.EndDateTime.Date.IsZero() && event.EndDateTime.DateTime.IsZero() {
		endDateTime.Date = event.EndDateTime.Date.Format("2006-01-02")
	}
	endDateTime.TimeZone = c.config.timeZone

	var startDateTime calendar.EventDateTime
	if event.StartDateTime.Date.IsZero() && !event.StartDateTime.DateTime.IsZero() {
		startDateTime.DateTime = event.StartDateTime.DateTime.Format(c.config.formatTime)
	} else if !event.StartDateTime.Date.IsZero() && event.StartDateTime.DateTime.IsZero() {
		startDateTime.Date = event.StartDateTime.Date.Format("2006-01-02")
	}
	startDateTime.TimeZone = c.config.timeZone

	var overrides []*calendar.EventReminder
	var override calendar.EventReminder
	override.Method = "popup"
	override.Minutes = 90
	overrides = append(overrides, &override)
	//reminders := &calendar.EventReminders{
	//	Overrides:  overrides,
	//	UseDefault: false,
	//}

	eventParsed := &calendar.Event{
		Description: event.Description,
		End:         &endDateTime,
		Location:    event.Location,
		Start:       &startDateTime,
		//Reminders:   reminders,
		Status:  "confirmed",
		Summary: event.Summary,
	}
	log.Println(eventParsed.End.Date)
	log.Println(eventParsed.Start.Date)

	resp, err := c.Service.Events.Insert(calendarID, eventParsed).Do()
	if err != nil {
		return "", err
	}

	return resp.Id, nil

}

func (c *GClient) List(calendarID string) ([]string, error) {

	resp, err := c.ListByTimeMax(calendarID, time.Date(1900, 01, 01, 00, 00, 00, 00, &time.Location{}))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *GClient) ListByTimeMax(calendarID string, timeMax time.Time) ([]string, error) {

	resp, err := c.Service.Events.List(calendarID).TimeMin(timeMax.Format(c.config.formatTime)).Do()
	if err != nil {
		return nil, err
	}

	var idList []string
	for _, event := range resp.Items {
		idList = append(idList, event.Id)
	}
	return idList, nil
}

func (c *GClient) Delete(calendarID string, eventId string) error {
	err := c.Service.Events.Delete(calendarID, eventId).Do()
	if err != nil {
		return err
	}
	return nil
}

func validateDates(v *validator.Validator, event *GEvent) {

	v.Check(bothDateNotFull(&event.StartDateTime), "Start-date", "chose Date or DateTime, cannot input both")
	v.Check(bothDateNotEmpty(&event.StartDateTime), "Start-date", "at least one date must be present")
	v.Check(bothDateNotFull(&event.EndDateTime), "end-date", "chose Date or DateTime, cannot input both")
	v.Check(bothDateNotEmpty(&event.EndDateTime), "end-date", "at least one date must be present")
}

func bothDateNotEmpty(date *GEventDateTime) bool {
	if date.Date.IsZero() && date.DateTime.IsZero() {
		return false
	}
	return true
}

func bothDateNotFull(date *GEventDateTime) bool {
	if !date.Date.IsZero() && !date.DateTime.IsZero() {
		return false
	}
	return true
}

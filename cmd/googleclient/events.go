package googleclient

import (
	"errors"
	"github.com/enrico-laboratory/go-validator"
	"google.golang.org/api/calendar/v3"
	"time"
)

type GEventService struct {
	service *calendar.Service
	config  *config
}

type GEventModel struct {
	EventID       string
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

func (c *GEventService) Insert(calendarID string, event *GEventModel) (string, error) {
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

	resp, err := c.service.Events.Insert(calendarID, eventParsed).Do()
	if err != nil {
		return "", err
	}

	return resp.Id, nil

}

func (c *GEventService) List(calendarID string) ([]GEventModel, error) {

	resp, err := c.ListByTimeMin(calendarID, time.Date(1900, 01, 01, 00, 00, 00, 00, &time.Location{}))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *GEventService) ListByTimeMin(calendarID string, timeMax time.Time) ([]GEventModel, error) {

	resp, err := c.service.Events.List(calendarID).TimeMin(timeMax.Format(c.config.formatTime)).Do()
	if err != nil {
		return nil, err
	}

	var idList []string
	for _, event := range resp.Items {
		idList = append(idList, event.Id)
	}

	var gEvents []GEventModel

	for _, event := range resp.Items {

		var startDateTimeObject GEventDateTime
		var endDateTimeObject GEventDateTime

		if event.Start.Date != "" {
			layoutDate := "2006-01-02"
			startDate, err := time.Parse(layoutDate, event.Start.Date)
			if err != nil {
				return nil, err
			}
			endDate, err := time.Parse(layoutDate, event.Start.Date)
			if err != nil {
				return nil, err
			}
			startDateTimeObject.Date = startDate
			endDateTimeObject.Date = endDate
		} else {
			startDateTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
			if err != nil {
				return nil, err
			}
			endDateTime, err := time.Parse(time.RFC3339, event.End.DateTime)
			if err != nil {
				return nil, err
			}
			startDateTimeObject.DateTime = startDateTime
			endDateTimeObject.DateTime = endDateTime
		}

		gEvent := GEventModel{
			EventID:       event.Id,
			Description:   event.Description,
			EndDateTime:   endDateTimeObject,
			Location:      event.Location,
			StartDateTime: startDateTimeObject,
			Summary:       event.Summary,
		}
		gEvents = append(gEvents, gEvent)
	}

	return gEvents, nil
}

func (c *GEventService) Delete(calendarID string, eventId string) error {
	err := c.service.Events.Delete(calendarID, eventId).Do()
	if err != nil {
		return err
	}
	return nil
}

func validateDates(v *validator.Validator, event *GEventModel) {

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

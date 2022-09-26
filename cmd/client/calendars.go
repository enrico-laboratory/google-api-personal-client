package client

import "google.golang.org/api/calendar/v3"

type GCalendar struct {
	Description string
	Location    string
	Summary     string
	TimeZone    string
}

type GCalendarService struct {
	Service *calendar.Service
}

func (c *GCalendarService) Insert(cal *GCalendar) (string, error) {

	call := c.Service.Calendars.Insert(&calendar.Calendar{
		Description: cal.Description,
		Location:    cal.Location,
		Summary:     cal.Summary,
		TimeZone:    cal.TimeZone,
	})

	resp, err := call.Do()
	if err != nil {
		return "", err
	}

	return resp.Id, err
}

func (c *GCalendarService) Get(id string) (string, error) {
	response, err := c.Service.Calendars.Get(id).Do()
	if err != nil {
		return "", err
	}
	return response.Summary, nil
}

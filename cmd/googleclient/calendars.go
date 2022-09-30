package googleclient

import "google.golang.org/api/calendar/v3"

type GCalendarModel struct {
	Id          string
	Description string
	Location    string
	Summary     string
}

type GCalendarService struct {
	service *calendar.Service
	config  *config
}

func (c *GCalendarService) Insert(summary string) (string, error) {

	resp, err := c.service.Calendars.Insert(&calendar.Calendar{
		Summary: summary,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *GCalendarService) Get(calendarId string) (string, error) {
	resp, err := c.service.Calendars.Get(calendarId).Do()
	if err != nil {
		return "", err
	}
	return resp.Summary, nil
}

// Patch adds only default reminders (method: popup, minutes: 90) the argument "cal" is not considered
func (c *GCalendarService) Patch(calendarID string, cal *GCalendarModel) (string, error) {
	defaultReminder := &calendar.EventReminder{
		Method:  "popup",
		Minutes: 90,
	}
	var reminderList []*calendar.EventReminder
	reminderList = append(reminderList, defaultReminder)

	calendarList := &calendar.CalendarListEntry{
		//ForegroundColor:  cal.ColorId,
		DefaultReminders: reminderList,
		Description:      cal.Description,
		Id:               calendarID,
		Location:         cal.Location,
		//SummaryOverride: cal.Summary,
		//TimeZone: c.config.timeZone,
	}
	resp, err := c.service.CalendarList.Patch(calendarID, calendarList).Do()

	if err != nil {
		return "", err
	}

	return resp.Summary, err
}

func (c *GCalendarService) Delete(calendarID string) error {
	err := c.service.Calendars.Delete(calendarID).Do()
	if err != nil {
		return err
	}
	return nil
}

func (c *GCalendarService) List() ([]GCalendarModel, error) {
	list, err := c.service.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	var gCalendarList []GCalendarModel
	for _, calendar := range list.Items {
		var gCalendar GCalendarModel
		gCalendar.Description = calendar.Description
		gCalendar.Location = calendar.Location
		gCalendar.Summary = calendar.Summary
		gCalendar.Id = calendar.Id

		gCalendarList = append(gCalendarList, gCalendar)
	}

	return gCalendarList, nil

}

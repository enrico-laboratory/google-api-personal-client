package googleclient

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCalendar(t *testing.T) {

	credFileName := "credentials.json"

	ctx := context.Background()

	c, err := NewClient(credFileName, ctx)
	if err != nil {
		log.Fatal(err)
	}

	var calendarID string

	t.Run("INSERT Calendar", func(t *testing.T) {
		summary := "Test Calendar"
		result, err := c.GCalendar.Insert(summary)
		t.Log(result)
		assert.Empty(t, err)
		assert.NotEmpty(t, result)

		calendarID = result
	})

	t.Run("GET Calendar", func(t *testing.T) {
		result, err := c.GCalendar.Get(calendarID)
		expected := "Test Calendar"
		actual := result
		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("UPDATE Calendar", func(t *testing.T) {
		calendar := &GCalendarModel{
			Description: "Test description",
			Location:    "Unknown location",
			Summary:     "Test Calendar Override",
		}
		result, err := c.GCalendar.Patch(calendarID, calendar)
		expected := "Test Calendar"
		actual := result
		assert.Empty(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("DELETE Calendar", func(t *testing.T) {
		err := c.GCalendar.Delete(calendarID)
		assert.Empty(t, err)
	})

	t.Run("LIST Calendars", func(t *testing.T) {
		list, err := c.GCalendar.List()
		assert.Empty(t, err)
		assert.True(t, len(list) > 1)
		for _, calendar := range list {
			log.Println("Summary:", calendar.Summary)
		}
	})
}

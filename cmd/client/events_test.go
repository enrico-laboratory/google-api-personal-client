package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {

	credFileName := "credentials.json"
	calendarID := "37bf2cf949f9e9d908a0daf93fa16b8af477aa4a8e94c6e9de5e03711f0cf7ef@group.calendar.google.com"

	ctx := context.Background()

	c, err := NewClient(credFileName, calendarID, ctx)
	if err != nil {
		log.Fatal(err)
	}
	testEventDate := &GEvent{
		Description: "Short description",
		EndDateTime: GEventDateTime{
			Date: time.Now().Add(24 * time.Hour),
		},
		Location: "Test Location",
		StartDateTime: GEventDateTime{
			Date: time.Now().Add(24 * time.Hour), // Starts Tomorrow
		},
		Summary: "TEST Event Date",
	}

	testEventDateTime := &GEvent{
		Description: "Short description",
		EndDateTime: GEventDateTime{
			DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
		},
		Location: "Test Location",
		StartDateTime: GEventDateTime{
			DateTime: time.Now().Add(48 * time.Hour), // Starts in 2 days
		},
		Summary: "TEST Event DateTime",
	}

	var eventIDList []string

	t.Run("INSERT Event with Date", func(t *testing.T) {
		result, err := c.Insert(calendarID, testEventDate)

		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("INSERT Event with DateTime", func(t *testing.T) {
		result, err := c.Insert(calendarID, testEventDateTime)
		t.Log(result)
		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("LIST all events TimeMax, finds event in 2 days", func(t *testing.T) {
		result, err := c.ListByTimeMax(calendarID, time.Now().Add(46*time.Hour))
		t.Log(result[0])
		assert.Empty(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("LIST all events", func(t *testing.T) {
		result, err := c.List(calendarID)
		assert.Empty(t, err)
		assert.Equal(t, 2, len(result))

		eventIDList = append(eventIDList, result...)
	})

	t.Run("DELETE all events", func(t *testing.T) {
		for _, eventID := range eventIDList {
			err := c.Delete(calendarID, eventID)
			assert.Empty(t, err)
		}
		result, err := c.List(calendarID)
		assert.Empty(t, err)
		assert.True(t, len(result) == 0)
	})

	eventDateValidation := []GEventDateTime{
		{
			Date:     time.Now(),
			DateTime: time.Now(),
		},
		{
			Date:     time.Time{},
			DateTime: time.Time{},
		},
	}

	for _, date := range eventDateValidation {
		t.Run("INSERT test Start date validation", func(t *testing.T) {
			testEventValidation := &GEvent{
				Description: "Short description",
				EndDateTime: GEventDateTime{
					DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
				},
				Location:      "Test Location",
				StartDateTime: date,
				Summary:       "TEST Event DateTime",
			}
			result, err := c.Insert(calendarID, testEventValidation)
			t.Log(err)
			assert.NotEmpty(t, err)
			assert.Empty(t, result)
		})

		for _, date := range eventDateValidation {
			t.Run("INSERT test End Date validation", func(t *testing.T) {
				testEventValidation := &GEvent{
					Description: "Short description",
					EndDateTime: date,
					Location:    "Test Location",
					StartDateTime: GEventDateTime{
						DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
					},
					Summary: "TEST Event DateTime",
				}
				result, err := c.Insert(calendarID, testEventValidation)
				t.Log(err)
				assert.NotEmpty(t, err)
				assert.Empty(t, result)
			})

		}
	}
}

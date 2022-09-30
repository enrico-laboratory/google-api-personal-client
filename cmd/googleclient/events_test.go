package googleclient

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var calendarID string
var c *GClient

func TestMain(m *testing.M) {

	credFileName := "credentials.json"

	ctx := context.Background()

	var err error
	c, err = NewClient(credFileName, ctx)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := c.GCalendar.Insert("Test Event Calendar")
	if err != nil {
		log.Println(err)
	}
	calendarID = resp

	exitVal := m.Run()

	err = c.GCalendar.Delete(calendarID)
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitVal)
}

func TestEvents(t *testing.T) {

	testEventDate := &GEventModel{
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

	testEventDateTime := &GEventModel{
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

	var eventIDList []GEventModel

	t.Run("INSERT Event with Date", func(t *testing.T) {
		result, err := c.GEvent.Insert(calendarID, testEventDate)

		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("INSERT Event with DateTime", func(t *testing.T) {
		result, err := c.GEvent.Insert(calendarID, testEventDateTime)
		t.Log(result)
		assert.Empty(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("LIST all events TimeMax, finds event in 2 days", func(t *testing.T) {
		result, err := c.GEvent.ListByTimeMin(calendarID, time.Now().Add(46*time.Hour))
		t.Log(result[0])
		assert.Empty(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("LIST all events", func(t *testing.T) {
		result, err := c.GEvent.List(calendarID)
		assert.Empty(t, err)
		assert.Equal(t, 2, len(result))

		eventIDList = append(eventIDList, result...)
	})

	t.Run("DELETE all events", func(t *testing.T) {
		for _, event := range eventIDList {
			err := c.GEvent.Delete(calendarID, event.EventID)
			assert.Empty(t, err)
		}
		result, err := c.GEvent.List(calendarID)
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
			testEventValidation := &GEventModel{
				Description: "Short description",
				EndDateTime: GEventDateTime{
					DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
				},
				Location:      "Test Location",
				StartDateTime: date,
				Summary:       "TEST Event DateTime",
			}
			result, err := c.GEvent.Insert(calendarID, testEventValidation)
			t.Log(err)
			assert.NotEmpty(t, err)
			assert.Empty(t, result)
		})

		for _, date := range eventDateValidation {
			t.Run("INSERT test End Date validation", func(t *testing.T) {
				testEventValidation := &GEventModel{
					Description: "Short description",
					EndDateTime: date,
					Location:    "Test Location",
					StartDateTime: GEventDateTime{
						DateTime: time.Now().Add(48 * time.Hour).Add(2 * time.Hour),
					},
					Summary: "TEST Event DateTime",
				}
				result, err := c.GEvent.Insert(calendarID, testEventValidation)
				t.Log(err)
				assert.NotEmpty(t, err)
				assert.Empty(t, result)
			})

		}
	}
}

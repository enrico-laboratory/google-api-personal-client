package googleclient

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"time"
)

type config struct {
	formatTime string
	timeZone   string
}

type GClient struct {
	GEvent    GEventService
	GCalendar GCalendarService
}

func NewClient(ctx context.Context) (*GClient, error) {

	client, err := getClient()
	if err != nil {
		return nil, err
	}
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	cfg := &config{
		formatTime: time.RFC3339,
		timeZone:   "Europe/Amsterdam",
	}

	gClient := &GClient{
		GEvent: GEventService{
			service: srv,
			config:  cfg,
		},
		GCalendar: GCalendarService{
			service: srv,
			config:  cfg,
		},
	}

	return gClient, nil
}

// getClient with service account generated key
func getClient() (*http.Client, error) {

	f, err := ioutil.ReadFile("key.json")
	if err != nil {
		return nil, err
	}
	cred, err := google.CredentialsFromJSON(context.Background(), f, "https://www.googleapis.com/auth/calendar")
	c := oauth2.NewClient(context.Background(), cred.TokenSource)
	return c, nil
}

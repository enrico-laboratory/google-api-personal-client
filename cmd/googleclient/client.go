package googleclient

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type config struct {
	formatTime string
	timeZone   string
}

var jsonCredentials string

type GClient struct {
	GEvent    GEventService
	GCalendar GCalendarService
}

// NewClient create a client from a context and the credentials file path. The credentials file path
// can be omitted if the GOOGLE_CREDENTIALS env variable is set instead
func NewClient(ctx context.Context, keyPath string) (*GClient, error) {
	jsonCredentials = os.Getenv("GOOGLE_CREDENTIALS")
	if keyPath == "" {
		if jsonCredentials == "" {
			return nil, errors.New("one of key path or GOOGLE_CREDENTIALS env must be define")
		}
	}

	client, err := getClient(keyPath)
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
func getClient(keyPath string) (*http.Client, error) {
	var err error
	var credByte []byte

	if jsonCredentials != "" {
		credByte = []byte(jsonCredentials)
	} else {
		credByte, err = ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
	}

	cred, err := google.CredentialsFromJSON(context.Background(), credByte, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		return nil, err
	}
	c := oauth2.NewClient(context.Background(), cred.TokenSource)
	return c, nil
}

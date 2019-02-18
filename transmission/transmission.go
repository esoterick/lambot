package transmission

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// SessionHeader header for the http requests session id
const SessionHeader string = "X-Transmission-Session-Id"

// Tor is a base struct for our Transmission client
type Tor struct {
	URL                 string
	SessionID           string
	AuthorizationHeader string
	Client              *http.Client
	RetryCount          int
	RetryBackoff        time.Duration
}

// NewTor returns a new Tor struct for us to work with the API
func NewTor(url string, username string, password string) (*Tor, error) {
	req, _ := http.NewRequest("POST", url, nil)
	req.SetBasicAuth(username, password)
	authHeader := req.Header.Get("Authorization")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if len(res.Header[SessionHeader]) < 1 {
		return nil, errors.New("unable to get session id from response header")
	}

	t := Tor{
		URL:                 url,
		SessionID:           res.Header[SessionHeader][0],
		AuthorizationHeader: authHeader,
		Client:              client,
		RetryCount:          3,
		RetryBackoff:        1 * time.Second,
	}

	return &t, nil
}

// Post is a helper method to send a POST request to the Transmission RPC API
func (t *Tor) Post(payload io.Reader) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var err error
	var body []byte

	req, err := http.NewRequest("POST", t.URL, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", t.AuthorizationHeader)
	req.Header.Add(SessionHeader, t.SessionID)

	for i := 0; i < t.RetryCount; i++ {
		resp, e := t.Client.Do(req.WithContext(ctx))
		if e != nil {
			if err == nil {
				err = e
			} else {
				err = errors.New(err.Error() + "\n" + e.Error())
			}
			time.Sleep(t.RetryBackoff)
			continue
		}

		if resp != nil {
			if resp.StatusCode == 0 || resp.StatusCode >= 400 {
				e := fmt.Errorf("caught retryable non-good status code: %d", resp.StatusCode)
				if err == nil {
					err = e
				} else {
					err = errors.New(err.Error() + "\n" + e.Error())
				}

				time.Sleep(t.RetryBackoff)
				continue
			}

			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)

			if err != nil {
				return nil, err
			}

			if body != nil {
				return body, nil
			}
		}
	}

	return body, nil
}

// GetTorrents gets the torrents...
func (t *Tor) GetTorrents() ([]Torrent, error) {
	req := Request{
		Method: "torrent-get",
		Arguments: Arguments{
			Fields: []string{
				"id",
				"name",
				"status",
				"rateDownload",
				"rateUpload",
			},
		},
	}

	b, err := req.MarshalJSON()
	if err != nil {
		return nil, err
	}

	payload := bytes.NewReader(b)

	p, err := t.Post(payload)
	if err != nil {
		return nil, err
	}

	var resp TorrentsResponse
	err = resp.UnmarshalJSON(p)
	if err != nil {
		return nil, err
	}

	return resp.Arguments.Torrents, nil
}

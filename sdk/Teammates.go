package sendgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Teammates is a Sendgrid Teammates.
type Teammates struct {
	Email   string `json:"email,omitempty"`
	IsAdmin bool   `json:"is_admin,omitempty"`
	// TODO:
	// Scopes  []string `json:"scopes,omitempty"`
	// PendingID
}

func parseTeammate(respBody string) (*Teammates, RequestError) {
	var body Teammates
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		log.Printf("[DEBUG] [parseTeammate] failed parsing teammate, response body: %s", respBody)

		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseTeammates(respBody string) ([]Teammates, RequestError) {
	var body []Teammates
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		log.Printf("[DEBUG] [parseTeammates] failed parsing teammates, response body: %s", respBody)

		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateTeammates creates a teammates and returns it.
func (c *Client) CreateTeammates(email string, isAdmin bool, scopes []string) (*Teammates, RequestError) {

	if email == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrEmailRequired}
	}

	// TODO: Switch if teammates exists and pending
	respBody, statusCode, err := c.Post("POST", "/teammates", Teammates{
		Email:   email,
		IsAdmin: isAdmin,
		// Scopes:  scopes,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed creating teammates: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingTeammates, statusCode, respBody),
		}
	}

	return parseTeammate(respBody)
}

// ReadTeammates retreives a teammates and returns it.
func (c *Client) ReadTeammates(username string) ([]Teammates, RequestError) {
	if username == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	fmt.Println("failed reading teammates: %w", username)

	endpoint := "/teammates/" + username

	respBody, statusCode, err := c.Get("GET", endpoint)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed reading teammates: %w", err),
		}
	}

	return parseTeammates(respBody)
}

// UpdateTeammates enables/disables a teammates.
// TODO:
func (c *Client) UpdateTeammates(username string, disabled bool) (bool, RequestError) {

	return false, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// DeleteTeammates deletes a teammates
func (c *Client) DeleteTeammates(username string) (bool, RequestError) {
	if username == "" {
		return false, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	respBody, statusCode, err := c.Get("DELETE", "/teammates/"+username)
	if err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed deleting teammates: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound { // ignore not found
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w: statusCode: %d, respBody: %s", err, statusCode, respBody),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}

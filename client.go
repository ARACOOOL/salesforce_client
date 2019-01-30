package salesforce_client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	version       = "0.0.1"
	EnvProduction = "prod"
	EnvStaging    = "stg"
	productionUrl = "https://login.salesforce.com/services/oauth2/token"
	stagingUrl    = "https://test.salesforce.com/services/oauth2/token"
)

type Auth struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
}

type Request struct {
	Method   string
	Endpoint string
	Body     io.Reader
}

type Response struct {
	Body          []byte
	StatusCode    int
	ContentLength int64
}

type Client struct {
	Http      *http.Client
	Env       string
	Version   string
	UserAgent string
	baseUrl   string
	token     string
}

// Find returns one specific object
func (c *Client) Find(object, ID string, r interface{}) error {
	response, err := c.send(Request{
		"GET",
		"/services/data/v" + c.Version + "/sobjects/" + object + "/" + ID,
		nil,
	})
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if r != nil {
		if w, ok := r.(io.Writer); ok {
			_, _ = io.Copy(w, response.Body)
		} else {
			err = json.NewDecoder(response.Body).Decode(r)
		}
	}

	return nil
}

// Query returns the query result
func (c *Client) Query(query string, r interface{}) error {
	response, err := c.send(Request{
		"GET",
		"/services/data/v" + c.Version + "/query?q=" + query,
		nil,
	})
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if r != nil {
		if w, ok := r.(io.Writer); ok {
			_, _ = io.Copy(w, response.Body)
		} else {
			err = json.NewDecoder(response.Body).Decode(r)
		}
	}

	return nil
}

// Create creates an object and returns the result
func (c *Client) Create(object string, params *Params, r interface{}) error {
	body, _ := json.Marshal(params.GetFields())
	response, err := c.send(Request{
		"POST",
		"/services/data/v" + c.Version + "/sobjects/" + object,
		bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if r != nil {
		if w, ok := r.(io.Writer); ok {
			_, _ = io.Copy(w, response.Body)
		} else {
			err = json.NewDecoder(response.Body).Decode(r)
		}
	}

	return nil
}

// Update updates specific object
func (c *Client) Update(object string, ID string, params *Params) error {
	body, _ := json.Marshal(params.GetFields())
	_, err := c.send(Request{
		"PATCH",
		"/services/data/v" + c.Version + "/sobjects/" + object + "/" + ID,
		bytes.NewReader(body),
	})
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes specific object
func (c *Client) Delete(object string, ID string) error {
	_, err := c.send(Request{
		"DELETE",
		"/services/data/v" + c.Version + "/sobjects/" + object + "/" + ID,
		nil,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) send(r Request) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, c.baseUrl+r.Endpoint, r.Body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.Http.Do(req)
	if err != nil {
		panic(err)
	}

	return resp, nil
}

// Auth gets the auth token and the base url
// This method should be called first after the creation of the client
func (c *Client) Auth(auth Auth) error {
	var loginUrl string

	if c.Env == EnvProduction {
		loginUrl = productionUrl
	} else {
		loginUrl = stagingUrl
	}

	form := url.Values{}
	form.Add("username", auth.Username)
	form.Add("password", auth.Password)
	form.Add("client_id", auth.ClientID)
	form.Add("client_secret", auth.ClientSecret)
	form.Add("grant_type", "password")

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	authResponse := &AuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(authResponse)
	if err != nil {
		return err
	}

	c.baseUrl = authResponse.InstanceUrl
	c.token = authResponse.AccessToken

	return nil
}

// NewClient creates a new SalesForce API client
func NewClient(env, apiVersion string) (*Client, error) {
	return &Client{
		Env:       env,
		Version:   apiVersion,
		UserAgent: "SalesForce_Go_Client/" + version,
		Http: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type API struct {
	Client *http.Client
	Token  string
}

func (c API) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	var host string
	if customEndpoint := os.Getenv("C3PM_API_ENDPOINT"); customEndpoint == "" {
		host = "https://c3pm.herokuapp.com/v1"
	} else {
		host = customEndpoint
	}
	url := host + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}
	return req, err
}

func (c API) fetch(method string, path string, body io.Reader, data interface{}) error {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		return handleHttpError(resp)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}

func (c API) send(method string, path string, buf io.Reader) error {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("file", "package.tar")
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	_, err = part.Write(contents)
	if err != nil {
		return err
	}
	w.Close()
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		return handleHttpError(resp)
	}
	return nil
}

func handleHttpError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var message string
	var parsedBody struct {
		Error string `json:"error"`
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		message = string(body)
	} else {
		message = parsedBody.Error
	}

	return fmt.Errorf("Client error: '%s' failed (%d): '%s'", resp.Request.URL, resp.StatusCode, message)
}

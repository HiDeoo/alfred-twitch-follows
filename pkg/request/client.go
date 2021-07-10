package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient httpClient
	headers    http.Header
}

type Response struct {
	Data       []byte
	StatusCode int
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(baseURL string) Client {
	return Client{baseURL, &http.Client{Timeout: time.Second * 5}, make(http.Header)}
}

func (client *Client) SetHeaders(headers http.Header) {
	client.headers = headers
}

func (client *Client) SetClient(httpClient httpClient) {
	client.httpClient = httpClient
}

func (client *Client) Get(path string, queryParams url.Values) (*Response, error) {
	return client.request(http.MethodGet, path, queryParams, nil)
}

func (client *Client) Post(path string, queryParams url.Values, bodyJSON map[string]string) (*Response, error) {
	return client.request(http.MethodPost, path, queryParams, bodyJSON)
}

func (client *Client) request(
	method,
	path string,
	queryParams url.Values,
	bodyJSON map[string]string,
) (*Response, error) {
	var bodyBuffer io.Reader = nil

	if bodyJSON != nil {
		bodyData, err := json.Marshal(bodyJSON)

		if err != nil {
			return nil, err
		}

		bodyBuffer = bytes.NewBuffer(bodyData)
	}

	req, err := http.NewRequest(method, client.baseURL+path+"?"+queryParams.Encode(), bodyBuffer)

	if err != nil {
		return nil, err
	}

	req.Header = client.headers

	if bodyJSON != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return &Response{body, res.StatusCode}, nil
}

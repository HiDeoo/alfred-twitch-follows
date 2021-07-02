package request

import (
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
	req, err := http.NewRequest(http.MethodGet, client.baseURL+path+"?"+queryParams.Encode(), nil)

	if err != nil {
		return nil, err
	}

	req.Header = client.headers

	res, err := client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	return &Response{body, res.StatusCode}, nil
}

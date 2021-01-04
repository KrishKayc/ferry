package httprequest

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

//HTTPRequest represents the apps request
type HTTPRequest struct {
	URL       string
	Path      string
	AuthToken string
	Params    map[string]string
}

//Send sends the request
func (httpreq *HTTPRequest) Send() []byte {
	client := &http.Client{}
	resp, err := client.Do(httpreq.get())
	HandleError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	HandleError(err)

	return body
}

//NewHTTPRequest ..
func NewHTTPRequest(url string, path string, authToken string, params map[string]string) *HTTPRequest {
	return &HTTPRequest{URL: url, Path: path, AuthToken: authToken, Params: params}
}

func (httpreq *HTTPRequest) get() *http.Request {
	var finalPath string
	bearer := "Basic " + httpreq.AuthToken
	if httpreq.Params != nil {
		var endPoint *url.URL
		endPoint, err := url.Parse(httpreq.URL)
		HandleError(err)

		endPoint.Path += httpreq.Path
		parameters := url.Values{}

		for k, v := range httpreq.Params {
			parameters.Add(k, v)
		}

		endPoint.RawQuery = parameters.Encode()
		finalPath = endPoint.String()

	} else {
		finalPath = httpreq.URL + httpreq.Path
	}

	req, err := http.NewRequest("GET", finalPath, nil)
	req.Header.Add("Authorization", bearer)
	HandleError(err)

	return req
}

//HandleError handles errors
func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

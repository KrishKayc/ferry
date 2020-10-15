package httprequest

// JiraClient represents a basic API client for Jira Rest API
type JiraClient struct {
	URL       string
	AuthToken string
}

// NewClient create a new instance of API client
func NewClient(URL, authToken string) *JiraClient {
	return &JiraClient{
		URL,
		authToken,
	}
}

// Get process the Jira Rest API authenticated request
func (c *JiraClient) Get(path string, params map[string]string) []byte {
	req := NewHTTPRequest(c.URL, path, c.AuthToken, params)

	return req.Send()
}

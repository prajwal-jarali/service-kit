package api

type HttpResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

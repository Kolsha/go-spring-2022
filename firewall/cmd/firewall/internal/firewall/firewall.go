package firewall

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var _ http.RoundTripper = (*firewall)(nil)

type firewall struct {
	config *Rules
}

func NewFirewall(config *Rules) http.RoundTripper {
	f := &firewall{
		config: config,
	}
	return f
}

func (f *firewall) RoundTrip(request *http.Request) (*http.Response, error) {
	rule := f.getRule(request.RequestURI)

	if !rule.isRequestAllowed(request) {
		return getBlockedQuery(), nil
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if !rule.isResponseAllowed(response) {
		return getBlockedQuery(), nil
	}

	return response, nil
}

func (f *firewall) getRule(requestURI string) *Rule {
	for _, rule := range f.config.Rules {
		if rule.Endpoint == requestURI {
			return &rule
		}
	}
	return nil
}

func (r *Rule) isRequestAllowed(request *http.Request) bool {
	if r == nil {
		return true
	}

	// check UserAgent
	UserAgent := request.UserAgent()
	for _, v := range r.ForbiddenUserAgents {
		matchString, err := regexp.MatchString(v, UserAgent)
		if matchString || err != nil {
			return false
		}
	}

	// check RequiredHeaders
	for _, v := range r.RequiredHeaders {
		if request.Header.Get(v) == "" {
			return false
		}
	}

	// check ForbiddenHeaders
	for _, v := range r.ForbiddenHeaders {
		fhPair := strings.SplitN(v, ": ", 2)
		fieldName, fieldRegex := fhPair[0], fhPair[1]
		matchString, err := regexp.MatchString(fieldRegex, request.Header.Get(fieldName))
		if matchString || err != nil {
			return false
		}
	}

	// get body
	if request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return false
		}

		request.Body = ioutil.NopCloser(bytes.NewReader(body))

		// check MaxRequestLengthBytes
		if isBodyExceedLimits(body, r.MaxRequestLengthBytes) {
			return false
		}

		// check ForbiddenRequestRe
		if isBodyForbidden(string(body), r.ForbiddenRequestRe) {
			return false
		}
	}

	return true
}

func (r *Rule) isResponseAllowed(response *http.Response) bool {
	if r == nil {
		return true
	}

	// check ForbiddenResponseCodes
	for _, v := range r.ForbiddenResponseCodes {
		if response.StatusCode == v {
			return false
		}
	}

	// get body
	if response.Body != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false
		}

		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// check MaxRequestLengthBytes
		if isBodyExceedLimits(body, r.MaxResponseLengthBytes) {
			return false
		}

		// check ForbiddenRequestRe
		if isBodyForbidden(string(body), r.ForbiddenResponseRe) {
			return false
		}
	}

	return true
}

func getBlockedQuery() *http.Response {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("Forbidden")),
		StatusCode: 403,
	}
}

func isBodyExceedLimits(body []byte, limit int) bool {
	if limit <= 0 {
		return false
	}

	return len(body) > limit
}

func isBodyForbidden(body string, ForbiddenRes []string) bool {
	if len(body) <= 0 {
		return false
	}
	for _, v := range ForbiddenRes {
		matchString, err := regexp.MatchString(v, body)
		if matchString || err != nil {
			return true
		}
	}

	return false
}

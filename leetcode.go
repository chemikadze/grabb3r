package grabb3r

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type leetcodeSource struct {
	username         string
	password         string
	httpClient       *http.Client
	requestInterval  time.Duration
	rateLimitBackoff time.Duration
	pageSize         int
}

// Response type from leetcode submissions API
type leetcodeSubmissionsDump struct {
	SubmissionsDump []leetcodeSolution `json:"submissions_dump"`
	HasNext         bool               `json:"has_next"`
	LastKey         string             `json:"last_key"`
}

// Leetcode-specific error detail message
type errorDetail struct {
	Detail string
}

// Leetcode-specific Solution and SolutionDesc implementation
// Instead metadata list, API exposes full information about solution,
// including source code itself, so there is no need to do extra
// roundtrip to download solutions.
type leetcodeSolution struct {
	SubmissionId  int64 `json:"id"`
	Lang          string
	Time          string
	Timestamp     int64
	StatusDisplay string `json:"status_display"`
	Runtime       string
	Url           string
	IsPending     string `json:"is_pending"`
	Title         string
	Memory        string
	SolutionCode  string `json:"code"`
}

func (s *leetcodeSolution) String() string {
	return fmt.Sprintf("%v", s.SubmissionId)
}

func (s *leetcodeSolution) ProblemName() string {
	return s.Title
}

func (s *leetcodeSolution) Equals(other SolutionDesc) bool {
	if o, ok := other.(*leetcodeSolution); ok {
		return s.SubmissionId == o.SubmissionId
	} else {
		return false
	}
}

func (s *leetcodeSolution) Desc() SolutionDesc {
	return s
}

func (s *leetcodeSolution) Code() string {
	return s.SolutionCode
}

func (s *leetcodeSolution) Language() Language {
	return Language(s.Lang)
}

func (s *leetcodeSolution) SubmittedTime() time.Time {
	return time.Unix(s.Timestamp, 0)
}

// Leetcode-specific HTTP error type with extra details
type HttpError struct {
	Status int
	Body   string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("unexpected response: %v", e.Status)
}

// constructor for LeetCode SolutionSource
func NewLeetCodeSource(user string, password string) SolutionSource {
	src := &leetcodeSource{username: user, password: password}
	src.init()
	return src
}

func (lc *leetcodeSource) init() {
	if lc.httpClient != nil {
		return
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err) // typically, never should be null
	}
	lc.httpClient = &http.Client{Jar: jar}
	lc.requestInterval, _ = time.ParseDuration("2.5s")
	lc.rateLimitBackoff, _ = time.ParseDuration("15s")
	lc.pageSize = 20
}

func (lc *leetcodeSource) loginUrl() string {
	return "https://leetcode.com/accounts/login/"
}

func (lc *leetcodeSource) listSolutionsUrl(offset int, limit int, lastKey string) string {
	return fmt.Sprintf("https://leetcode.com/api/submissions/?offset=%v&limit=%v&lastkey=%v", offset, limit, lastKey)
}

func (lc *leetcodeSource) csrfToken(resp *http.Response) (string, error) {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrftoken" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("failed to get csrftoken from login page")
}

func (lc *leetcodeSource) makeHttpRequest(req *http.Request) (*http.Response, error) {
	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return resp, nil
}

func (lc *leetcodeSource) makeHttpRequestAndRead(req *http.Request) (*http.Response, []byte, error) {
	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, nil, err
	}
	return resp, data, nil
}

func (lc *leetcodeSource) Login() error {
	// get CSRF token
	req, err := http.NewRequest("GET", lc.loginUrl(), nil)
	if err != nil {
		return err
	}
	resp, err := lc.makeHttpRequest(req)
	if err != nil {
		return err
	}
	csrfToken, err := lc.csrfToken(resp)
	if err != nil {
		return err
	}

	// authorize
	formBody := url.Values{}
	formBody.Add("csrfmiddlewaretoken", csrfToken)
	formBody.Add("login", lc.username)
	formBody.Add("password", lc.password)
	formBody.Add("next", "/problems")
	req, err = http.NewRequest("POST", lc.loginUrl(), strings.NewReader(formBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("origin", "http://leetcode.com")
	req.Header.Add("referer", lc.loginUrl())
	req.Header.Add("x-csrftoken", csrfToken)
	resp, body, err := lc.makeHttpRequestAndRead(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("status: %v body:\n %v", resp.StatusCode, string(body))
		return fmt.Errorf("unexpected status code from login: %v %v", resp.StatusCode, resp.Status)
	}
	return nil
}

func isThrottled(statusCode int, errDetail string) bool {
	return statusCode == http.StatusForbidden &&
		errDetail == "You do not have permission to perform this action."
}

func (lc *leetcodeSource) ListSolutions() (chan SolutionDesc, chan error) {
	resChan := make(chan SolutionDesc)
	errChan := make(chan error)
	closeWithError := func(err error) {
		errChan <- err
		close(resChan)
		close(errChan)
	}
	go func() {
		lastKey := ""
		pageSize := lc.pageSize
		for hasNext, offset := true, 0; hasNext; {
			pageUrl := lc.listSolutionsUrl(offset, pageSize, lastKey)
			log.Printf("Requesting %s", pageUrl)
			req, err := http.NewRequest("GET", pageUrl, nil)
			if err != nil {
				closeWithError(err)
				return
			}
			// make request and validate response code, handling server-side throttling
			resp, body, err := lc.makeHttpRequestAndRead(req)
			if err != nil {
				closeWithError(err)
				return
			} else if resp.StatusCode != http.StatusOK {
				errDetail := &errorDetail{}
				err = json.Unmarshal(body, errDetail)
				if isThrottled(resp.StatusCode, errDetail.Detail) {
					time.Sleep(lc.rateLimitBackoff)
					log.Printf("retrying after throttling...")
					continue
				} else {
					err = HttpError{Status: resp.StatusCode, Body: string(body)}
					closeWithError(err)
					return
				}
			}
			// process successful response
			submissions := &leetcodeSubmissionsDump{}
			if err := json.Unmarshal(body, &submissions); err != nil {
				closeWithError(err)
				return
			}
			hasNext = submissions.HasNext
			lastKey = submissions.LastKey
			offset += pageSize
			for _, submission := range submissions.SubmissionsDump {
				resChan <- &submission
			}
			time.Sleep(lc.requestInterval)
		}
		close(resChan)
		close(errChan)
	}()
	return resChan, errChan
}

func (*leetcodeSource) GetSolution(id SolutionDesc) (Solution, error) {
	if ls, ok := id.(*leetcodeSolution); ok {
		return ls, nil
	} else {
		return nil, fmt.Errorf("not a leetcode solution: %v", id.String())
	}
}

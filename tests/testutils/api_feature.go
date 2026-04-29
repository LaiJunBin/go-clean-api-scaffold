package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/google/go-cmp/cmp"

	"go-clean-api-scaffold/internal/config"
)

type APIFeature struct {
	resp   *http.Response
	config *config.Config
}

func NewAPIFeature() *APIFeature {
	return &APIFeature{
		config: config.NewTestConfig(),
	}
}

func (a *APIFeature) buildAPIURL(endpoint string) string {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:" + strconv.Itoa(a.config.Server.Port),
	}

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		u.Path = endpoint
		return u.String()
	}

	u.Path = endpointURL.Path
	u.RawQuery = endpointURL.RawQuery
	return u.String()
}

func (a *APIFeature) ResetResponse(sc *godog.Scenario) {
	fmt.Printf("[feature][reset] scenario=%q id=%s\n", sc.Name, sc.Id)
	a.closeResponse()
}

func (a *APIFeature) closeResponse() {
	if a.resp != nil && a.resp.Body != nil {
		_ = a.resp.Body.Close()
	}
	a.resp = nil
}

func (a *APIFeature) ISendrequestTo(method, endpoint string) error {
	a.closeResponse()

	req, err := http.NewRequest(method, a.buildAPIURL(endpoint), nil)
	if err != nil {
		return err
	}

	a.resp, err = http.DefaultClient.Do(req)
	return err
}

func (a *APIFeature) TheResponseCodeShouldBe(code int) error {
	if a.resp == nil {
		return errors.New("response is nil, please send a request first")
	}

	if a.resp.StatusCode != code {
		return fmt.Errorf("expected response code %d, got %d", code, a.resp.StatusCode)
	}

	return nil
}

func (a *APIFeature) TheResponseShouldMatchJSON(body *godog.DocString) error {
	if a.resp == nil {
		return errors.New("response is nil, please send a request first")
	}

	actualBody, err := io.ReadAll(a.resp.Body)
	if err != nil {
		return err
	}

	var actual interface{}
	if err := json.Unmarshal(actualBody, &actual); err != nil {
		return fmt.Errorf("actual response is not valid json: %w", err)
	}

	var expected interface{}
	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return fmt.Errorf("expected response is not valid json: %w", err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		return fmt.Errorf("response json mismatch (-want +got):\n%s", diff)
	}

	return nil
}

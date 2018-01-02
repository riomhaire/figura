package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/riomhaire/figura/entities"
	"github.com/riomhaire/figura/frameworks"
	"github.com/riomhaire/figura/usecases"
)

func createTestRegistry(configurationStorage usecases.ConfigurationStorage, fileStorage usecases.Storage) *usecases.Registry {
	logger := frameworks.ConsoleLogger{}

	registry := usecases.Registry{}
	configuration := usecases.Configuration{}
	configuration.Version = "TEST"
	configuration.Application = "Figura"

	registry.Configuration = configuration
	registry.Logger = logger
	registry.ConfigurationStorage = configurationStorage
	registry.ConfigurationReader = usecases.NewConfigurationReader(&registry)

	registry.Storage = fileStorage

	return &registry
}

func TestUnknownApplication(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/configuration/unknown", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(UnknownApplicationStorage{}, AuthenticatedFileReader{})
	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleReadConfig)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestNotImplimented(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/configuration/unknown", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(UnknownApplicationStorage{}, AuthenticatedFileReader{})
	registry.ConfigurationReader = usecases.ConfigurationInteractor(NotImplimentedReader{})

	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleReadConfig)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestNotAuthenticated(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/configuration/known", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, AuthenticatedFileReader{})
	registry.ConfigurationReader = NotAuthentictedReader{}
	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleReadConfig)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestKnownApplication(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/configuration/known", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, AuthenticatedFileReader{})
	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleReadConfig)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestHealthHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v2/authentication/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, AuthenticatedFileReader{})
	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleHealth)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestStatisticsHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v2/authentication/statistics", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, AuthenticatedFileReader{})
	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleStatistics)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check for presence of common fields pid and uptime
	responseMap := make(map[string]interface{})
	err = json.NewDecoder(rr.Body).Decode(&responseMap)
	if err != nil {
		t.Fatal(err)
	}

	for _, val := range []string{"pid", "uptime"} {
		if _, ok := responseMap[val]; ok {
			//prsent
		} else {
			// missing - fail
			t.Fatal(errors.New("Expected parameter missing"))
		}

	}

}

func TestHostAppender(t *testing.T) {
	// Check that extra headers are added
	rr := httptest.NewRecorder()
	api := RestAPI{}
	api.AddWorkerHeader(rr, nil, nil)
	if len(rr.Header().Get("X-WORKER")) == 0 {
		t.Fail()
	}
	// Check that non-existent header not there
	if len(rr.Header().Get("X-WORKER2")) != 0 {
		t.Fail()
	}
}

func TestVersionAppender(t *testing.T) {
	// Check that extra headers are added
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, AuthenticatedFileReader{})
	api := NewRestAPI(registry)
	api.AddWorkerVersion(rr, nil, nil)
	if len(rr.Header().Get("X-WORKER-VERSION")) == 0 {
		t.Fail()
	}
	// Check that non-existent header not there
	if len(rr.Header().Get("X-WORKER2-VERSION")) != 0 {
		t.Fail()
	}
}

func TestCoorsAppender(t *testing.T) {
	// Check that extra headers are added
	rr := httptest.NewRecorder()
	api := RestAPI{}
	api.AddCoorsHeader(rr, nil, nil)
	if len(rr.Header().Get("Access-Control-Allow-Origin")) == 0 {
		t.Fail()
	}
}

func TestExtractInvalidAuthorizationHeader(t *testing.T) {
	// Parsing
	_, err := extractAuthorization("bad")
	if err == nil || !strings.Contains(err.Error(), "Not Authorized") {
		t.Fail()
	}
}

func TestExtractValidAuthorizationHeader(t *testing.T) {
	token := "VALID-TOKEN"
	validBearer := fmt.Sprintf("%v%v", bearerPrefix, token)
	bearer, err := extractAuthorization(validBearer)
	if err != nil {
		t.Fail()
	}
	if token != bearer {
		t.Fail()
	}
}

func TestExtractNonBearerAuthorizationHeader(t *testing.T) {
	token := "VALID-TOKEN"
	invalidBearer := fmt.Sprintf("%v%v", "something ", token)
	_, err := extractAuthorization(invalidBearer)
	if err == nil {
		t.Fail()
	}
}

func TestReadApplicationFile(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/configuration/known/somefile.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, AuthenticatedFileReader{})

	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleReadFile)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// OK check content is "this is a test"
	content := rr.Body.String()
	expected := "this is a test"
	if content != expected {
		t.Errorf("Content expected '%s' got '%s'", expected, content)
	}
}

func TestBadReadApplicationFile(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/api/v1/configuration/known/unknown.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	registry := createTestRegistry(KnownApplicationStorage{}, BadAuthenticatedFileReader{})

	restAPI := NewRestAPI(registry)
	handler := http.HandlerFunc(restAPI.HandleReadFile)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// *********************************************************************************
// Test implimentation to simulate back end

type UnknownApplicationStorage struct{}

func (t UnknownApplicationStorage) Lookup(application string) entities.ApplicationConfiguration {
	return entities.ApplicationConfiguration{
		ResultType: entities.UnknownApplication,
		Message:    fmt.Sprintf("Application '%v' is not registered within configuration service.", application),
	}
}

type KnownApplicationStorage struct{}

func (t KnownApplicationStorage) Lookup(application string) entities.ApplicationConfiguration {
	return entities.ApplicationConfiguration{
		ResultType: entities.NoError,
		Message:    fmt.Sprintf("Success"),
		Data:       []byte("entry: value"),
	}
}

type NotImplimentedReader struct{}

func (c NotImplimentedReader) Lookup(authorization, application string) entities.ApplicationConfiguration {
	return entities.ApplicationConfiguration{
		ResultType: entities.NotImplimentedError,
		Message:    fmt.Sprintf("Not Implimented"),
	}

}

type NotAuthentictedReader struct{}

func (c NotAuthentictedReader) Lookup(authorization, application string) entities.ApplicationConfiguration {
	return entities.ApplicationConfiguration{
		ResultType: entities.AuthenticationError,
		Message:    "Not Authenticated",
	}

}

type AuthenticatedFileReader struct{}

func (c AuthenticatedFileReader) Locate(application, filename string) (io.Reader, error) {

	return strings.NewReader("this is a test"), nil

}

type BadAuthenticatedFileReader struct{}

func (c BadAuthenticatedFileReader) Locate(application, filename string) (io.Reader, error) {

	return strings.NewReader(""), errors.New("Some error")

}

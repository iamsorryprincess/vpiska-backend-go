package v1

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/config"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/storage"
)

type testData struct {
	Name                string
	Url                 string
	Method              string
	Body                string
	RequestContentType  string
	ResponseContentType string
	ExpectedStatusCode  int
	CheckBody           bool
	ExpectedBody        string
	Handler             http.HandlerFunc
}

var testHandler *Handler

func TestHandler(t *testing.T) {
	defer os.RemoveAll("logs")
	defer os.RemoveAll("media")
	appLogger, logFile, err := logger.NewZeroLogger()
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	configuration, err := config.Parse("../../../../configs/test.yml")
	if err != nil {
		log.Fatal(err)
	}

	tokenManager := auth.NewJwtManager(configuration.JWT.Key, configuration.JWT.Issuer, configuration.JWT.Audience, time.Minute*5)
	repositories, cleaner, err := repository.NewRepositories(configuration.Database.ConnectionString, configuration.Database.DbName)
	if err != nil {
		log.Fatal(err)
	}

	defer cleaner.Clean()
	hashManager, err := hash.NewPasswordHashManager(configuration.Hash.Key)
	if err != nil {
		log.Fatal(err)
	}

	fileStorage, err := storage.NewLocalFileStorage("media")
	if err != nil {
		log.Fatal(err)
	}

	services, err := service.NewServices(appLogger, repositories, hashManager, tokenManager, fileStorage)
	if err != nil {
		log.Fatal(err)
	}

	testHandler = NewHandler(appLogger, services, tokenManager)
	t.Run("test users", TestUsers)
}

func testHandlerMethod(testData testData, t *testing.T) {
	request := httptest.NewRequest(testData.Method, testData.Url, bytes.NewBuffer([]byte(testData.Body)))
	request.Header.Set("Content-Type", testData.RequestContentType)
	recorder := httptest.NewRecorder()
	testData.Handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	if response.Header.Get("Content-Type") != testData.ResponseContentType {
		t.Fatalf("Incorrect Content-Type\nExpected Content-Type: %s\nActual Content-Type: %s\n", testData.ResponseContentType, response.Header.Get("Content-Type"))
	}

	if response.StatusCode != testData.ExpectedStatusCode {
		t.Fatalf("Incorrect status code\nExpected status code: %d\nActual status code: %d\n", testData.ExpectedStatusCode, response.StatusCode)
	}

	if testData.CheckBody {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		result := string(data)
		if result != testData.ExpectedBody {
			t.Errorf("Incorrect response body\nExpected body: %s\nActual body: %s", testData.ExpectedBody, result)
		}
	}
}

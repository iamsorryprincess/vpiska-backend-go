package v1

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/config"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
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
	AuthHeader          string
	RequestContentType  string
	ResponseContentType string
	ExpectedStatusCode  int
	CheckBody           bool
	ExpectedBody        string
	Handler             http.HandlerFunc
}

var testHandler *Handler
var testUserAccessToken string
var testEventId string
var testEventId10 string
var testEventId25 string
var testEventId50 string
var testEventId75 string
var testEventId100 string

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

	password, err := hashManager.HashPassword("string")
	if err != nil {
		log.Fatal(err)
	}

	userId, err := repositories.Users.CreateUser(context.Background(), domain.User{
		Name:     "integration_tests",
		Phone:    "1111111111",
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	testUserAccessToken, err = tokenManager.GetAccessToken(auth.TokenData{ID: userId, Name: "integration_tests"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = repositories.Users.CreateUser(context.Background(), domain.User{
		Name:     "integration_tests_events",
		Phone:    "9090909090",
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	testEventId, err = repositories.Events.CreateEvent(context.Background(), domain.Event{
		OwnerID: "owner_id",
		Name:    "integration_tests",
		Address: "integration_tests",
		State:   domain.EventStateOpened,
		Coordinates: domain.Coordinates{
			X: 99999,
			Y: 99999,
		},
		Users: []domain.UserInfo{
			{ID: "owner_id"},
		},
		Media: []domain.MediaInfo{
			{
				ID:          "media_id",
				ContentType: "image/jpeg",
			},
		},
		ChatMessages: []domain.ChatMessage{
			{
				UserID:   "test_id_1",
				UserName: "test_events_1",
				Message:  "test message",
			},
			{
				UserID:   "test_id_2",
				UserName: "test_events_2",
				Message:  "another message",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = initEventsForRangeTest(repositories.Events)
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
	t.Run("users", TestUsers)
	t.Run("events", TestEvents)
}

func testHandlerMethod(testData testData, t *testing.T) {
	request := httptest.NewRequest(testData.Method, testData.Url, bytes.NewBuffer([]byte(testData.Body)))
	request.Header.Set("Content-Type", testData.RequestContentType)
	if testData.AuthHeader != "" {
		request.Header.Set("Authorization", "Bearer "+testData.AuthHeader)
	}
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

func initEventsForRangeTest(events repository.Events) error {
	id, err := events.CreateEvent(context.Background(), domain.Event{
		OwnerID: "owner_id_range_1",
		Name:    "integration_tests_range_1",
		Address: "integration_tests_range_1",
		State:   domain.EventStateOpened,
		Coordinates: domain.Coordinates{
			X: 10,
			Y: 10,
		},
		Users:        []domain.UserInfo{},
		Media:        []domain.MediaInfo{},
		ChatMessages: []domain.ChatMessage{},
	})
	if err != nil {
		return err
	}
	testEventId10 = id

	testEventId25, err = events.CreateEvent(context.Background(), domain.Event{
		OwnerID: "owner_id_range_2",
		Name:    "integration_tests_range_2",
		Address: "integration_tests_range_2",
		State:   domain.EventStateOpened,
		Coordinates: domain.Coordinates{
			X: 25,
			Y: 25,
		},
		Users:        []domain.UserInfo{},
		Media:        []domain.MediaInfo{},
		ChatMessages: []domain.ChatMessage{},
	})
	if err != nil {
		return err
	}

	testEventId50, err = events.CreateEvent(context.Background(), domain.Event{
		OwnerID: "owner_id_range_3",
		Name:    "integration_tests_range_3",
		Address: "integration_tests_range_3",
		State:   domain.EventStateOpened,
		Coordinates: domain.Coordinates{
			X: 50,
			Y: 50,
		},
		Users:        []domain.UserInfo{},
		Media:        []domain.MediaInfo{},
		ChatMessages: []domain.ChatMessage{},
	})
	if err != nil {
		return err
	}

	testEventId75, err = events.CreateEvent(context.Background(), domain.Event{
		OwnerID: "owner_id_range_4",
		Name:    "integration_tests_range_4",
		Address: "integration_tests_range_4",
		State:   domain.EventStateOpened,
		Coordinates: domain.Coordinates{
			X: 75,
			Y: 75,
		},
		Users:        []domain.UserInfo{},
		Media:        []domain.MediaInfo{},
		ChatMessages: []domain.ChatMessage{},
	})
	if err != nil {
		return err
	}

	testEventId100, err = events.CreateEvent(context.Background(), domain.Event{
		OwnerID: "owner_id_range_5",
		Name:    "integration_tests_range_5",
		Address: "integration_tests_range_5",
		State:   domain.EventStateOpened,
		Coordinates: domain.Coordinates{
			X: 100,
			Y: 100,
		},
		Users:        []domain.UserInfo{},
		Media:        []domain.MediaInfo{},
		ChatMessages: []domain.ChatMessage{},
	})
	if err != nil {
		return err
	}

	return nil
}

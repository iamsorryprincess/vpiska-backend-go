package v1

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEvents(t *testing.T) {
	t.Run("create", TestCreateEvent)
	t.Run("get", TestGetEvent)
	t.Run("range", TestRangeEvents)
}

func TestCreateEvent(t *testing.T) {
	tests := []testData{
		{
			Name:                "unauthorized",
			Url:                 "/api/v1/events/create",
			Method:              http.MethodPost,
			Body:                `{}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"unauthorized"}],"result":null}`,
			Handler:             testHandler.jwtAuth(testHandler.createEvent),
		},
		{
			Name:                "empty body",
			Url:                 "/api/v1/events/create",
			Method:              http.MethodPost,
			Body:                `{}`,
			AuthHeader:          testUserAccessToken,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"NameIsEmpty"},{"errorCode":"AddressIsEmpty"},{"errorCode":"CoordinatesIsEmpty"}],"result":null}`,
			Handler:             testHandler.jwtAuth(testHandler.createEvent),
		},
		{
			Name:                "success",
			Url:                 "/api/v1/events/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test_events_create","address":"test_events_create","coordinates":{"x":99999,"y":99999}}`,
			AuthHeader:          testUserAccessToken,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			Handler:             testHandler.jwtAuth(testHandler.createEvent),
		},
		{
			Name:                "owner already has event",
			Url:                 "/api/v1/events/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test_events_create","address":"test_events_create","coordinates":{"x":0,"y":0}}`,
			AuthHeader:          testUserAccessToken,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"OwnerAlreadyHasEvent"}],"result":null}`,
			Handler:             testHandler.jwtAuth(testHandler.createEvent),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testHandlerMethod(test, t)
		})
	}
}

func TestGetEvent(t *testing.T) {
	tests := []testData{
		{
			Name:                "empty body",
			Url:                 "/api/v1/events/get",
			Method:              http.MethodPost,
			Body:                `{}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"IdIsEmpty"}],"result":null}`,
			Handler:             testHandler.getEventByID,
		},
		{
			Name:                "invalid id",
			Url:                 "/api/v1/events/get",
			Method:              http.MethodPost,
			Body:                `{"eventId":"123"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"InvalidIdFormat"}],"result":null}`,
			Handler:             testHandler.getEventByID,
		},
		{
			Name:                "success",
			Url:                 "/api/v1/events/get",
			Method:              http.MethodPost,
			Body:                fmt.Sprintf(`{"eventId":"%s"}`, testEventId),
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":{"id":"%s","ownerId":"owner_id","name":"integration_tests","address":"integration_tests","coordinates":{"x":99999,"y":99999},"usersCount":1,"media":[{"id":"media_id","contentType":"image/jpeg"}],"chatMessages":[{"userId":"test_id_1","userName":"test_events_1","userImageId":"","message":"test message"},{"userId":"test_id_2","userName":"test_events_2","userImageId":"","message":"another message"}]}}`, testEventId),
			Handler:             testHandler.getEventByID,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testHandlerMethod(test, t)
		})
	}
}

func TestRangeEvents(t *testing.T) {
	tests := []testData{
		{
			Name:                "empty result",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 0,"y": 0},"horizontalRange": 0,"verticalRange": 0}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":true,"errors":null,"result":[]}`,
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 20 coords 0",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 0,"y": 0},"horizontalRange": 20,"verticalRange": 20}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_1","usersCount":0,"coordinates":{"x":10,"y":10}}]}`, testEventId10),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 52 coords 0",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 0,"y": 0},"horizontalRange": 52,"verticalRange": 52}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_1","usersCount":0,"coordinates":{"x":10,"y":10}},{"id":"%s","name":"integration_tests_range_2","usersCount":0,"coordinates":{"x":25,"y":25}}]}`, testEventId10, testEventId25),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 102 coords 0",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 0,"y": 0},"horizontalRange": 102,"verticalRange": 102}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_1","usersCount":0,"coordinates":{"x":10,"y":10}},{"id":"%s","name":"integration_tests_range_2","usersCount":0,"coordinates":{"x":25,"y":25}},{"id":"%s","name":"integration_tests_range_3","usersCount":0,"coordinates":{"x":50,"y":50}}]}`, testEventId10, testEventId25, testEventId50),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 152 coords 0",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 0,"y": 0},"horizontalRange": 152,"verticalRange": 152}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_1","usersCount":0,"coordinates":{"x":10,"y":10}},{"id":"%s","name":"integration_tests_range_2","usersCount":0,"coordinates":{"x":25,"y":25}},{"id":"%s","name":"integration_tests_range_3","usersCount":0,"coordinates":{"x":50,"y":50}},{"id":"%s","name":"integration_tests_range_4","usersCount":0,"coordinates":{"x":75,"y":75}}]}`, testEventId10, testEventId25, testEventId50, testEventId75),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 202 coords 0",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 0,"y": 0},"horizontalRange": 202,"verticalRange": 202}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_1","usersCount":0,"coordinates":{"x":10,"y":10}},{"id":"%s","name":"integration_tests_range_2","usersCount":0,"coordinates":{"x":25,"y":25}},{"id":"%s","name":"integration_tests_range_3","usersCount":0,"coordinates":{"x":50,"y":50}},{"id":"%s","name":"integration_tests_range_4","usersCount":0,"coordinates":{"x":75,"y":75}},{"id":"%s","name":"integration_tests_range_5","usersCount":0,"coordinates":{"x":100,"y":100}}]}`, testEventId10, testEventId25, testEventId50, testEventId75, testEventId100),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 25 coords 13",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 13,"y": 13},"horizontalRange": 25,"verticalRange": 25}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_1","usersCount":0,"coordinates":{"x":10,"y":10}},{"id":"%s","name":"integration_tests_range_2","usersCount":0,"coordinates":{"x":25,"y":25}}]}`, testEventId10, testEventId25),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 2 coords 50",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 50,"y": 50},"horizontalRange": 2,"verticalRange": 2}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_3","usersCount":0,"coordinates":{"x":50,"y":50}}]}`, testEventId50),
			Handler:             testHandler.getEventsByRange,
		},
		{
			Name:                "range 52 coords 50",
			Url:                 "/api/v1/events/range",
			Method:              http.MethodPost,
			Body:                `{"coordinates": {"x": 50,"y": 50},"horizontalRange": 52,"verticalRange": 52}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        fmt.Sprintf(`{"isSuccess":true,"errors":null,"result":[{"id":"%s","name":"integration_tests_range_2","usersCount":0,"coordinates":{"x":25,"y":25}},{"id":"%s","name":"integration_tests_range_3","usersCount":0,"coordinates":{"x":50,"y":50}},{"id":"%s","name":"integration_tests_range_4","usersCount":0,"coordinates":{"x":75,"y":75}}]}`, testEventId25, testEventId50, testEventId75),
			Handler:             testHandler.getEventsByRange,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testHandlerMethod(test, t)
		})
	}
}

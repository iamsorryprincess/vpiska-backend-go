package v1

import (
	"net/http"
	"testing"
)

func TestUsers(t *testing.T) {
	t.Run("create user", TestCreateUser)
	t.Run("login user", TestLoginUser)
}

func TestCreateUser(t *testing.T) {
	tests := []testData{
		{
			Name:                "empty phone",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name": "test","password": "123456","confirmPassword": "123456"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PhoneIsEmpty"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "invalid phone format",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name": "test","phone":"123","password": "123456","confirmPassword": "123456"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PhoneRegexInvalid"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "empty name",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"phone":"9374113516","password": "123456","confirmPassword": "123456"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"NameIsEmpty"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "empty password",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test","phone":"9374113516"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PasswordIsEmpty"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "password length invalid",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test","phone":"9374113516","password":"123"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PasswordLengthInvalid"},{"errorCode":"ConfirmPasswordInvalid"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "empty string",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                ``,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"InternalError"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "empty body",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"NameIsEmpty"},{"errorCode":"PhoneIsEmpty"},{"errorCode":"PasswordIsEmpty"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "success create",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test","phone":"9374113516","password":"string","confirmPassword":"string"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           false,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "name and phone already exist",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test","phone":"9374113516","password":"string","confirmPassword":"string"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"NameAndPhoneAlreadyUse"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "name already exist",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name":"test","phone":"9374113515","password":"string","confirmPassword":"string"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"NameAlreadyUse"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
		{
			Name:                "phone already exist",
			Url:                 "/api/v1/users/create",
			Method:              http.MethodPost,
			Body:                `{"name":"qweasd","phone":"9374113516","password":"string","confirmPassword":"string"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PhoneAlreadyUse"}],"result":null}`,
			Handler:             testHandler.createUser,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testHandlerMethod(test, t)
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []testData{
		{
			Name:                "empty phone",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                `{"password": "123456"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PhoneIsEmpty"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
		{
			Name:                "empty password",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                `{"phone":"9090909090"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PasswordIsEmpty"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
		{
			Name:                "password length invalid",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                `{"phone":"9374113516","password":"123"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PasswordLengthInvalid"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
		{
			Name:                "empty string",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                ``,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"InternalError"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
		{
			Name:                "empty body",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                `{}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"PhoneIsEmpty"},{"errorCode":"PasswordIsEmpty"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
		{
			Name:                "user not found",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                `{"phone":"1111111100","password":"string"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"UserNotFound"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
		{
			Name:                "invalid password",
			Url:                 "/api/v1/users/login",
			Method:              http.MethodPost,
			Body:                `{"phone":"1111111111","password":"stringss"}`,
			RequestContentType:  contentTypeJSON,
			ResponseContentType: contentTypeJSON,
			ExpectedStatusCode:  http.StatusOK,
			CheckBody:           true,
			ExpectedBody:        `{"isSuccess":false,"errors":[{"errorCode":"InvalidPassword"}],"result":null}`,
			Handler:             testHandler.loginUser,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			testHandlerMethod(test, t)
		})
	}
}

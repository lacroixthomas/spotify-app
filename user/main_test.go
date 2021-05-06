package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zmb3/spotify"
)

type mockSpotifyClient struct {
	err  error
	user spotify.User
}

func (c *mockSpotifyClient) GetUsersPublicProfile(userID spotify.ID) (*spotify.User, error) {
	return &c.user, c.err
}

func (c *mockSpotifyClient) CurrentUser() (*spotify.PrivateUser, error) {
	var pu spotify.PrivateUser
	pu.DisplayName = c.user.DisplayName
	pu.ID = c.user.ID
	pu.Images = c.user.Images

	return &pu, c.err
}

func getRequestMock(err error, user spotify.User) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, CLIENT_CONTEXT, &mockSpotifyClient{
		err:  err,
		user: user,
	})
	return r.WithContext(ctx)
}

func Test_userHandler(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "should get user properly",
			req: getRequestMock(nil, spotify.User{
				DisplayName: "thomas",
				ID:          "thomas123",
				Images:      []spotify.Image{},
			}),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"thomas","id":"thomas123","image":""}`,
		},
		{
			name: "should get user properly with image",
			req: getRequestMock(nil, spotify.User{
				DisplayName: "thomas",
				ID:          "thomas123",
				Images:      []spotify.Image{{URL: "url"}},
			}),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"thomas","id":"thomas123","image":"url"}`,
		},
		{
			name:         "should error on spotify api call",
			req:          getRequestMock(errors.New("cannot get user"), spotify.User{}),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			userHandler(rr, tt.req)
			if res := rr.Code; res != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					res, tt.expectedCode)
			}
			if tt.expectedCode < 500 {
				if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						strings.TrimSpace(rr.Body.String()), tt.expectedBody)
				}
			}
		})
	}
}

func Test_userFromHandler(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name: "should get specific user properly",
			req: getRequestMock(nil, spotify.User{
				DisplayName: "thomas",
				ID:          "thomas123",
				Images:      []spotify.Image{{URL: "url"}},
			}),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"thomas","id":"thomas123","image":"url"}`,
		},
		{
			name:         "should error on spotify api call",
			req:          getRequestMock(errors.New("cannot get user"), spotify.User{}),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			userFromHandler(rr, tt.req)
			if res := rr.Code; res != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					res, tt.expectedCode)
			}
			if tt.expectedCode < 500 {
				if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v",
						strings.TrimSpace(rr.Body.String()), tt.expectedBody)
				}
			}
		})
	}
}

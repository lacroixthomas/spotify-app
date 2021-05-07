package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/zmb3/spotify"
)

type mockSpotifyClient struct {
	err      error
	playlist spotify.SimplePlaylistPage
}

func (c *mockSpotifyClient) CurrentUsersPlaylists() (*spotify.SimplePlaylistPage, error) {
	return &c.playlist, c.err
}

func getRequestMock(err error, playlist spotify.SimplePlaylistPage) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, CLIENT_CONTEXT, &mockSpotifyClient{
		err:      err,
		playlist: playlist,
	})
	return r.WithContext(ctx)
}

func Test_reducePlaylist(t *testing.T) {
	type args struct {
		playlistResp *spotify.SimplePlaylistPage
	}
	tests := []struct {
		name string
		args args
		want []playlistItem
	}{
		{
			name: "should handle empty playlists",
			args: args{
				playlistResp: &spotify.SimplePlaylistPage{
					Playlists: []spotify.SimplePlaylist{},
				},
			},
			want: []playlistItem{},
		},
		{
			name: "should get all playlists",
			args: args{
				playlistResp: &spotify.SimplePlaylistPage{
					Playlists: []spotify.SimplePlaylist{
						{
							Name: "test-name",
							ID:   "ID",
							URI:  "uri:...",
							Owner: spotify.User{
								DisplayName: "Thomas",
							},
						},
					},
				},
			},
			want: []playlistItem{
				{
					Name:      "test-name",
					OwnerName: "Thomas",
					ID:        "ID",
					URI:       "uri:...",
				},
			},
		},
		{
			name: "should get all the playlists with images",
			args: args{
				playlistResp: &spotify.SimplePlaylistPage{
					Playlists: []spotify.SimplePlaylist{
						{
							Name: "test-name",
							ID:   "ID",
							URI:  "uri:...",
							Owner: spotify.User{
								DisplayName: "Thomas",
							},
						},
						{
							Name: "test-name-2",
							ID:   "ID-2",
							URI:  "uri:2",
							Owner: spotify.User{
								DisplayName: "Thomas2",
							},
							Images: []spotify.Image{
								{
									URL: "http://...",
								},
							},
						},
					},
				},
			},
			want: []playlistItem{
				{
					Name:      "test-name",
					OwnerName: "Thomas",
					ID:        "ID",
					URI:       "uri:...",
				},
				{
					Name:      "test-name-2",
					OwnerName: "Thomas2",
					ID:        "ID-2",
					URI:       "uri:2",
					Image:     "http://...",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reducePlaylist(tt.args.playlistResp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reducePlaylist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_playlistHandler(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
		expectedBody string
	}{
		{
			name:         "should get 0 playlist",
			args:         args{req: getRequestMock(nil, spotify.SimplePlaylistPage{})},
			expectedCode: http.StatusOK,
			expectedBody: `[]`,
		},
		{
			name: "should get playlists properly",
			args: args{req: getRequestMock(nil, spotify.SimplePlaylistPage{
				Playlists: []spotify.SimplePlaylist{
					{
						Name: "test-name",
						ID:   "ID",
						URI:  "uri:...",
						Owner: spotify.User{
							DisplayName: "Thomas",
						},
					},
				},
			})},
			expectedCode: http.StatusOK,
			expectedBody: `[{"image":"","name":"test-name","owner_name":"Thomas","ID":"ID","uri":"uri:..."}]`,
		},
		{
			name: "should get playlists properly with images",
			args: args{req: getRequestMock(nil, spotify.SimplePlaylistPage{
				Playlists: []spotify.SimplePlaylist{
					{
						Name: "test-name",
						ID:   "ID",
						URI:  "uri:...",
						Owner: spotify.User{
							DisplayName: "Thomas",
						},
						Images: []spotify.Image{
							{
								URL: "http://...",
							},
						},
					},
				},
			})},
			expectedCode: http.StatusOK,
			expectedBody: `[{"image":"http://...","name":"test-name","owner_name":"Thomas","ID":"ID","uri":"uri:..."}]`,
		},
		{
			name:         "should error on spotify api call",
			args:         args{req: getRequestMock(errors.New("could not get playlists"), spotify.SimplePlaylistPage{})},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `[]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			playlistHandler(rr, tt.args.req)
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

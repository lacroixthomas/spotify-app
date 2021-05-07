package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/zmb3/spotify"
)

func getTimeFromString(str string) time.Time {
	date, _ := time.Parse(spotify.DateLayout, str)
	return date
}

type mockSpotifyClient struct {
	err    error
	player spotify.CurrentlyPlaying
}

func (c *mockSpotifyClient) PlayerCurrentlyPlaying() (*spotify.CurrentlyPlaying, error) {
	return &c.player, c.err
}

func (c *mockSpotifyClient) PlayOpt(opt *spotify.PlayOptions) error {
	return c.err
}

func (c *mockSpotifyClient) Pause() error {
	return c.err
}
func (c *mockSpotifyClient) Next() error {
	return c.err
}
func (c *mockSpotifyClient) Previous() error {
	return c.err
}

func getRequestMock(err error, player spotify.CurrentlyPlaying, withBody bool) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	if withBody {
		r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	}
	ctx := r.Context()
	ctx = context.WithValue(ctx, CLIENT_CONTEXT, &mockSpotifyClient{
		err:    err,
		player: player,
	})
	return r.WithContext(ctx)
}

func Test_reducePlayer(t *testing.T) {
	type args struct {
		playerResp *spotify.CurrentlyPlaying
	}
	tests := []struct {
		name string
		args args
		want player
	}{
		{
			name: "",
			args: args{
				&spotify.CurrentlyPlaying{
					Playing: true,
					Item: &spotify.FullTrack{
						Album: spotify.SimpleAlbum{
							Name:                 "album",
							ReleaseDate:          "2020-12-15",
							ReleaseDatePrecision: "day",
						},
						SimpleTrack: spotify.SimpleTrack{
							Name: "test",
							ID:   "id",
							Artists: []spotify.SimpleArtist{
								{
									Name: "artist name",
								},
								{
									Name: "thomas",
								},
							},
						},
					},
				},
			},
			want: player{
				IsPlaying:   true,
				AlbumName:   "album",
				ArtistsName: []string{"artist name", "thomas"},
				MusicName:   "test",
				ID:          "id",
				ReleaseDate: getTimeFromString("2020-12-15"),
				Duration:    0,
				Progress:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reducePlayer(tt.args.playerResp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reducePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_playerHandler(t *testing.T) {
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
			name: "should get the current player",
			args: args{
				req: getRequestMock(nil, spotify.CurrentlyPlaying{
					Playing: true,
					Item: &spotify.FullTrack{
						Album: spotify.SimpleAlbum{
							Name:                 "album",
							ReleaseDate:          "2020-12-15",
							ReleaseDatePrecision: "day",
						},
						SimpleTrack: spotify.SimpleTrack{
							Name: "test",
							ID:   "id",
							Artists: []spotify.SimpleArtist{
								{
									Name: "artist name",
								},
								{
									Name: "thomas",
								},
							},
						},
					},
				}, false),
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"is_playing":true,"album_name":"album","artists_name":["artist name","thomas"],"music_name":"test","ID":"id","release_date":"2020-12-15T00:00:00Z","progress":0,"duration":0}`,
		},
		{
			name: "should error on spotify api call",
			args: args{
				req: getRequestMock(errors.New("could not fetch player"), spotify.CurrentlyPlaying{}, false),
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			playerHandler(rr, tt.args.req)
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

func Test_playMusicHandler(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name:         "should start the music",
			args:         args{req: getRequestMock(nil, spotify.CurrentlyPlaying{}, true)},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should error decoding body",
			args:         args{req: getRequestMock(errors.New("could not start the music"), spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "should error on spotify api call",
			args:         args{req: getRequestMock(errors.New("could not start the music"), spotify.CurrentlyPlaying{}, true)},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			playMusicHandler(rr, tt.args.req)
			if res := rr.Code; res != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					res, tt.expectedCode)
			}
		})
	}
}

func Test_pauseMusicHandler(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name:         "should pause the music",
			args:         args{req: getRequestMock(nil, spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should error on spotify api call",
			args:         args{req: getRequestMock(errors.New("could not pause music"), spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			pauseMusicHandler(rr, tt.args.req)
			if res := rr.Code; res != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					res, tt.expectedCode)
			}
		})
	}
}

func Test_nextMusicHandler(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name:         "should go to next music",
			args:         args{req: getRequestMock(nil, spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should error on spotify api call",
			args:         args{req: getRequestMock(errors.New("could not go to next music"), spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			nextMusicHandler(rr, tt.args.req)
			if res := rr.Code; res != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					res, tt.expectedCode)
			}
		})
	}
}

func Test_prevMusicHandler(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name:         "should go to previous music",
			args:         args{req: getRequestMock(nil, spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusOK,
		},
		{
			name:         "should error on spotify api call",
			args:         args{req: getRequestMock(errors.New("could not go to previous music"), spotify.CurrentlyPlaying{}, false)},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			prevMusicHandler(rr, tt.args.req)
			if res := rr.Code; res != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					res, tt.expectedCode)
			}
		})
	}
}

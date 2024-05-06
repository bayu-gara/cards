package rest

import (
	"net/http"
	"net/url"
	"testing"

	//usecase
	deckuc "github.com/bayu-gara/cards/internal/usecase/deck"

	//external
	"github.com/golang/mock/gomock"
)

type CustomResponseWriter struct{}

func (w *CustomResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
}

func Test_createDeck(t *testing.T) {
	ctl := gomock.NewController(t)
	deckUsecaseMock := deckuc.NewMockUsecase(ctl)
	defer ctl.Finish()

	deckuc.UsecaseObj = deckUsecaseMock
	urlObj, _ := url.Parse("http://localhost:8080/v1/deck/create?shuffle=false&cards=AS,AH,AC,AD")

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Success create a deck",
			args: args{
				w: &CustomResponseWriter{},
				r: &http.Request{
					URL: urlObj,
				},
			},
			mock: func() {
				deckUsecaseMock.EXPECT().CreateDeck(
					gomock.Any(),
					false,
					[]string{"AS", "AH", "AC", "AD"},
				).Return(deckuc.Deck{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			createDeck(tt.args.w, tt.args.r)
		})
	}
}

func Test_opendDeck(t *testing.T) {
	ctl := gomock.NewController(t)
	deckUsecaseMock := deckuc.NewMockUsecase(ctl)
	defer ctl.Finish()

	deckuc.UsecaseObj = deckUsecaseMock
	urlObj, _ := url.Parse("http://localhost:8080/v1/deck/open?deck_id=fa704a87-b45f-40a4-bf96-93c3f3b0f33b")

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Success open a deck",
			args: args{
				w: &CustomResponseWriter{},
				r: &http.Request{
					URL: urlObj,
				},
			},
			mock: func() {
				deckUsecaseMock.EXPECT().OpenDeck(
					gomock.Any(),
					"fa704a87-b45f-40a4-bf96-93c3f3b0f33b",
				).Return(deckuc.Deck{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			opendDeck(tt.args.w, tt.args.r)
		})
	}
}

func Test_drawCard(t *testing.T) {
	ctl := gomock.NewController(t)
	deckUsecaseMock := deckuc.NewMockUsecase(ctl)
	defer ctl.Finish()

	deckuc.UsecaseObj = deckUsecaseMock
	urlObj, _ := url.Parse("http://localhost:8080/v1/deck/draw?deck_id=fa704a87-b45f-40a4-bf96-93c3f3b0f33b&count=2")

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		mock func()
	}{
		{
			name: "Success draw a card",
			args: args{
				w: &CustomResponseWriter{},
				r: &http.Request{
					URL: urlObj,
				},
			},
			mock: func() {
				deckUsecaseMock.EXPECT().DrawCard(
					gomock.Any(),
					"fa704a87-b45f-40a4-bf96-93c3f3b0f33b",
					2,
				).Return(deckuc.Deck{}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			drawCard(tt.args.w, tt.args.r)
		})
	}
}

func Test_writeSuccess(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		data interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Nil data",
			args: args{
				w:    &CustomResponseWriter{},
				data: nil,
			},
		},
		{
			name: "Default",
			args: args{
				w:    &CustomResponseWriter{},
				data: `{"name": "Bayu Anggara"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeSuccess(tt.args.w, tt.args.data)
		})
	}
}

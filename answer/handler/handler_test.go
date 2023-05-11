package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sivchari/cagows/answer/handler"
	"github.com/sivchari/cagows/answer/memory"
	"github.com/sivchari/cagows/answer/model"
	"github.com/sivchari/cagows/answer/repository"
	"github.com/sivchari/cagows/answer/router"
)

func Test_handler_List(t *testing.T) {
	type fields struct {
		repo repository.Repository
		hook func(repo repository.Repository)
	}
	tests := []struct {
		name   string
		fields fields
		want   []*model.Todo
	}{
		{
			name: "正常系: リスト取得",
			fields: fields{
				repo: memory.New(),
				hook: func(repo repository.Repository) {
					repo.Add(&model.Todo{
						Title:     "test",
						Completed: false,
					})
					repo.Add(&model.Todo{
						Title:     "test2",
						Completed: false,
					})
				},
			},
			want: []*model.Todo{
				{
					ID:        1,
					Title:     "test",
					Completed: false,
				},
				{
					ID:        2,
					Title:     "test2",
					Completed: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.Routing(tt.fields.repo)
			tsrv := httptest.NewServer(r)
			t.Cleanup(tsrv.Close)
			tt.fields.hook(tt.fields.repo)
			res, err := http.Get(fmt.Sprintf("%s/list", tsrv.URL))
			if err != nil {
				t.Errorf("http.Get() error = %v", err)
				return
			}
			if res.StatusCode != http.StatusOK {
				t.Errorf("http.Get() status code = %v, want %v", res.StatusCode, http.StatusOK)
			}
			var got []*model.Todo
			readedRes, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll() error = %v", err)
				return
			}
			if err := json.Unmarshal(readedRes, &got); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("handler.List() mismatch (-got +want):\n%s", diff)
				return
			}
			res.Body.Close()
		})
	}
}

func Test_handler_Get(t *testing.T) {
	type fields struct {
		repo repository.Repository
		hook func(repo repository.Repository)
	}
	tests := []struct {
		name   string
		id     int
		fields fields
		want   *model.Todo
	}{
		{
			name: "正常系: ID指定で取得",
			id:   1,
			fields: fields{
				repo: memory.New(),
				hook: func(repo repository.Repository) {
					repo.Add(&model.Todo{
						Title:     "test",
						Completed: false,
					})
				},
			},
			want: &model.Todo{
				ID:        1,
				Title:     "test",
				Completed: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.Routing(tt.fields.repo)
			tsrv := httptest.NewServer(r)
			t.Cleanup(tsrv.Close)
			tt.fields.hook(tt.fields.repo)
			res, err := http.Get(fmt.Sprintf("%s/get?id=%d", tsrv.URL, tt.id))
			if err != nil {
				t.Errorf("http.Get() error = %v", err)
				return
			}
			if res.StatusCode != http.StatusOK {
				t.Errorf("http.Get() status code = %v, want %v", res.StatusCode, http.StatusOK)
			}
			var got *model.Todo
			readedRes, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll() error = %v", err)
				return
			}
			if err := json.Unmarshal(readedRes, &got); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("handler.Get() mismatch (-got +want):\n%s", diff)
				return
			}
			res.Body.Close()
		})
	}
}

func Test_handler_Add(t *testing.T) {
	type fields struct {
		repo repository.Repository
	}
	tests := []struct {
		name       string
		fields     fields
		want       *model.Todo
		wantStatus int
	}{
		{
			name: "正常系: 追加",
			fields: fields{
				repo: memory.New(),
			},
			want: &model.Todo{
				ID:        1,
				Title:     "test",
				Completed: false,
			},
			wantStatus: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.Routing(tt.fields.repo)
			tsrv := httptest.NewServer(r)
			t.Cleanup(tsrv.Close)
			res, err := http.Post(
				fmt.Sprintf("%s/add", tsrv.URL),
				"application/json",
				strings.NewReader(`{"title":"test"}`),
			)
			if err != nil {
				t.Errorf("http.Get() error = %v", err)
				return
			}
			if res.StatusCode != tt.wantStatus {
				t.Errorf("http.Get() status code = %v, want %v", res.StatusCode, tt.want)
			}
			var got *model.Todo
			readedRes, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll() error = %v", err)
				return
			}
			res.Body.Close()
			if err := json.Unmarshal(readedRes, &got); err != nil {
				t.Errorf("json.Unmarshal() error = %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("handler.Add() mismatch (-got +want):\n%s", diff)
				return
			}
		})
	}
}

func Test_handler_Complete(t *testing.T) {
	type fields struct {
		repo repository.Repository
		hook func(repo repository.Repository)
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "正常系: 完了",
			fields: fields{
				repo: memory.New(),
				hook: func(repo repository.Repository) {
					repo.Add(&model.Todo{
						Title:     "test",
						Completed: false,
					})
				},
			},
			want: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.Routing(tt.fields.repo)
			tsrv := httptest.NewServer(r)
			t.Cleanup(tsrv.Close)
			tt.fields.hook(tt.fields.repo)
			res, err := http.Post(
				fmt.Sprintf("%s/complete", tsrv.URL),
				"application/json",
				strings.NewReader(`{"id":1}`),
			)
			if err != nil {
				t.Errorf("http.Get() error = %v", err)
				return
			}
			if res.StatusCode != tt.want {
				t.Errorf("http.Get() status code = %v, want %v", res.StatusCode, http.StatusOK)
			}
			res.Body.Close()
		})
	}
}

func Test_GetID(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "正常系: ID指定",
			args: args{
				r: httptest.NewRequest(
					http.MethodGet,
					"http://dummy.com/get?id=1",
					nil,
				),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "異常系: ID指定なし",
			args: args{
				r: httptest.NewRequest(
					http.MethodGet,
					"http://dummy.com/get",
					nil,
				),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.GetID(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getID() = %v, want %v", got, tt.want)
			}
		})
	}
}

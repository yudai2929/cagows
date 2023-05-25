package memory

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sivchari/cagows/train/model"
)

func TestTodoMemory_List(t *testing.T) {
	tests := []struct {
		name    string
		hook    func(m *TodoMemory)
		want    []*model.Todo
		wantErr bool
	}{
		{
			name: "正常系: データ取得",
			hook: func(m *TodoMemory) {
				m.mem[1] = &model.Todo{
					ID:        1,
					Title:     "test1",
					Completed: false,
				}
				m.mem[2] = &model.Todo{
					ID:        2,
					Title:     "test2",
					Completed: false,
				}
			},
			want: []*model.Todo{
				{
					ID:        1,
					Title:     "test1",
					Completed: false,
				},
				{
					ID:        2,
					Title:     "test2",
					Completed: false,
				},
			},
			wantErr: false,
		},
		{
			name:    "正常系: データがない場合",
			hook:    func(m *TodoMemory) {},
			want:    []*model.Todo{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TodoMemory{
				mem: map[int]*model.Todo{},
			}
			tt.hook(m)
			got := m.List()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("TodoMemory.List() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestTodoMemory_Get(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		hook    func(m *TodoMemory)
		want    *model.Todo
		wantErr bool
	}{
		{
			name: "正常系: データ取得",
			id:   1,
			hook: func(m *TodoMemory) {
				m.mem[1] = &model.Todo{
					ID:        1,
					Title:     "test1",
					Completed: false,
				}
			},
			want: &model.Todo{
				ID:        1,
				Title:     "test1",
				Completed: false,
			},
			wantErr: false,
		},
		{
			name:    "正常系: データがない場合",
			id:      1,
			hook:    func(m *TodoMemory) {},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TodoMemory{
				mem: map[int]*model.Todo{},
			}
			tt.hook(m)
			got, err := m.Get(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoMemory.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("TodoMemory.Get() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestTodoMemory_Add(t *testing.T) {
	tests := []struct {
		name    string
		todo    *model.Todo
		want    *model.Todo
		wantErr bool
	}{
		// TODO: ここにTodoを追加できた場合のテストを書く
		{
			name: "正常系: データ追加",
			todo: &model.Todo{
				Title:     "test1",
				Completed: false,
			},
			want: &model.Todo{
				ID:        1,
				Title:     "test1",
				Completed: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TodoMemory{
				mem: map[int]*model.Todo{},
			}
			got := m.Add(tt.todo)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("TodoMemory.Add() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestTodoMemory_Complete(t *testing.T) {
	tests := []struct {
		name    string
		hook    func(m *TodoMemory)
		id      int
		wantErr bool
	}{
		{
			name: "正常系: データ更新",
			hook: func(m *TodoMemory) {
				m.mem[1] = &model.Todo{
					ID:        1,
					Title:     "test1",
					Completed: false,
				}
			},
			id:      1,
			wantErr: false,
		},
		{
			name:    "異常系: データがない場合",
			hook:    func(m *TodoMemory) {},
			id:      1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TodoMemory{
				mem: map[int]*model.Todo{},
			}
			tt.hook(m)
			if err := m.Complete(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("TodoMemory.Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

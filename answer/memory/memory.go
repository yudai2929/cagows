package memory

import (
	"errors"
	"sort"
	"sync"

	"github.com/sivchari/cagows/answer/model"
	"github.com/sivchari/cagows/answer/repository"
)

type TodoMemory struct {
	mem map[int]*model.Todo
	// sync.RWMutextは複数のgoroutineから同時にアクセスされることを想定した構造体で
	// 読み込みはロック関係なく読み込めるが、書き込みはロックをかけて行う必要がある
	sync.RWMutex
}

// ここでrepository.Repository interfaceをRepository structが実装しているかチェックしている
var _ repository.Repository = (*TodoMemory)(nil)

func New() repository.Repository {
	return &TodoMemory{
		mem: map[int]*model.Todo{},
	}
}

func (m *TodoMemory) List() []*model.Todo {
	m.RLock()
	defer m.RUnlock()

	todos := make([]*model.Todo, 0, len(m.mem))
	for _, todo := range m.mem {
		todos = append(todos, todo)
	}
	// mapをforで回すと順番がランダムになるのでID順にソートする
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})
	return todos
}

func (m *TodoMemory) Get(id int) (*model.Todo, error) {
	m.RLock()
	defer m.RUnlock()

	todo, ok := m.mem[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return todo, nil
}

func (m *TodoMemory) Add(todo *model.Todo) *model.Todo {
	m.Lock()
	defer m.Unlock()

	id := len(m.mem) + 1
	todo.ID = id
	m.mem[id] = todo
	return todo
}

func (m *TodoMemory) Complete(id int) error {
	m.Lock()
	defer m.Unlock()

	todo, ok := m.mem[id]
	if !ok {
		return errors.New("not found")
	}
	todo.Completed = true
	return nil
}

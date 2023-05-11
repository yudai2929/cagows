package repository

import "github.com/sivchari/cagows/train/model"

type Repository interface {
	List() []*model.Todo
	Get(id int) (*model.Todo, error)
	Add(todo *model.Todo) *model.Todo
	Complete(id int) error
}

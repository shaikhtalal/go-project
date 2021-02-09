package internal

import (
	"context"
	"fmt"

	"bcdpkg.in/go-project/pkg/param"
	"bcdpkg.in/go-project/pkg/todo/graph/model"
	"github.com/pkg/errors"
)

type todo struct {
	ID   uint   `gorm:"column:id;primary_key"  json:"id"`
	Text string `gorm:"column:text;type:varchar(255)" json:"text"`
	Done bool   `gorm:"column:isDone" json:"done"`
}

func CreateTodo(ctx context.Context, input *model.NewTodo) (*model.Todo, error) {
	p, _ := param.EjectParamForResolver(ctx)
	t := todo{
		Text: input.Text,
	}
	found := p.Client.Db.HasTable(&todo{})
	if !found {
		p.Client.Db.AutoMigrate(&todo{})
	}
	if err := p.Client.Db.Create(&t).Error; err != nil {
		return nil, errors.Wrap(err, "failed to add into database")
	}
	return &model.Todo{
		ID:   fmt.Sprint(t.ID),
		Text: t.Text,
	}, nil
}

func GetAllTodo(ctx context.Context) ([]*model.Todo, error) {
	p, _ := param.EjectParamForResolver(ctx)
	t := []todo{}
	if err := p.Client.Db.Find(&t).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get from database")
	}
	todos := make([]*model.Todo, len(t))
	for i, todo := range t {
		todos[i] = &model.Todo{
			ID:   fmt.Sprint(todo.ID),
			Text: todo.Text,
			Done: todo.Done,
		}
	}
	// this is how we call the hasura Graphql
	// p.GraphQLClient.GetUsers(ctx)
	return todos, nil
}

//TableName return Table Name in DB
func (m *todo) TableName() string {
	return "todo"
}

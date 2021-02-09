package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"bcdpkg.in/go-project/pkg/internal"
	"bcdpkg.in/go-project/pkg/todo/graph/generated"
	"bcdpkg.in/go-project/pkg/todo/graph/model"
	"github.com/pkg/errors"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	res, err := internal.CreateTodo(ctx, &input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve and create the todo")
	}
	return res, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	res, err := internal.GetAllTodo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve and get the todo")
	}
	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

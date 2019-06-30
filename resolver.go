package graphqlblog

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateBlog(ctx context.Context, input *NewBlog) (*Blog, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateBlog(ctx context.Context, id *string, input *NewBlog) (*Blog, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteBlog(ctx context.Context, id *string) ([]*Blog, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Entries(ctx context.Context, search *string) ([]*Blog, error) {
	panic("not implemented")
}

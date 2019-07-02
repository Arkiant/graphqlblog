package graphqlblog

import (
	"context"
	"fmt"
	"io"

	"github.com/arkiant/graphqlblog/blogclient"

	"github.com/arkiant/grpc-go-course/blog/blogpb"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

// CreateBlog create a blog connected by grpc
func (r *mutationResolver) CreateBlog(ctx context.Context, input *NewBlog) (*Blog, error) {
	c, err := blogclient.Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer blogclient.Close()

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: input.AuthorID,
			Title:    input.Title,
			Content:  input.Content,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create blog: %v", err)
	}

	blogResult := res.GetBlog()
	return &Blog{
		ID:       blogResult.GetId(),
		AuthorID: blogResult.GetAuthorId(),
		Title:    blogResult.GetTitle(),
		Content:  blogResult.GetContent(),
	}, nil
}

func (r *mutationResolver) UpdateBlog(ctx context.Context, id *string, input *NewBlog) (*Blog, error) {
	c, err := blogclient.Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer blogclient.Close()

	if *id == "" || id == nil {
		return nil, fmt.Errorf("Cannot parse a null ID")
	}

	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       *id,
			AuthorId: input.AuthorID,
			Title:    input.Title,
			Content:  input.Content,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot update blog id %s: %v", *id, err)
	}

	blog := res.GetBlog()

	return &Blog{
		ID:       blog.GetId(),
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}, nil
}
func (r *mutationResolver) DeleteBlog(ctx context.Context, id *string) ([]*Blog, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

// Entries retrieve all entries from database
func (r *queryResolver) Entries(ctx context.Context, search *string) ([]*Blog, error) {

	c, err := blogclient.Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer blogclient.Close()

	result := make([]*Blog, 0)

	// if search is empty retrieve all collection
	if *search == "" {
		stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
		if err != nil {
			return nil, fmt.Errorf("Error while calling ListBlog RPC: %v", err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("Something happened: %v", err)
			}

			blog := res.GetBlog()

			result = append(result, &Blog{
				ID:       blog.GetId(),
				AuthorID: blog.GetAuthorId(),
				Title:    blog.GetTitle(),
				Content:  blog.GetContent(),
			})
		}
	} else {
		// if search has id retrieve a single collection
		res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: *search})
		if err != nil {
			return nil, fmt.Errorf("Something happened: %v", err)
		}

		blog := res.GetBlog()

		result = append(result, &Blog{
			ID:       blog.GetId(),
			AuthorID: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		})

	}

	return result, nil
}

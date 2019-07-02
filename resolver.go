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

func newBlogToPbBlog(id *string, input *NewBlog) *blogpb.Blog {
	blog := &blogpb.Blog{
		AuthorId: input.AuthorID,
		Title:    input.Title,
		Content:  input.Content,
	}

	if id != nil {
		blog.Id = *id
	}

	return blog
}

func pbBlogToBlog(input *blogpb.Blog) *Blog {
	return &Blog{
		ID:       input.GetId(),
		AuthorID: input.GetAuthorId(),
		Title:    input.GetTitle(),
		Content:  input.GetContent(),
	}
}

type mutationResolver struct{ *Resolver }

// CreateBlog create a blog connected by grpc
func (r *mutationResolver) CreateBlog(ctx context.Context, input *NewBlog) (*Blog, error) {
	res, err := blogclient.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: newBlogToPbBlog(nil, input),
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create blog: %v", err)
	}

	blogResult := res.GetBlog()
	return pbBlogToBlog(blogResult), nil
}

// UpdateBlog update a blog by id
func (r *mutationResolver) UpdateBlog(ctx context.Context, id *string, input *NewBlog) (*Blog, error) {
	if *id == "" || id == nil {
		return nil, fmt.Errorf("Cannot parse a null ID")
	}

	res, err := blogclient.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlogToPbBlog(id, input),
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot update blog id %s: %v", *id, err)
	}

	blog := res.GetBlog()

	return pbBlogToBlog(blog), nil
}

// DeleteBlog delete a blog by id
func (r *mutationResolver) DeleteBlog(ctx context.Context, id *string) ([]*Blog, error) {
	c, err := blogclient.Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer blogclient.Close()

	if *id == "" || id == nil {
		return nil, fmt.Errorf("Cannot parse a null ID")
	}

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: *id})
	if err != nil {
		return nil, fmt.Errorf("Something happened: %v", err)
	}

	blog := res.GetBlog()
	result := make([]*Blog, 0)
	result = append(result, pbBlogToBlog(blog))

	_, errDelete := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: *id})
	if errDelete != nil {
		return nil, fmt.Errorf("Cannot delete id %s, error: %v", *id, errDelete)
	}

	return result, nil
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

			result = append(result, pbBlogToBlog(blog))
		}
	} else {
		// if search has id retrieve a single collection
		res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: *search})
		if err != nil {
			return nil, fmt.Errorf("Something happened: %v", err)
		}

		blog := res.GetBlog()

		result = append(result, pbBlogToBlog(blog))

	}

	return result, nil
}

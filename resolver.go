package graphqlblog

import (
	"context"
	"io"
	"log"

	"github.com/arkiant/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
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

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(":50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	result := make([]*Blog, 0)

	// if search is empty retrieve all collection
	if *search == "" {
		stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
		if err != nil {
			log.Fatalf("Error while calling ListBlog RPC: %v\n", err)
		}

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Something happened: %v\n", err)
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
			log.Fatalf("Something happened: %v\n", err)
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

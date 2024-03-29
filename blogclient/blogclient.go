/*
 *
 * Copyright 2019.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//Package blogclient implements blog microservice client connection
package blogclient

import (
	"context"
	"fmt"
	"io"

	"github.com/arkiant/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

const (
	address = ":50051"
)

var (
	cc *grpc.ClientConn
)

// Connect function connect a singleton connection to client
// and returns a blog service client
func Connect() (blogpb.BlogServiceClient, error) {
	opts := grpc.WithInsecure()
	if cc == nil {
		cc, err := grpc.Dial(address, opts)
		if err != nil {
			return nil, fmt.Errorf("Could not connect: %v", err)
		}
		return blogpb.NewBlogServiceClient(cc), nil
	}

	return blogpb.NewBlogServiceClient(cc), nil
}

// Close the connection if exists
func Close() error {
	if cc != nil {
		return cc.Close()
	}
	return nil
}

// CreateBlog create a blog connected by grpc
func CreateBlog(ctx context.Context, blog *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	c, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer Close()

	return c.CreateBlog(context.Background(), blog)
}

// UpdateBlog update a blog
func UpdateBlog(ctx context.Context, blog *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	c, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer Close()

	return c.UpdateBlog(context.Background(), blog)
}

// DeleteBlog update a blog
func DeleteBlog(ctx context.Context, id *string) (*blogpb.Blog, error) {
	c, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer Close()

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: *id})
	if err != nil {
		return nil, fmt.Errorf("Something happened: %v", err)
	}

	blog := res.GetBlog()

	_, errDelete := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: *id})
	if errDelete != nil {
		return nil, fmt.Errorf("Cannot delete id %s, error: %v", *id, errDelete)
	}

	return blog, nil
}

// ListBlog retrieve all from database
func ListBlog(ctx context.Context, blog *blogpb.ListBlogRequest) ([]*blogpb.Blog, error) {
	c, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer Close()

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		return nil, fmt.Errorf("Error while calling ListBlog RPC: %v", err)
	}

	result := make([]*blogpb.Blog, 0)

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Something happened: %v", err)
		}

		blog := res.GetBlog()

		result = append(result, blog)

	}

	return result, nil
}

// ReadBlog retrieve a single collection
func ReadBlog(ctx context.Context, blog *blogpb.ReadBlogRequest) (*blogpb.Blog, error) {
	c, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("Could not connect: %v", err)
	}
	defer Close()

	res, err := c.ReadBlog(context.Background(), blog)
	if err != nil {
		return nil, fmt.Errorf("Something happened: %v", err)
	}

	return res.GetBlog(), nil
}

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

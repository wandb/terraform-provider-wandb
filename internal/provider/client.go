package provider

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/machinebox/graphql"
)

type GraphQLClientWithHeaders struct {
	client  *graphql.Client
	headers http.Header
}

func NewGraphQLClientWithHeaders(endpoint string, headers http.Header) *GraphQLClientWithHeaders {
	client := graphql.NewClient(endpoint)
	client.Log = func(s string) { log.Println(s) }
	return &GraphQLClientWithHeaders{
		client:  client,
		headers: headers,
	}
}

// Run wraps the graphql.Client's Run method to include headers.
func (c *GraphQLClientWithHeaders) Run(ctx context.Context, req *graphql.Request, resp interface{}) error {
	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Set(key, strings.TrimSpace(value))
		}
	}
	return c.client.Run(ctx, req, resp)
}

func (c *GraphQLClientWithHeaders) Log(s string) {
	c.client.Log(s)
}

package grpc_client

import (
	"context"
	"fmt"

	pb "article_service/grpc_go" // путь к сгенерированному grpc коду

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SearchClient struct {
	conn   *grpc.ClientConn
	client pb.SearchServiceClient
}

func NewSearchClient(addr string) (*SearchClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	client := pb.NewSearchServiceClient(conn)
	return &SearchClient{conn: conn, client: client}, nil
}

func (c *SearchClient) Close() error {
	return c.conn.Close()
}

func (c *SearchClient) IndexArticle(ctx context.Context, id int32, title, content string) error {
	req := &pb.ArticleEmbeddingRequest{
		Id:      id,
		Title:   title,
		Content: content,
	}
	_, err := c.client.IndexArticle(ctx, req)
	if err != nil {
		return fmt.Errorf("IndexArticle gRPC error: %w", err)
	}
	return nil
}

func (c *SearchClient) SemanticSearch(ctx context.Context, query string, limit int32) ([]int32, error) {
	req := &pb.SearchRequest{
		Query: query,
		Limit: limit,
	}
	resp, err := c.client.SemanticSearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("SemanticSearch gRPC error: %w", err)
	}
	return resp.ArticleIds, nil
}

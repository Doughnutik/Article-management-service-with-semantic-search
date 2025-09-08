package handler

import (
	"context"
	"fmt"

	"article_service/api"
	"article_service/internal/grpc_client"
	"article_service/internal/storage"
)

type articleHandler struct {
	db           *storage.Postgres
	searchClient *grpc_client.SearchClient
}

func NewArticleHandler(db *storage.Postgres, searchClient *grpc_client.SearchClient) api.Handler {
	return &articleHandler{
		db:           db,
		searchClient: searchClient,
	}
}

// GET /articles
func (h *articleHandler) ArticlesGet(ctx context.Context, params api.ArticlesGetParams) (api.ArticlesGetRes, error) {
	if v, ok := params.Page.Get(); ok && v < 1 {
		return &api.BadRequest{}, nil
	}
	if v, ok := params.Limit.Get(); ok && v < 1 {
		return &api.BadRequest{}, nil
	}

	res, err := h.db.ListArticles(ctx, params.Page.Or(1), params.Limit.Or(3))
	if err != nil {
		return &api.InternalServerError{}, nil
	}

	return res, nil
}

// GET /articles/{id}
func (h *articleHandler) ArticlesIDGet(ctx context.Context, params api.ArticlesIDGetParams) (api.ArticlesIDGetRes, error) {
	if params.ID < 0 {
		return &api.BadRequest{}, nil
	}
	res, err := h.db.GetArticle(ctx, params.ID)

	if err != nil {
		return &api.InternalServerError{}, nil
	}
	return res, nil
}

// DELETE /articles/{id}
func (h *articleHandler) ArticlesIDDelete(ctx context.Context, params api.ArticlesIDDeleteParams) (api.ArticlesIDDeleteRes, error) {
	if params.ID < 0 {
		return &api.BadRequest{}, nil
	}

	_, err := h.db.Pool.Exec(ctx, "DELETE FROM articles WHERE id=$1", params.ID)
	if err != nil {
		return &api.InternalServerError{}, nil
	}

	return &api.ArticlesIDDeleteNoContent{}, nil
}

// PUT /articles/{id}
func (h *articleHandler) ArticlesIDPut(ctx context.Context, req *api.ArticleUpdate, params api.ArticlesIDPutParams) (api.ArticlesIDPutRes, error) {
	if params.ID < 0 {
		return &api.BadRequest{}, nil
	}
	if req == nil {
		return &api.BadRequest{}, nil
	}

	res, err := h.db.UpdateArticle(ctx, params.ID, req)

	if err != nil {
		return &api.InternalServerError{}, nil
	}
	return res, nil
}

// POST /articles
func (h *articleHandler) ArticlesPost(ctx context.Context, req *api.ArticleCreate) (api.ArticlesPostRes, error) {
	if req == nil {
		return &api.BadRequest{}, nil
	}
	if req.Title == "" || req.Content == "" || req.Author == "" {
		return &api.BadRequest{}, nil
	}

	res, err := h.db.CreateArticle(ctx, req)
	if err != nil {
		return &api.InternalServerError{}, nil
	}
	return res, nil
}

// POST /articles/search
func (h *articleHandler) ArticlesSearchPost(ctx context.Context, req *api.SearchRequest) (api.ArticlesSearchPostRes, error) {
	if req == nil || req.Query == "" {
		return &api.BadRequest{}, nil
	}

	if v, ok := req.Limit.Get(); ok && v < 1 {
		return &api.BadRequest{}, nil
	}

	//TODO тут надо брать из БД статьи по вернувшимся ID
	IDs, err := h.searchClient.SemanticSearch(ctx, req.Query, int32(v))
	if err != nil {
		// если gRPC не ответил или ошибка на сервере — возвращаем 500
		return &api.InternalServerError{}, nil
	}

	fmt.Println(IDs)
	res := api.ArticlesSearchPostOKApplicationJSON{}
	return &res, nil
}

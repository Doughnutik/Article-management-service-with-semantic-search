package handler

import (
	"context"
	"fmt"
	"time"

	"article_service/api"
	"article_service/internal/grpc_client"
)

type articleHandler struct {
	// Тут можно хранить зависимости: БД, логгер, клиент SearchService и т.д.
	searchClient *grpc_client.SearchClient
}

func NewArticleHandler(searchClient *grpc_client.SearchClient) api.Handler {
	return &articleHandler{searchClient: searchClient}
}

// GET /articles
func (h *articleHandler) ArticlesGet(ctx context.Context, params api.ArticlesGetParams) (api.ArticlesGetRes, error) {
	if v, ok := params.Page.Get(); ok && v < 1 {
		return &api.BadRequest{}, nil
	}
	if v, ok := params.Limit.Get(); ok && v < 1 {
		return &api.BadRequest{}, nil
	}

	//TODO логику запроса

	return &api.ArticlesGetOKApplicationJSON{}, nil
}

// GET /articles/{id}
func (h *articleHandler) ArticlesIDGet(ctx context.Context, params api.ArticlesIDGetParams) (api.ArticlesIDGetRes, error) {
	if params.ID < 0 {
		return &api.BadRequest{}, nil
	}

	//TODO логику запроса

	article := api.Article{
		ID:        api.NewOptInt(params.ID),
		Title:     api.NewOptString("DemoTitle"),
		Content:   api.NewOptString("DemoContent"),
		Author:    api.NewOptString("DemoAuthor"),
		UpdatedAt: api.NewOptDateTime(time.Now()),
		Tags:      []string{"demotag1", "demotag2"},
	}
	return &article, nil
}

// DELETE /articles/{id}
func (h *articleHandler) ArticlesIDDelete(ctx context.Context, params api.ArticlesIDDeleteParams) (api.ArticlesIDDeleteRes, error) {
	if params.ID < 0 {
		return &api.BadRequest{}, nil
	}

	//TODO логику запроса

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

	//TODO логику запроса

	article := api.Article{
		ID:        api.NewOptInt(params.ID),
		Title:     api.NewOptString("UpdatedTitle"),
		Content:   api.NewOptString("UpdatedContent"),
		Author:    api.NewOptString("UpdatedAuthor"),
		UpdatedAt: api.NewOptDateTime(time.Now()),
		Tags:      []string{"updatedtag1", "updatedtag2"},
	}
	return &article, nil
}

// POST /articles
func (h *articleHandler) ArticlesPost(ctx context.Context, req *api.ArticleCreate) (api.ArticlesPostRes, error) {
	if req == nil {
		return &api.BadRequest{}, nil
	}
	if req.Title == "" || req.Content == "" || req.Author == "" {
		return &api.BadRequest{}, nil
	}

	//TODO логику запроса

	article := api.Article{
		ID:        api.NewOptInt(0),
		Title:     api.NewOptString(req.Title),
		Content:   api.NewOptString(req.Content),
		Author:    api.NewOptString(req.Author),
		UpdatedAt: api.NewOptDateTime(time.Now()),
		Tags:      req.Tags,
	}

	if err := h.searchClient.IndexArticle(ctx, int32(article.ID.Value), req.Title, req.Content); err != nil {
		return &api.InternalServerError{}, nil
	}

	return &article, nil
}

// POST /articles/search
func (h *articleHandler) ArticlesSearchPost(ctx context.Context, req *api.SearchRequest) (api.ArticlesSearchPostRes, error) {
	if req == nil || req.Query == "" {
		return &api.BadRequest{}, nil
	}

	if v, ok := req.Limit.Get(); ok && v < 1 {
		return &api.BadRequest{}, nil
	}

	//TODO логику запроса

	v, ok := req.Limit.Get()
	if !ok {
		v = 3
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

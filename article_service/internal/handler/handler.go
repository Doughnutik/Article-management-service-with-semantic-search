package handler

import (
	"context"
	"time"

	"article_service/api"
)

type articleHandler struct {
	// Тут можно хранить зависимости: БД, логгер, клиент SearchService и т.д.
}

func NewArticleHandler() api.Handler {
	return &articleHandler{}
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

	return &api.ArticlesSearchPostOKApplicationJSON{}, nil
}

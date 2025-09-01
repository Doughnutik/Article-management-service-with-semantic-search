package handler

import (
	"context"
	"errors"
	"time"

	"article_service/api"
	"article_service/internal/store"
)

// HandlerImpl реализует интерфейс api.Handler
type HandlerImpl struct {
	db *store.Store
}

func NewHandler(db *store.Store) *HandlerImpl {
	return &HandlerImpl{db: db}
}

// ArticlesPost — создание статьи
func (h *HandlerImpl) ArticlesPost(ctx context.Context, req *api.ArticleCreate) (api.ArticlesPostRes, error) {
	if req.Title == "" || req.Content == "" || req.Author == "" {
		return api.ArticlesPostRes{}, errors.New("bad request")
	}

	article := store.Article{
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		Tags:      req.Tags,
		UpdatedAt: time.Now(),
	}

	id, err := h.db.CreateArticle(ctx, &article)
	if err != nil {
		return api.ArticlesPostRes{}, err
	}

	res := api.Article{
		Id:        id,
		Title:     article.Title,
		Content:   article.Content,
		Author:    article.Author,
		UpdatedAt: article.UpdatedAt,
		Tags:      article.Tags,
	}

	return api.ArticlesPostRes{Article: &res}, nil
}

// ArticlesGet — получение списка статей
func (h *HandlerImpl) ArticlesGet(ctx context.Context, params api.ArticlesGetParams) (api.ArticlesGetRes, error) {
	articles, err := h.db.ListArticles(ctx, params.Page, params.Limit, params.Tags)
	if err != nil {
		return api.ArticlesGetRes{}, err
	}

	var res []api.Article
	for _, a := range articles {
		res = append(res, api.Article{
			Id:        a.ID,
			Title:     a.Title,
			Content:   a.Content,
			Author:    a.Author,
			UpdatedAt: a.UpdatedAt,
			Tags:      a.Tags,
		})
	}

	return api.ArticlesGetRes{Articles: res}, nil
}

// Остальные методы реализуются аналогично: ArticlesIDGet, ArticlesIDPut, ArticlesIDDelete

// ArticlesSearchPost — поиск по смыслу
func (h *HandlerImpl) ArticlesSearchPost(ctx context.Context, req *api.SearchRequest) (api.ArticlesSearchPostRes, error) {
	// здесь будем асинхронно вызывать SearchService по gRPC
	return api.ArticlesSearchPostRes{}, errors.New("not implemented")
}

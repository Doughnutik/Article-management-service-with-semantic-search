package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"article_service/api"
)

func (p *Postgres) CreateArticle(ctx context.Context, create *api.ArticleCreate) (*api.Article, error) {
	var id int
	cur_time := time.Now()
	err := p.Pool.QueryRow(
		ctx,
		`INSERT INTO articles (title, content, author, updated_at, tags)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id`,
		create.Title, create.Content, create.Author, cur_time, create.Tags,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	article := api.Article{
		ID:        api.NewOptInt(id),
		Title:     api.NewOptString(create.Title),
		Content:   api.NewOptString(create.Content),
		Author:    api.NewOptString(create.Author),
		UpdatedAt: api.NewOptDateTime(cur_time),
		Tags:      create.Tags,
	}
	return &article, nil
}

func (p *Postgres) GetArticle(ctx context.Context, id int) (*api.Article, error) {
	row := p.Pool.QueryRow(
		ctx,
		`SELECT id, title, content, author, updated_at, tags FROM articles WHERE id = $1`,
		id,
	)
	var a api.Article
	err := row.Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.UpdatedAt, &a.Tags)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (p *Postgres) UpdateArticle(ctx context.Context, id int, upd *api.ArticleUpdate) (*api.Article, error) {
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if v, ok := upd.Title.Get(); ok {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, v)
		argID++
	}
	if v, ok := upd.Content.Get(); ok {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argID))
		args = append(args, v)
		argID++
	}
	if v, ok := upd.Author.Get(); ok {
		setParts = append(setParts, fmt.Sprintf("author = $%d", argID))
		args = append(args, v)
		argID++
	}
	if upd.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argID))
		args = append(args, upd.Tags)
		argID++
	}

	setParts = append(setParts, "updated_at = now()")

	if len(setParts) == 1 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`UPDATE articles SET %s WHERE id = $%d`, strings.Join(setParts, ", "), argID)
	args = append(args, id)

	res, err := p.Pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("article with id %d not found", id)
	}

	article, err := p.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	return article, nil

}

func (p *Postgres) DeleteArticle(ctx context.Context, id int) error {
	_, err := p.Pool.Exec(ctx, `DELETE FROM articles WHERE id = $1`, id)
	return err
}

func (p *Postgres) ListArticles(ctx context.Context, limit int) ([]api.Article, error) {
	rows, err := p.Pool.Query(
		ctx,
		`SELECT id, title, content, author, updated_at, tags
		 FROM articles
		 ORDER BY updated_at DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []api.Article
	for rows.Next() {
		var a api.Article
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Author, &a.UpdatedAt, &a.Tags)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}

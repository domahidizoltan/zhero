package schemametadata

import (
	"context"

	"github.com/domahidizoltan/zhero/pkg/database"
)

type repo interface {
	Save(context.Context, Schema) error
}

type Service struct {
	repo repo
}

func NewService(repo repo) Service {
	return Service{
		repo: repo,
	}
}

func (m Service) Save(ctx context.Context, schema Schema) error {
	return database.InTx(ctx, func(ctx context.Context) error {
		return m.repo.Save(ctx, schema)
	})
}

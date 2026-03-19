package page

import (
	"context"

	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/domahidizoltan/zhero/pkg/paging"
)

type (
	pageRepo interface {
		Insert(context.Context, Page, string) (string, error)
		Update(context.Context, string, Page, string) error
		GetPageBySchemaNameAndIdentifier(context.Context, string, string, bool) (*Page, error)
		List(context.Context, string, ListOptions, bool) ([]Page, paging.Meta, error)
		Enable(context.Context, string, string, bool) error
		Delete(context.Context, string, string) error
		GetEnabledSchemaNames(context.Context) ([]string, error)
	}
	routeSvc interface {
		AssignRoute(ctx context.Context, customRoute, pageKey string) error
	}
)

type Service struct {
	pageRepo pageRepo
	routeSvc routeSvc
}

func NewService(repo pageRepo, routeSvc routeSvc) Service {
	return Service{
		pageRepo: repo,
		routeSvc: routeSvc,
	}
}

func (s Service) Create(ctx context.Context, page Page, idField string) (string, error) {
	createdID := ""
	if err := database.InTx(ctx, func(ctx context.Context) error {
		var err error
		if createdID, err = s.pageRepo.Insert(ctx, page, idField); err != nil {
			return err
		}

		if page.Route != "" {
			pageKey := page.SchemaName + "/" + createdID
			if err := s.routeSvc.AssignRoute(ctx, page.Route, pageKey); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return "", err
	}

	return createdID, nil
}

func (s Service) Update(ctx context.Context, identifier string, page Page, idField string) error {
	if err := database.InTx(ctx, func(ctx context.Context) error {
		if err := s.pageRepo.Update(ctx, identifier, page, idField); err != nil {
			return err
		}

		if page.Route != "" {
			pageKey := page.SchemaName + "/" + identifier
			if err := s.routeSvc.AssignRoute(ctx, page.Route, pageKey); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s Service) GetPageBySchemaNameAndIdentifier(ctx context.Context, schemaName, identifier string, onlyEnabled bool) (*Page, error) {
	return s.pageRepo.GetPageBySchemaNameAndIdentifier(ctx, schemaName, identifier, onlyEnabled)
}

func (s Service) List(ctx context.Context, schemaName string, opts ListOptions, onlyEnabled bool) ([]Page, paging.Meta, error) {
	return s.pageRepo.List(ctx, schemaName, opts, onlyEnabled)
}

func (s Service) Enable(ctx context.Context, schemaName, identifier string, enable bool) error {
	return database.InTx(ctx, func(ctx context.Context) error {
		return s.pageRepo.Enable(ctx, schemaName, identifier, enable)
	})
}

func (s Service) Delete(ctx context.Context, schemaName, identifier string) error {
	return database.InTx(ctx, func(ctx context.Context) error {
		return s.pageRepo.Delete(ctx, schemaName, identifier)
	})
}

func (s Service) GetEnabledSchemaNames(ctx context.Context) ([]string, error) {
	return s.pageRepo.GetEnabledSchemaNames(ctx)
}

package route

import (
	"context"
	"errors"
	"strings"

	"github.com/domahidizoltan/zhero/pkg/database"
	"github.com/domahidizoltan/zhero/pkg/url"
)

type (
	repo interface {
		Create(ctx context.Context, route, page string) error
		GetByRoute(ctx context.Context, route string) (*Route, error)
		GetLatestVersion(ctx context.Context, page string) (*Route, error)
	}
	Service struct {
		repo repo
	}
)

func NewService(repo repo) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) AssignRoute(ctx context.Context, customRoute, pageKey string) error {
	slug, err := s.slugifyRoute(customRoute)
	if err != nil {
		return err
	}

	return database.InTx(ctx, func(ctx context.Context) error {
		return s.repo.Create(ctx, slug, pageKey)
	})
}

func (s Service) GetValidSlug(ctx context.Context, customRoute string) (string, error) {
	slug, err := s.slugifyRoute(customRoute)
	if err != nil {
		return "", err
	}

	route, err := s.repo.GetByRoute(ctx, slug)
	if err != nil {
		return "", err
	}

	if route != nil {
		return "", errors.New("route already exists")
	}
	return slug, err
}

func (s Service) slugifyRoute(customRoute string) (string, error) {
	customRoute = strings.Trim(customRoute, "/")
	if customRoute == "" {
		return "", errors.New("route cannot be empty")
	}

	if len(customRoute) > 255 {
		return "", errors.New("route is too long (max 255 characters)")
	}

	parts := make([]string, 0, strings.Count(customRoute, "/")+1)
	for _, r := range strings.Split(customRoute, "/") {
		parts = append(parts, url.Slugify(r))
	}
	return "/" + strings.Join(parts, "/"), nil
}

func (s Service) GetByRoute(ctx context.Context, route string) (*Route, error) {
	return s.repo.GetByRoute(ctx, route)
}

func (s Service) GetLatestVersion(ctx context.Context, pageKey string) (*Route, error) {
	return s.repo.GetLatestVersion(ctx, pageKey)
}

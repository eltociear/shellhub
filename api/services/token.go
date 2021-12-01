package services

import (
	"context"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type TokenService interface {
	ListToken(ctx context.Context, tenantID string) ([]models.Token, error)
	CreateToken(ctx context.Context, tenantID string) (*models.Token, error)
	GetToken(ctx context.Context, tenantID string, ID string) (*models.Token, error)
	DeleteToken(ctx context.Context, tenantID string, ID string) error
	UpdateToken(ctx context.Context, tenantID string, ID string, token *models.APITokenUpdate) error
}

func (s *service) ListToken(ctx context.Context, tenantID string) ([]models.Token, error) {
	_, err := s.store.NamespaceGet(ctx, tenantID)
	if err != nil {
		return nil, ErrNamespaceNotFound
	}

	return s.store.TokenListAPIToken(ctx, tenantID)
}

func (s *service) CreateToken(ctx context.Context, tenantID string) (*models.Token, error) {
	_, err := s.store.NamespaceGet(ctx, tenantID)
	if err != nil {
		return nil, ErrNamespaceNotFound
	}

	return s.store.TokenCreateAPIToken(ctx, tenantID)
}

func (s *service) GetToken(ctx context.Context, tenantID string, id string) (*models.Token, error) {
	token, err := s.store.TokenGetAPIToken(ctx, tenantID, id)
	if err != nil {
		if err == store.ErrNoDocuments {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return token, nil
}

func (s *service) DeleteToken(ctx context.Context, tenantID string, id string) error {
	_, err := s.GetToken(ctx, tenantID, id)
	if err != nil {
		return ErrNotFound
	}

	return s.store.TokenDeleteAPIToken(ctx, tenantID, id)
}

func (s *service) UpdateToken(ctx context.Context, tenantID string, id string, request *models.APITokenUpdate) error {
	_, err := s.GetToken(ctx, tenantID, id)
	if err != nil {
		return ErrNotFound
	}

	return s.store.TokenUpdateAPIToken(ctx, tenantID, id, request)
}
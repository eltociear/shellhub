package store

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type PublicKeyStore interface {
	PublicKeyList(ctx context.Context, pagination paginator.Query) ([]models.PublicKey, int, error)
	PublicKeyGet(ctx context.Context, fingerprint string, tenantID string) (*models.PublicKey, error)
	PublicKeyCreate(ctx context.Context, key *models.PublicKey) error
	PublicKeyUpdate(ctx context.Context, fingerprint string, tenantID string, key *models.PublicKeyUpdate) (*models.PublicKey, error)
	PublicKeyDelete(ctx context.Context, fingerprint string, tenantID string) error
}

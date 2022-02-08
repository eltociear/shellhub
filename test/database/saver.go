package database

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/models"
)

type Saver interface {
	InsertUsers(ctx context.Context, users []models.User) error
	InsertNamespaces(ctx context.Context, namespaces []models.Namespace) error
}

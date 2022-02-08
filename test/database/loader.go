package database

import "github.com/shellhub-io/shellhub/pkg/models"

type Loader interface {
	LoadUsers(path string) []models.User
	LoadNamespaces(path string) []models.Namespace
}

package saver

import (
	"context"
	"fmt"

	"github.com/shellhub-io/shellhub/pkg/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Saver struct {
	Connection *mongo.Client
	Database   string
}

func insert(ctx context.Context, conn *mongo.Client, database, collection string, data interface{}) error {
	_, err := conn.Database(database).Collection(collection, nil).InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("saver: %w", err)
	}

	return nil
}

func (s *Saver) InsertUsers(ctx context.Context, data []models.User) error {
	for _, i := range data {
		err := insert(ctx, s.Connection, s.Database, "users", i)
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Saver) InsertNamespaces(ctx context.Context, data []models.Namespace) error {
	for _, i := range data {
		err := insert(ctx, s.Connection, s.Database, "namespaces", i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Saver) Close(ctx context.Context) error {
	log.Println("Deleting database...")
	err := s.Connection.Database(s.Database).Drop(ctx)
	if err != nil {
		err = fmt.Errorf("saver: %w", err)
		return err
	}

	log.Println("Database deleted!")

	return nil
}

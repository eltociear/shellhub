package database

import (
	"context"
	"fmt"

	"github.com/shellhub-io/shellhub/test/database/json/loader"
	"github.com/shellhub-io/shellhub/test/database/mongo/saver"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoURI = "172.21.0.3"
const MongoDatabase = "test"

func Populate(ctx context.Context) error {
	l := loader.Loader{}

	users, err := l.LoadUsers("./users.json")
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Debugln(err)
		return err
	}
	namespaces, err := l.LoadNamespaces("./namespaces.json")
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Configuration files loaded!")

	log.Println("Trying to connect to database test")

	client := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", MongoURI))
	connection, err := mongo.Connect(ctx, client)
	defer connection.Disconnect(ctx)
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Pinging...")

	if err = connection.Ping(ctx, nil); err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Database connected!")

	s := saver.Saver{
		Connection: connection,
		Database:   MongoDatabase,
	}

	log.Println("Inserting data loaded to the test database...")

	err = s.InsertUsers(ctx, users)
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	err = s.InsertNamespaces(ctx, namespaces)
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Data inserted!")

	return nil
}

func Clean(ctx context.Context) error {
	log.Println("Trying to connect to database test")

	client := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", MongoURI))
	connection, err := mongo.Connect(ctx, client)
	defer connection.Disconnect(ctx)
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Pinging...")

	if err = connection.Ping(ctx, nil); err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Database connected!")

	log.Println("Cleaning the test database")

	err = connection.Database(MongoDatabase).Drop(ctx)
	if err != nil {
		err = fmt.Errorf("database: %w", err)
		log.Errorln(err)
		return err
	}

	log.Println("Database cleaned!")

	return nil
}

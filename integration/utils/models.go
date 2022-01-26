package utils

import (
	"github.com/shellhub-io/shellhub/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Name      string             `bson:"name"`
	Confirmed bool               `bson:"confirmed"`
	Password  string             `bson:"password"`
}

type Login struct {
	Name      string
	Token     string
	User      *User
	Namespace *models.Namespace
}

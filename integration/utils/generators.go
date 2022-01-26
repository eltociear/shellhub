package utils

import (
	"fmt"
    "context"
    "log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/shellhub-io/shellhub/pkg/validator"
    "github.com/shellhub-io/shellhub/pkg/models"
)

func NewUser(db *mongo.Database, num int) *User {
    user := &User{
		ID: func() primitive.ObjectID {
			return primitive.NewObjectID()
		}(),
		Email:     fmt.Sprintf("user%d.email.com", num),
		Name:      fmt.Sprintf("user%d", num),
		Username:  fmt.Sprintf("username%d", num),
		Confirmed: true,
		Password:  validator.HashPassword("senha"),
	}

	_, _ = db.Collection("users", nil).InsertOne(context.TODO(), user)
    return user
}

func NewNamespace(db *mongo.Database, num int, owner *User, members []User, roles []string) *models.Namespace {
    ns := &models.Namespace{
		Name:     fmt.Sprintf("namespace%d", num),
		TenantID: fmt.Sprintf("xxx%d", num),
		Members:  append([]models.Member{{ID: owner.ID.Hex(), Role: "owner", Username: owner.Username}}, func (members []User, roles []string) []models.Member {
            if len(members) < 1 || len(members) != len(roles) {
                return []models.Member{}
            }

            vRole := map[string]bool{
                "administrator": true,
                "observer": true,
                "operator": true,
            }

            var users []models.Member
            for i, member := range members {
                _, ok := vRole[roles[i]]
                if !ok {
                    return []models.Member{}
                }

                users = append(users, models.Member{
                    ID: member.ID.Hex(),
                    Role: roles[i],
                    Username: member.Username,
                })
            }

            return users
        }(members, roles)...),
    }

    _, _ = db.Collection("namespaces", nil).InsertOne(context.TODO(), ns)
    return ns
}

func NewLogin(db *mongo.Database, ns *models.Namespace, user *User) (*Login, error) {
	token, err := TokenFromAuthentication(user.Username, "senha")
	if err != nil {
		log.Println(token, err)
		return nil, err
	}

	token2, err := TenantSwitch(token, fmt.Sprintf(ns.TenantID))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	name := fmt.Sprintf("%s-%s", user.Name, ns.Name)

	return &Login{
		Name:      name,
		User:      user,
		Namespace: ns,
		Token:     token2,
	}, nil
}

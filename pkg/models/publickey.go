package models

import (
	"errors"
	"reflect"
	"time"
)

const (
	// PublicKeyKindHostname is a PublicKeyKind.Key of PublicKeyKind what represent a device's hostname.
	PublicKeyKindHostname = "hostname"
	// PublicKeyKindTag is a PublicKeyKind.Key of PublicKeyKind what represent device's tags.
	PublicKeyKindTag = "tag"
)

// PublicKeyKind contains the category of a Public Key.
//
// A PublicKeyKind can be eiter PublicKeyKindHostname or PublicKeyKindTag. The first type, PublicKeyKindHostname, accepts
// in the Value field a string, a device's Hostname and the second, PublicKeyKindTag, a slice of strings what contains
// the tags required for a device to be accessible by that PublicKey.
type PublicKeyKind struct {
	Key   string      `json:"key" bson:"key" validate:"oneof=hostname tag"`
	Value interface{} `json:"value" bson:"value" validate:"required"`
}

// GetHostname returns, if the PublicKeyKind.Key is PublicKeyKindHostname, the hostname.
func (k *PublicKeyKind) GetHostname() (string, error) {
	if k.Key != PublicKeyKindHostname {
		return "", errors.New("the PublicKey kind is is not valid for this method")
	}

	value := reflect.ValueOf(k.Value)
	if value.Kind() != reflect.String {
		return "", errors.New("value is not hostname")
	}

	hostname := value.Interface().(string)

	//v := validator.New()
	//_ = v.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
	//	_, err := regexp.Compile(fl.Field().String())
	//
	//	return err == nil
	//})

	return hostname, nil
}

// GetTags returns, if the PublicKeyKind.Key is PublicKeyKindTag, the first 3 tags.
func (k *PublicKeyKind) GetTags() ([]string, error) {
	if k.Key != PublicKeyKindTag {
		return nil, errors.New("the PublicKey kind is is not valid for this method")
	}

	value := reflect.ValueOf(k.Value)
	if value.Kind() != reflect.Slice {
		return nil, errors.New("value is not tags")
	}

	tags := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		tags[i] = value.Index(i).Interface().(string)
	}

	return tags, nil
}

type PublicKeyFields struct {
	Name     string        `json:"name"`
	Username string        `json:"username" bson:"username,omitempty"`
	Kind     PublicKeyKind `json:"kind" bson:"kind" validate:"required"`
}

type PublicKey struct {
	Data            []byte    `json:"data"`
	Fingerprint     string    `json:"fingerprint"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	TenantID        string    `json:"tenant_id" bson:"tenant_id"`
	PublicKeyFields `bson:",inline"`
}

type PublicKeyUpdate struct {
	PublicKeyFields `bson:",inline"`
}

type PublicKeyAuthRequest struct {
	Fingerprint string `json:"fingerprint"`
	Data        string `json:"data"`
}

type PublicKeyAuthResponse struct {
	Signature string `json:"signature"`
}

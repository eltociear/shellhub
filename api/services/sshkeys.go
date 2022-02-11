package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/clock"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/validator"
	"golang.org/x/crypto/ssh"
)

type SSHKeysService interface {
	EvaluateKey(ctx context.Context, key *models.PublicKey, device *models.Device, name string) (bool, error)
	ListPublicKeys(ctx context.Context, pagination paginator.Query) ([]models.PublicKey, int, error)
	GetPublicKey(ctx context.Context, fingerprint, tenant string) (*models.PublicKey, error)
	CreatePublicKey(ctx context.Context, key *models.PublicKey, tenant string) error
	UpdatePublicKey(ctx context.Context, fingerprint, tenant string, key *models.PublicKeyUpdate) (*models.PublicKey, error)
	DeletePublicKey(ctx context.Context, fingerprint, tenant string) error
	CreatePrivateKey(ctx context.Context) (*models.PrivateKey, error)
}

type Request struct {
	Namespace string
}

// EvaluateKey evaluate if a public key
func (s *service) EvaluateKey(ctx context.Context, key *models.PublicKey, device *models.Device, name string) (bool, error) {
	_, err := validator.ValidateStruct(key)
	if err != nil {
		return false, ErrInvalidFormat
	}

	if err != nil {
		return false, err
	}

	switch key.Kind.Key {
	case models.PublicKeyKindHostname:
		hostname, err := key.Kind.GetHostname()
		if err != nil {
			return false, err
		}

		x, err := regexp.MatchString(hostname, device.Name)
		if err != nil {
			return false, err
		}

		y, err := regexp.MatchString(key.Username, name)
		if err != nil {
			return false, err
		}

		return x && y, nil
	case models.PublicKeyKindTag:
		tags, err := key.Kind.GetTags()
		if err != nil {
			return false, err
		}

		for _, tag := range tags {
			if !contains(device.Tags, tag) {
				return false, errors.New("a tag is not valid for this Public Key")
			}
		}

		return true, nil
	}

	return false, errors.New("this kind of public key kind is not valid")
}

func (s *service) GetPublicKey(ctx context.Context, fingerprint, tenant string) (*models.PublicKey, error) {
	return s.store.PublicKeyGet(ctx, fingerprint, tenant)
}

func (s *service) CreatePublicKey(ctx context.Context, key *models.PublicKey, tenant string) error {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(key.Data) //nolint:dogsled
	if err != nil {
		return ErrInvalidFormat
	}

	_, err = validator.ValidateStruct(key)
	if err != nil {
		return ErrInvalidFormat
	}

	key.CreatedAt = clock.Now()
	key.Fingerprint = ssh.FingerprintLegacyMD5(pub)

	switch key.Kind.Key {
	case models.PublicKeyKindHostname:
		//hostname, err := key.Kind.GetHostname()
		//if err != nil {
		//	return err
		//}
		//
		//_, err = validator.ValidateVar(hostname, "regexp") // TODO: remove this static validation from here.
		//if err != nil {
		//	return err
		//}

		break
	case models.PublicKeyKindTag:
		tags, err := key.Kind.GetTags()
		if err != nil {
			return err
		}

		for _, tag := range tags {
			if !validator.ValidateFieldTag(tag) {
				return fmt.Errorf("%s is not a valid tag", tag)
			}
		}

		break
	}

	// Check if the public key already exists.
	returned, err := s.store.PublicKeyGet(ctx, key.Fingerprint, tenant)
	if err != nil && err != store.ErrNoDocuments {
		return err
	}

	if returned != nil {
		return ErrDuplicateFingerprint
	}

	err = s.store.PublicKeyCreate(ctx, key)
	if err != nil {
		return err
	}

	return err
}

func (s *service) ListPublicKeys(ctx context.Context, pagination paginator.Query) ([]models.PublicKey, int, error) {
	return s.store.PublicKeyList(ctx, pagination)
}

func (s *service) UpdatePublicKey(ctx context.Context, fingerprint, tenant string, key *models.PublicKeyUpdate) (*models.PublicKey, error) {
	return s.store.PublicKeyUpdate(ctx, fingerprint, tenant, key)
}

func (s *service) DeletePublicKey(ctx context.Context, fingerprint, tenant string) error {
	return s.store.PublicKeyDelete(ctx, fingerprint, tenant)
}

func (s *service) CreatePrivateKey(ctx context.Context) (*models.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	pubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	privateKey := &models.PrivateKey{
		Data: pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		}),
		Fingerprint: ssh.FingerprintLegacyMD5(pubKey),
		CreatedAt:   clock.Now(),
	}

	if err := s.store.PrivateKeyCreate(ctx, privateKey); err != nil {
		return nil, err
	}

	return privateKey, nil
}

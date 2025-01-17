package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"regexp"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/clock"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/shellhub-io/shellhub/pkg/validator"
	"golang.org/x/crypto/ssh"
)

type SSHKeysService interface {
	EvaluateKeyFilter(ctx context.Context, key *models.PublicKey, dev models.Device) (bool, error)
	EvaluateKeyUsername(ctx context.Context, key *models.PublicKey, username string) (bool, error)
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

func (s *service) EvaluateKeyFilter(ctx context.Context, key *models.PublicKey, dev models.Device) (bool, error) {
	if key.Filter.Hostname != "" {
		ok, err := regexp.MatchString(key.Filter.Hostname, dev.Name)
		if err != nil {
			return false, err
		}

		return ok, nil
	} else if len(key.Filter.Tags) > 0 {
		for _, tag := range dev.Tags {
			if contains(key.Filter.Tags, tag) {
				return true, nil
			}
		}

		return false, nil
	}

	return true, nil
}

func (s *service) EvaluateKeyUsername(ctx context.Context, key *models.PublicKey, username string) (bool, error) {
	if key.Username == "" {
		return true, nil
	}

	ok, err := regexp.MatchString(key.Username, username)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *service) GetPublicKey(ctx context.Context, fingerprint, tenant string) (*models.PublicKey, error) {
	_, err := s.store.NamespaceGet(ctx, tenant)
	if err != nil {
		return nil, NewErrNamespaceNotFound(tenant, err)
	}

	return s.store.PublicKeyGet(ctx, fingerprint, tenant)
}

func (s *service) CreatePublicKey(ctx context.Context, key *models.PublicKey, tenant string) error {
	if err := key.Validate(); err != nil {
		data, _ := validator.GetInvalidFieldsValues(err)

		return NewErrPublicKeyInvalid(data, nil)
	}

	// Checks if public key filter type is Tags.
	// If it is, checks if there are, at least, one tag on the public key filter and if the all tags exist on database.
	if key.Filter.Tags != nil {
		if len(key.Filter.Tags) == 0 {
			return NewErrPublicKeyTagsEmpty(nil)
		}

		tags, _, err := s.store.TagsGet(ctx, tenant)
		if err != nil {
			return NewErrTagEmpty(tenant, err)
		}

		for _, tag := range key.Filter.Tags {
			if !contains(tags, tag) {
				return NewErrTagNotFound(tag, nil)
			}
		}
	}

	key.CreatedAt = clock.Now()

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(key.Data) //nolint:dogsled
	if err != nil {
		return NewErrPublicKeyDataInvalid(key.Data, nil)
	}

	key.Fingerprint = ssh.FingerprintLegacyMD5(pubKey)

	returnedKey, err := s.store.PublicKeyGet(ctx, key.Fingerprint, tenant)
	if err != nil && err != store.ErrNoDocuments {
		return NewErrPublicKeyNotFound(key.Fingerprint, err)
	}

	if returnedKey != nil {
		return NewErrPublicKeyDuplicated([]string{key.Fingerprint}, err)
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
	if err := key.Validate(); err != nil {
		data, _ := validator.GetInvalidFieldsValues(err)

		return nil, NewErrPublicKeyInvalid(data, nil)
	}

	// Checks if public key filter type is Tags. If it is, checks if there are, at least, one tag on the public key
	// filter and if the all tags exist on database.
	if key.Filter.Tags != nil {
		if len(key.Filter.Tags) == 0 {
			return nil, NewErrPublicKeyTagsEmpty(nil)
		}

		tags, _, err := s.store.TagsGet(ctx, tenant)
		if err != nil {
			return nil, NewErrTagEmpty(tenant, err)
		}

		for _, tag := range key.Filter.Tags {
			if !contains(tags, tag) {
				return nil, NewErrTagNotFound(tag, nil)
			}
		}
	}

	return s.store.PublicKeyUpdate(ctx, fingerprint, tenant, key)
}

func (s *service) DeletePublicKey(ctx context.Context, fingerprint, tenant string) error {
	_, err := s.store.NamespaceGet(ctx, tenant)
	if err != nil {
		return NewErrNamespaceNotFound(tenant, err)
	}

	_, err = s.store.PublicKeyGet(ctx, fingerprint, tenant)
	if err != nil {
		return NewErrPublicKeyNotFound(fingerprint, err)
	}

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

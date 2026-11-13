package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type recordingClient struct {
	encryptInput *kms.EncryptInput
	decryptInput *kms.DecryptInput
}

func (c *recordingClient) Encrypt(_ context.Context, params *kms.EncryptInput, _ ...func(*kms.Options)) (*kms.EncryptOutput, error) {
	c.encryptInput = params

	return &kms.EncryptOutput{CiphertextBlob: params.Plaintext}, nil
}

func (c *recordingClient) Decrypt(_ context.Context, params *kms.DecryptInput, _ ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	c.decryptInput = params

	return &kms.DecryptOutput{Plaintext: params.CiphertextBlob}, nil
}

func TestRepository_NoContext_LeavesEncryptionContextNil(t *testing.T) {
	for name, encCtx := range map[string]map[string]string{
		"nil map":   nil,
		"empty map": {},
	} {
		t.Run(name, func(t *testing.T) {
			c := &recordingClient{}
			r := NewRepository(c, "key")

			if _, err := r.Encrypt(context.TODO(), "test", encCtx); err != nil {
				t.Fatalf("error encrypting: %v", err)
			}

			if _, err := r.Decrypt("dGVzdA==", encCtx); err != nil {
				t.Fatalf("error decrypting: %v", err)
			}

			if c.encryptInput.EncryptionContext != nil {
				t.Errorf("expected nil encrypt context, got %#v", c.encryptInput.EncryptionContext)
			}

			if c.decryptInput.EncryptionContext != nil {
				t.Errorf("expected nil decrypt context, got %#v", c.decryptInput.EncryptionContext)
			}
		})
	}
}

func TestRepository_WithContext_PassesItThrough(t *testing.T) {
	c := &recordingClient{}
	r := NewRepository(c, "key")
	encCtx := map[string]string{"purpose": "develop-freeze"}

	if _, err := r.Encrypt(context.TODO(), "test", encCtx); err != nil {
		t.Fatalf("error encrypting: %v", err)
	}

	if _, err := r.Decrypt("dGVzdA==", encCtx); err != nil {
		t.Fatalf("error decrypting: %v", err)
	}

	if !reflect.DeepEqual(c.encryptInput.EncryptionContext, encCtx) {
		t.Errorf("expected encrypt context %v, got %v", encCtx, c.encryptInput.EncryptionContext)
	}

	if !reflect.DeepEqual(c.decryptInput.EncryptionContext, encCtx) {
		t.Errorf("expected decrypt context %v, got %v", encCtx, c.decryptInput.EncryptionContext)
	}
}

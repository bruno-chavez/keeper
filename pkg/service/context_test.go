package service

import (
	"context"
	"reflect"
	"testing"
)

type recordingCrypto struct {
	encryptCtx map[string]string
	decryptCtx map[string]string
}

func (r *recordingCrypto) Encrypt(_ context.Context, data string, encCtx map[string]string) (string, error) {
	r.encryptCtx = encCtx

	return data, nil
}

func (r *recordingCrypto) Decrypt(data string, encCtx map[string]string) (string, error) {
	r.decryptCtx = encCtx

	return data, nil
}

func TestService_PassesEncryptionContextThrough(t *testing.T) {
	encCtx := map[string]string{"purpose": "develop-freeze"}
	rc := &recordingCrypto{}
	s := NewService(rc, rc)

	if _, err := s.Encrypt(context.TODO(), "a = b", encCtx); err != nil {
		t.Fatalf("error encrypting: %v", err)
	}

	if _, err := s.DecryptTerraform("a = b", encCtx); err != nil {
		t.Fatalf("error decrypting: %v", err)
	}

	if !reflect.DeepEqual(rc.encryptCtx, encCtx) || !reflect.DeepEqual(rc.decryptCtx, encCtx) {
		t.Errorf("context not passed through: encrypt=%v decrypt=%v", rc.encryptCtx, rc.decryptCtx)
	}
}

func TestService_NoContextStaysNil(t *testing.T) {
	rc := &recordingCrypto{}
	s := NewService(rc, rc)

	if _, err := s.Encrypt(context.TODO(), "a = b", nil); err != nil {
		t.Fatalf("error encrypting: %v", err)
	}

	if _, err := s.DecryptTerraform("a = b", nil); err != nil {
		t.Fatalf("error decrypting: %v", err)
	}

	if rc.encryptCtx != nil || rc.decryptCtx != nil {
		t.Errorf("expected nil contexts, got encrypt=%v decrypt=%v", rc.encryptCtx, rc.decryptCtx)
	}
}

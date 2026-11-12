package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestParseEncryptionContext_Valid(t *testing.T) {
	tests := map[string]map[string]string{
		"":                                    nil,
		"purpose:develop-freeze":              {"purpose": "develop-freeze"},
		"purpose:develop-freeze,team:infra":   {"purpose": "develop-freeze", "team": "infra"},
		" purpose : develop-freeze ":          {"purpose": "develop-freeze"},
		"arn:arn:aws:kms:us-east-1:1:key/abc": {"arn": "arn:aws:kms:us-east-1:1:key/abc"},
	}

	for input, expected := range tests {
		got, err := parseEncryptionContext(input)
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", input, err)
			continue
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("input %q: expected %v, got %v", input, expected, got)
		}
	}
}

func TestParseEncryptionContext_Invalid(t *testing.T) {
	for _, input := range []string{"purpose", "purpose:", ":value", "a:b,", "a:b,c", "a:b,a:c", ","} {
		got, err := parseEncryptionContext(input)
		if err == nil {
			t.Errorf("input %q: expected an error, got %v", input, got)
		}
	}
}

type recordingEncryptor struct {
	called bool
	encCtx map[string]string
}

func (r *recordingEncryptor) Encrypt(_ context.Context, _ string, encCtx map[string]string) (string, error) {
	r.called = true
	r.encCtx = encCtx

	return "ok", nil
}

func postEncrypt(h *Handler, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", target, strings.NewReader(`a = "b"`))
	req.Header.Set("Content-Type", "text/plain")
	res := httptest.NewRecorder()
	h.Encrypt(res, req)

	return res
}

func TestHandler_Encrypt_WithoutContextParam(t *testing.T) {
	e := &recordingEncryptor{}
	res := postEncrypt(NewHandler(e), "/encrypt")

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 but got %d", res.Code)
	}

	if e.encCtx != nil {
		t.Errorf("expected nil context, got %v", e.encCtx)
	}
}

func TestHandler_Encrypt_WithContextParam(t *testing.T) {
	e := &recordingEncryptor{}
	res := postEncrypt(NewHandler(e), "/encrypt?context=purpose:develop-freeze,team:infra")

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 but got %d", res.Code)
	}

	expected := map[string]string{"purpose": "develop-freeze", "team": "infra"}
	if !reflect.DeepEqual(e.encCtx, expected) {
		t.Errorf("expected %v, got %v", expected, e.encCtx)
	}
}

func TestHandler_Encrypt_InvalidContextParam(t *testing.T) {
	e := &recordingEncryptor{}
	res := postEncrypt(NewHandler(e), "/encrypt?context=malformed")

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d", res.Code)
	}

	if res.Body.String() != invalidContextMessage {
		t.Errorf("expected invalid context message, got %q", res.Body.String())
	}

	if e.called {
		t.Error("encryptor must not be called on invalid context")
	}
}

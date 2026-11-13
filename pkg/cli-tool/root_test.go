package cli_tool

import (
	"reflect"
	"testing"
)

func TestRoot_ContextFlag(t *testing.T) {
	t.Cleanup(func() {
		_ = root.PersistentFlags().Set("context", "")
		encryptionContext = nil
	})

	if encryptionContext != nil {
		t.Fatalf("expected nil context by default, got %v", encryptionContext)
	}

	err := root.PersistentFlags().Parse([]string{"--context", "purpose=develop-freeze", "--context", "team=infra,env=prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{"purpose": "develop-freeze", "team": "infra", "env": "prod"}
	if !reflect.DeepEqual(encryptionContext, expected) {
		t.Errorf("expected %v, got %v", expected, encryptionContext)
	}
}

func TestRoot_ContextFlagInvalid(t *testing.T) {
	if err := root.PersistentFlags().Parse([]string{"--context", "novalue"}); err == nil {
		t.Error("expected an error for a pair without '='")
	}
}

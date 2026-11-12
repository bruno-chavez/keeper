package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (r Repository) Encrypt(ctx context.Context, data string, encCtx map[string]string) (string, error) {
	input := &kms.EncryptInput{
		KeyId: aws.String(r.encryptionKeyName),
	}

	input.Plaintext = []byte(data)

	// Only set when non-empty so requests without context stay identical to the previous behavior
	if len(encCtx) > 0 {
		input.EncryptionContext = encCtx
	}

	output, err := r.client.Encrypt(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error in repository encrypting string: %w", err)
	}

	return r.encoder.EncodeToString(output.CiphertextBlob), err
}

package util

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func UUID2Base64(id uuid.UUID) string {
	return base64.RawURLEncoding.EncodeToString(id)
}

package sign

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/LexusEgorov/goMetrics/internal/middleware"
)

type sign struct {
	key string
}

func (s sign) Sign(data []byte) string {
	if s.key == "" {
		return ""
	}

	h := hmac.New(sha256.New, []byte(s.key))
	h.Write(data)

	return string(h.Sum(nil))
}

func (s sign) Verify(data []byte, sign string) bool {
	if s.key == "" {
		return true
	}

	return s.Sign(data) == sign
}

func NewSign(key string) middleware.Signer {
	return sign{
		key: key,
	}
}

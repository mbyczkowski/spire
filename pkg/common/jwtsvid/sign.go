package jwtsvid

import (
	"crypto"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spiffe/spire/pkg/common/idutil"
)

const (
	keyIDHeader = "kid"
)

func SignToken(spiffeID string, audience []string, expires time.Time, signer crypto.Signer, kid string) (string, error) {
	if err := idutil.ValidateSpiffeID(spiffeID, idutil.AllowAnyTrustDomainWorkload()); err != nil {
		return "", err
	}

	audience = pruneEmptyValues(audience)

	if expires.IsZero() {
		return "", errors.New("expiration is required")
	}
	if len(audience) == 0 {
		return "", errors.New("audience is required")
	}
	if len(kid) == 0 {
		return "", errors.New("kid is required")
	}

	claims := jwt.MapClaims{
		"sub": spiffeID,
		"exp": expires.Unix(),
		"aud": audienceClaim(audience),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(signingMethodES256, claims)
	token.Header[keyIDHeader] = kid
	signedToken, err := token.SignedString(signer)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func pruneEmptyValues(values []string) []string {
	pruned := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			pruned = append(pruned, value)
		}
	}
	return pruned
}

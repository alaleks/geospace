// Package authentication provides authentication in app and performs checking auth data against.
package authentication

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alaleks/geospace/internal/server/config"
	"github.com/alaleks/geospace/internal/server/database"
	"github.com/golang-jwt/jwt"
	"github.com/golang-module/dongle"
)

const (
	Expiration = time.Second * 259200 // expiration time token validity period
)

// typical errors
var (
	ErrInvalidClaim = errors.New("invalid token claim")
)

// Auth contains db instance, cipher and secret key for JWT.
type Auth struct {
	db        *database.DB
	cipher    *dongle.Cipher
	secretJWT string
}

// Init performs initialization pointer of the Auth instance.
func Init(db *database.DB, cfgSecure config.Secure) *Auth {
	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.CBC)              // CBC、CFB、OFB、CTR、ECB
	cipher.SetPadding(dongle.PKCS7)         // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cipher.SetKey(cfgSecure.GetKeyCipher()) // key must from 1 to 56 bytes
	cipher.SetIV(cfgSecure.GetIVCipher())   // iv must be 8 bytes

	return &Auth{
		db:        db,
		cipher:    cipher,
		secretJWT: cfgSecure.GetSecretJWT(),
	}
}

// CheckPass performs check password from user with password from database
// and returns false if passwords do not match.
func (a *Auth) CheckPass(passFromUser, passFromDB string) bool {
	return strings.TrimSpace(passFromUser) ==
		dongle.Decrypt.FromHexString(passFromDB).ByBlowfish(a.cipher).ToString()
}

// EncryptPass performs encrypt password from user.
func (a *Auth) EncryptPass(pass string) string {
	return dongle.Encrypt.FromString(pass).ByBlowfish(a.cipher).ToHexString()
}

// GetTokenJWT performs generate token jwt for user.
func (a *Auth) GetTokenJWT(uid int) (string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["uid"] = uid
	claims["exp"] = now.Add(Expiration).Unix()

	tokenString, err := tokenByte.SignedString([]byte(a.secretJWT))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

// CheckToken perfoms validate jwt token.
func (a *Auth) CheckToken(token string) error {
	tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(a.secretJWT), nil
	})
	if err != nil {
		return err
	}

	_, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return ErrInvalidClaim
	}

	return nil
}

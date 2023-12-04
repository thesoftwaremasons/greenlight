package data

import (
	"context"
	"database/sql"
	"github.com/thesoftwaremasons/greenlight/internal/validator"
	"time"
)

// ValidateTokenPlaintext Check that the plaintext token has been provided and is exactly 26 bytes long.
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	return token, err
}
func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO Tokens(hash,user_id,expiry,scope)
			VALUES ($1,$2,$3,$4)
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{token.Hash, token.UserId, token.Expiry, token.Scope}
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m TokenModel) Delete(scope string, userId int64) error {

	query := `DELETE FROM Tokens where user_id=$1 AND scope=$2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, userId, scope)
	return err

}

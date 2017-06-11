package main

import (
	"github.com/asticode/go-astimail"
	"github.com/asticode/go-astitools/string"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Vars
var (
	errNotFound = errors.New("not.found")
	storage     Storage
)

// Base represents a base model
type Base struct {
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"udated_at"`
}

// Email represents an email
type Email struct {
	Base
	Addr        string         `db:"addr"`
	ID          int            `db:"id"`
	Token       string         `db:"token"`
	UserID      int            `db:"user_id"`
	ValidatedAt mysql.NullTime `db:"validated_at"`
}

// User represents a user
type User struct {
	Base
	ClientPublicKey  *astimail.PublicKey  `db:"client_public_key"`
	ID               int                  `db:"id"`
	ServerPrivateKey *astimail.PrivateKey `db:"server_private_key"`
	Username         string               `db:"username"`
}

// Storage represents a storage
type Storage interface {
	EmailCreate(email string, u *User) (token string, err error)
	EmailFetchWithValidationToken(token string) (e *Email, err error)
	EmailValidate(e *Email) (err error)
	UserCreate(cltPubKey *astimail.PublicKey, srvPrvKey *astimail.PrivateKey) error
	UserFetchWithEmail(email string) (*User, error)
	UserFetchWithKey(key *astimail.PublicKey) (*User, error)
	UserUpdate(u *User, cltPubKey *astimail.PublicKey, srvPrvKey *astimail.PrivateKey) error
}

// storageMySQL represents a MySQL storage
type storageMySQL struct {
	db *sqlx.DB
}

// newStorageMySQL builds a new mysql storage
func newStorageMySQL(db *sqlx.DB) *storageMySQL {
	return &storageMySQL{db: db}
}

// EmailCreate creates an email
func (s *storageMySQL) EmailCreate(email string, u *User) (token string, err error) {
	token = astistring.RandomString(100)
	_, err = s.db.Exec("INSERT INTO email (addr, user_id, validation_token) VALUES ($1, $2, $3) ON DUPLICATE KEY UPDATE validation_token = $3", email, u.ID, token)
	return
}

// EmailFetchWithValidationToken fetches an email based on a validation token
func (s *storageMySQL) EmailFetchWithValidationToken(token string) (e *Email, err error) {
	err = s.db.Get(e, "SELECT * FROM email WHERE validation_token = $1 AND validated_at IS NULL LIMIT 1", token)
	return
}

// EmailValidate validates an email
func (s *storageMySQL) EmailValidate(e *Email) (err error) {
	_, err = s.db.Exec("UPDATE email SET validated_at = NOW() WHERE id = $1", e.ID)
	return
}

// UserCreate creates a user
func (s *storageMySQL) UserCreate(cltPubKey *astimail.PublicKey, srvPrvKey *astimail.PrivateKey) (err error) {
	_, err = s.db.Exec("INSERT INTO user (client_public_key_hash, client_public_key, server_private_key) VALUES ($1, $2)", cltPubKey.Hash(), cltPubKey.String(), srvPrvKey.String())
	return
}

// UserFetchWithEmail fetches a user based on an email
func (s *storageMySQL) UserFetchWithEmail(email string) (u *User, err error) {
	err = s.db.Get(u, "SELECT u.* FROM user u INNER JOIN email e ON u.id = e.user_id WHERE e.addr = $1 AND validated_at IS NOT NULL LIMIT 1", email)
	return
}

// UserFetchWithKey fetches a user based on a key
func (s *storageMySQL) UserFetchWithKey(key *astimail.PublicKey) (u *User, err error) {
	err = s.db.Get(u, "SELECT * FROM user WHERE client_public_key_hash = $1 LIMIT 1", key.Hash())
	return
}

// UserUpdate updates a user
func (s *storageMySQL) UserUpdate(u *User, cltPubKey *astimail.PublicKey, srvPrvKey *astimail.PrivateKey) (err error) {
	_, err = s.db.Exec("UPDATE user SET client_public_key = $1, server_private_key = $2 WHERE id = $3", cltPubKey.String(), srvPrvKey.String(), u.ID)
	return
}

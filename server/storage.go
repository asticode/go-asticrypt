package main

import (
	"database/sql"

	"github.com/asticode/go-astilog"
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
	UpdatedAt mysql.NullTime `db:"updated_at"`
}

// Email represents an email
type Email struct {
	Base
	Addr            string         `db:"addr"`
	ID              int            `db:"id"`
	Token           string         `db:"token"`
	UserID          int            `db:"user_id"`
	ValidatedAt     mysql.NullTime `db:"validated_at"`
	ValidationToken string         `db:"validation_token"`
}

// User represents a user
type User struct {
	Base
	ClientPublicKey     *astimail.PublicKey  `db:"client_public_key"`
	ClientPublicKeyHash []byte               `db:"client_public_key_hash"`
	ID                  int                  `db:"id"`
	ServerPrivateKey    *astimail.PrivateKey `db:"server_private_key"`
}

// Storage represents a storage
type Storage interface {
	EmailCreate(email string, u *User) (token string, err error)
	EmailFetchWithValidationToken(token string) (e *Email, err error)
	EmailList(u *User) (e []*Email, err error)
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
	astilog.Debug("Creating new email")
	token = astistring.RandomString(100)
	_, err = s.db.Exec("INSERT INTO email (addr, user_id, validation_token) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE validation_token = ?", email, u.ID, token, token)
	return
}

// EmailFetchWithValidationToken fetches an email based on a validation token
func (s *storageMySQL) EmailFetchWithValidationToken(token string) (e *Email, err error) {
	astilog.Debug("Fetching email with validation token")
	e = &Email{}
	if err = s.db.Get(e, "SELECT * FROM email WHERE validation_token = ? AND validated_at IS NULL LIMIT 1", token); err == sql.ErrNoRows {
		err = errNotFound
	}
	return
}

// EmailList lists the emails of a user
func (s *storageMySQL) EmailList(u *User) (e []*Email, err error) {
	astilog.Debug("Listing emails")
	e = []*Email{}
	err = s.db.Select(&e, "SELECT * FROM email WHERE user_id = ? AND validated_at IS NOT NULL", u.ID)
	return
}

// EmailValidate validates an email
func (s *storageMySQL) EmailValidate(e *Email) (err error) {
	astilog.Debug("Validating email")
	_, err = s.db.Exec("UPDATE email SET validated_at = NOW() WHERE id = ?", e.ID)
	return
}

// UserCreate creates a user
func (s *storageMySQL) UserCreate(cltPubKey *astimail.PublicKey, srvPrvKey *astimail.PrivateKey) (err error) {
	astilog.Debug("Creating new user")
	_, err = s.db.Exec("INSERT INTO user (client_public_key_hash, client_public_key, server_private_key) VALUES (?, ?, ?)", cltPubKey.Hash(), cltPubKey.String(), srvPrvKey.String())
	return
}

// UserFetchWithEmail fetches a user based on an email
func (s *storageMySQL) UserFetchWithEmail(email string) (u *User, err error) {
	astilog.Debug("Fetching user with email")
	u = &User{}
	if err = s.db.Get(u, "SELECT u.* FROM user u INNER JOIN email e ON u.id = e.user_id WHERE e.addr = ? AND validated_at IS NOT NULL LIMIT 1", email); err == sql.ErrNoRows {
		err = errNotFound
	}
	return
}

// UserFetchWithKey fetches a user based on a key
func (s *storageMySQL) UserFetchWithKey(key *astimail.PublicKey) (u *User, err error) {
	astilog.Debug("Fetching user with key")
	u = &User{}
	if err = s.db.Get(u, "SELECT * FROM user WHERE client_public_key_hash = ? LIMIT 1", key.Hash()); err == sql.ErrNoRows {
		err = errNotFound
	}
	return
}

// UserUpdate updates a user
func (s *storageMySQL) UserUpdate(u *User, cltPubKey *astimail.PublicKey, srvPrvKey *astimail.PrivateKey) (err error) {
	astilog.Debug("Updating user")
	_, err = s.db.Exec("UPDATE user SET client_public_key = ?, server_private_key = ? WHERE id = ?", cltPubKey.String(), srvPrvKey.String(), u.ID)
	return
}

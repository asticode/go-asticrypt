package main

import (
	"database/sql"

	"github.com/asticode/go-asticrypt"
	"github.com/asticode/go-astilog"
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

// Account represents an account
type Account struct {
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
	ClientPublicKey     *asticrypt.PublicKey  `db:"client_public_key"`
	ClientPublicKeyHash []byte                `db:"client_public_key_hash"`
	ID                  int                   `db:"id"`
	ServerPrivateKey    *asticrypt.PrivateKey `db:"server_private_key"`
}

// Storage represents a storage
type Storage interface {
	AccountCreate(account string, u *User) (token string, err error)
	AccountFetchWithValidationToken(token string) (e *Account, err error)
	AccountList(u *User) (e []*Account, err error)
	AccountValidate(e *Account) (err error)
	UserCreate(cltPubKey *asticrypt.PublicKey, srvPrvKey *asticrypt.PrivateKey) error
	UserFetchWithAccount(account string) (*User, error)
	UserFetchWithKey(key *asticrypt.PublicKey) (*User, error)
	UserUpdate(u *User, cltPubKey *asticrypt.PublicKey, srvPrvKey *asticrypt.PrivateKey) error
}

// storageMySQL represents a MySQL storage
type storageMySQL struct {
	db *sqlx.DB
}

// newStorageMySQL builds a new mysql storage
func newStorageMySQL(db *sqlx.DB) *storageMySQL {
	return &storageMySQL{db: db}
}

// AccountCreate creates an account
func (s *storageMySQL) AccountCreate(account string, u *User) (token string, err error) {
	astilog.Debug("Creating new account")
	token = astistring.RandomString(100)
	_, err = s.db.Exec("INSERT INTO account (addr, user_id, validation_token) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE validation_token = ?", account, u.ID, token, token)
	return
}

// AccountFetchWithValidationToken fetches an account based on a validation token
func (s *storageMySQL) AccountFetchWithValidationToken(token string) (e *Account, err error) {
	astilog.Debug("Fetching account with validation token")
	e = &Account{}
	if err = s.db.Get(e, "SELECT * FROM account WHERE validation_token = ? AND validated_at IS NULL LIMIT 1", token); err == sql.ErrNoRows {
		err = errNotFound
	}
	return
}

// AccountList lists the accounts of a user
func (s *storageMySQL) AccountList(u *User) (e []*Account, err error) {
	astilog.Debug("Listing accounts")
	e = []*Account{}
	err = s.db.Select(&e, "SELECT * FROM account WHERE user_id = ? AND validated_at IS NOT NULL", u.ID)
	return
}

// AccountValidate validates an account
func (s *storageMySQL) AccountValidate(e *Account) (err error) {
	astilog.Debug("Validating account")
	_, err = s.db.Exec("UPDATE account SET validated_at = NOW() WHERE id = ?", e.ID)
	return
}

// UserCreate creates a user
func (s *storageMySQL) UserCreate(cltPubKey *asticrypt.PublicKey, srvPrvKey *asticrypt.PrivateKey) (err error) {
	astilog.Debug("Creating new user")
	_, err = s.db.Exec("INSERT INTO user (client_public_key_hash, client_public_key, server_private_key) VALUES (?, ?, ?)", cltPubKey.Hash(), cltPubKey.String(), srvPrvKey.String())
	return
}

// UserFetchWithAccount fetches a user based on an account
func (s *storageMySQL) UserFetchWithAccount(account string) (u *User, err error) {
	astilog.Debug("Fetching user with account")
	u = &User{}
	if err = s.db.Get(u, "SELECT u.* FROM user u INNER JOIN account e ON u.id = e.user_id WHERE e.addr = ? AND validated_at IS NOT NULL LIMIT 1", account); err == sql.ErrNoRows {
		err = errNotFound
	}
	return
}

// UserFetchWithKey fetches a user based on a key
func (s *storageMySQL) UserFetchWithKey(key *asticrypt.PublicKey) (u *User, err error) {
	astilog.Debug("Fetching user with key")
	u = &User{}
	if err = s.db.Get(u, "SELECT * FROM user WHERE client_public_key_hash = ? LIMIT 1", key.Hash()); err == sql.ErrNoRows {
		err = errNotFound
	}
	return
}

// UserUpdate updates a user
func (s *storageMySQL) UserUpdate(u *User, cltPubKey *asticrypt.PublicKey, srvPrvKey *asticrypt.PrivateKey) (err error) {
	astilog.Debug("Updating user")
	_, err = s.db.Exec("UPDATE user SET client_public_key = ?, server_private_key = ? WHERE id = ?", cltPubKey.String(), srvPrvKey.String(), u.ID)
	return
}

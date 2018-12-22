// Copyright (c) 2018. Flying Jamnik Authors
// license that can be found in the LICENSE file.

package csrf

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	InternalServerErr      = errors.New("Internal server error!")
	UnauthorizedRequestErr = errors.New("Unauthorized request error!")
)

type CSRF struct {
	UserID     uint
	Start, End time.Time
	Token      string
	mux        sync.Mutex
}

func RegisterCSRF(userid uint) *CSRF {
	randomNumber := rand.Uint64()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(randomNumber))
	token, _ := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	return &CSRF{userid, time.Now().UTC(), time.Now().Add(15 * time.Minute).UTC(), string(token)}
}

func (c *CSRF) IsActive() bool {
	return time.Now().After(c.End)
}

func (c *CSRF) IsSameToken(token string) bool {
	return c.Token == token
}

func (c *CSRF) Lock() {
	c.mux.Lock()
}

func (c *CSRF) Unlock() {
	c.mux.Unlock()
}

package jwe

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang/glog"
	"gopkg.in/square/go-jose.v2"
	"sync"
)

// KeyHolder is responsible for generating, storing and sync encryption key used for token generation/decryption.
type KeyHolder interface {
	// Returns encrypter instance that can be used to encrypt data.
	Encrypter() jose.Encrypter
	// Returns encryption key that can be used to decrypt data.
	Key() *rsa.PrivateKey
	// Forces refresh of encryption key sync with k8s resource (secret).
	Refresh()
}

// Implements KeyHolder interface.
type rsaKeyHolder struct {
	key *rsa.PrivateKey
	mux sync.Mutex
}

// Encrypter implements key holder interface.
// Used encryption algorithms:
//    - Content encryption: AES-GCM (256)
//    - Key management: RSA-OAEP-SHA256
func (self *rsaKeyHolder) Encrypter() jose.Encrypter {
	publicKey := &self.Key().PublicKey
	encrypter, err := jose.NewEncrypter(jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256, Key: publicKey}, nil)
	if err != nil {
		panic(err)
	}

	return encrypter

}

// Key implements key holder interface.
func (self *rsaKeyHolder) Key() *rsa.PrivateKey {
	self.mux.Lock()
	defer self.mux.Unlock()
	return self.key
}

//TODO(wzt3309) sync with k8s resource.
func (self *rsaKeyHolder) Refresh() {}

//TODO(wzt3309) sync init encryption key with k8s resource (secret).
func (self *rsaKeyHolder) init() {
	self.initEncryptionKey()
}

func (self *rsaKeyHolder) initEncryptionKey() {
	glog.Info("Generating JWE encryption key")
	self.mux.Lock()
	defer self.mux.Unlock()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	self.key = privateKey
}

// NewRSAKeyHolder creates new KeyHolder instance.
func NewRSAKeyHolder() KeyHolder {
	holder := &rsaKeyHolder{}

	holder.init()
	return holder
}

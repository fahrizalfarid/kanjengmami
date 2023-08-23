package kanjengmami

import (
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	c := &Client{
		PrivateKey: "",
		PublicKey:  "",
	}
	err := c.loadPem()
	assert.Equal(t, ErrInvalidPublicKeyPath, err)

	c.PublicKey = "./public_key.pem"
	err = c.loadPem()
	assert.Equal(t, ErrInvalidPrivateKeyPath, err)

	c.PrivateKey = "./private_key.pem"
	err = c.loadPem()
	assert.Nil(t, err)

	var rsaPrivate *rsa.PrivateKey
	var rsaPublic *rsa.PublicKey

	assert.IsType(t, rsaPrivate, c.privateKey)
	assert.IsType(t, rsaPublic, c.publicKey)

	data := []byte("hallo world")
	c.publicKey = nil

	encryptedData, err := c.encryptPacket(data)
	assert.Nil(t, encryptedData)
	assert.Equal(t, ErrInvalidPublicKey, err)

	_ = c.loadPem()
	encryptedData, err = c.encryptPacket(data)
	assert.Nil(t, err)
	assert.Equal(t, 512, len(encryptedData))

	c.privateKey = nil

	decryptedData, err := c.decryptPacket(encryptedData)
	assert.Nil(t, decryptedData)
	assert.Equal(t, ErrInvalidPrivateKey, err)

	_ = c.loadPem()
	decryptedData, err = c.decryptPacket(encryptedData)
	assert.Nil(t, err)
	assert.Equal(t, "hallo world", string(decryptedData))
}

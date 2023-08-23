package kanjengmami

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func (c *Client) loadPem() error {
	if len(c.PublicKey) == 0 {
		return ErrInvalidPublicKeyPath
	}

	if len(c.PrivateKey) == 0 {
		return ErrInvalidPrivateKeyPath
	}

	privKey, err := os.ReadFile(c.PrivateKey)
	if err != nil {
		return err
	}
	pubKey, err := os.ReadFile(c.PublicKey)
	if err != nil {
		return err
	}

	privBlock, _ := pem.Decode(privKey)
	if privBlock == nil {
		return ErrInvalidPrivateBlock
	}

	privRsa, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return err
	}

	pubBlock, _ := pem.Decode(pubKey)
	if err != nil {
		return ErrInvalidPublicBlock
	}
	pubRsa, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return err
	}

	c.privateKey = privRsa
	c.publicKey = pubRsa.(*rsa.PublicKey)

	return nil
}

func (c *Client) encryptPacket(data []byte) ([]byte, error) {
	if c.publicKey == nil {
		return nil, ErrInvalidPublicKey
	}

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, c.publicKey, data, nil)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

func (c *Client) decryptPacket(encryptedData []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, ErrInvalidPrivateKey
	}

	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privateKey, encryptedData, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

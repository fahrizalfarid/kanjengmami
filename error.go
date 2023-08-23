package kanjengmami

import "errors"

var ErrInvalidServers = errors.New("invalid servers")
var ErrInvalidKey = errors.New("invalid key")
var ErrInvalidKeyAndServer = errors.New("invalid secret key or address")
var ErrInvalidPrivateBlock = errors.New("private block is nil")
var ErrInvalidPublicBlock = errors.New("public block is nil")
var ErrInvalidPrivateKey = errors.New("private key is nil")
var ErrInvalidPublicKey = errors.New("public key is nil")
var ErrInvalidPrivateKeyPath = errors.New("invalid private key path")
var ErrInvalidPublicKeyPath = errors.New("invalid public key path")

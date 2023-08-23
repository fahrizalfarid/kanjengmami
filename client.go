package kanjengmami

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/fahrizalfarid/kanjengmami/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	Address               []string
	SecretKey             string
	PrivateKey, PublicKey string
	privateKey            *rsa.PrivateKey
	publicKey             *rsa.PublicKey
}

func NewClient(address []string, SecretKey string,
	privateKeyPath, publicKeyPath string) *Client {
	c := &Client{
		Address:    address,
		SecretKey:  SecretKey,
		PrivateKey: privateKeyPath,
		PublicKey:  publicKeyPath,
	}

	if len(c.PublicKey) != 0 && len(c.PrivateKey) != 0 {
		err := c.loadPem()
		if err != nil {
			panic(err)
		}
	}

	addrs = c.Address
	initNameResolver()
	return c
}

func (c *Client) conn() (*grpc.ClientConn, error) {
	if len(c.SecretKey) == 0 || len(c.Address) == 0 {
		return nil, ErrInvalidKeyAndServer
	}

	cn, _ := grpc.Dial(fmt.Sprintf("%s:///%s", scheme, serviceName),
		grpc.WithDefaultServiceConfig(roundRobin),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))

	try := 0
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		state := cn.GetState()
		if state == 2 {
			break
		}
		if try >= 5 {
			return nil, ErrInvalidServers
		}
		try++
	}

	return cn, nil
}

func (c *Client) cachingClient(conn *grpc.ClientConn, ctx context.Context) (model.CachingClient, context.Context) {
	md := metadata.Pairs("secret-key", c.SecretKey)
	ctx = metadata.NewOutgoingContext(
		ctx, md,
	)
	return model.NewCachingClient(conn), ctx
}

func (c *Client) Set(ctx context.Context, req *model.CacheRequest) error {
	if req.TtlInSecond <= 0 {
		// set to one hour
		req.TtlInSecond = 60 * 60
	}

	conn, err := c.conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	m, ctx := c.cachingClient(conn, ctx)
	_, err = m.Put(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetWithSecure(ctx context.Context, req *model.CacheRequest) error {
	if req.TtlInSecond <= 0 {
		// set to one hour
		req.TtlInSecond = 60 * 60
	}

	encryptedData, err := c.encryptPacket(req.Data)
	if err != nil {
		return err
	}

	req.Data = encryptedData

	conn, err := c.conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	m, ctx := c.cachingClient(conn, ctx)
	_, err = m.Put(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(ctx context.Context, req *model.CacheRequestKey) (*model.CacheResponse, error) {
	if len(req.Key) == 0 {
		return nil, ErrInvalidKey
	}

	conn, err := c.conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	m, ctx := c.cachingClient(conn, ctx)
	data, err := m.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetWithSecure(ctx context.Context, req *model.CacheRequestKey) (*model.CacheResponse, error) {
	if len(req.Key) == 0 {
		return nil, ErrInvalidKey
	}

	conn, err := c.conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	m, ctx := c.cachingClient(conn, ctx)
	data, err := m.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	decryptedData, err := c.decryptPacket(data.Data)
	if err != nil {
		return nil, err
	}

	return &model.CacheResponse{
		Key:  data.Key,
		Data: decryptedData,
	}, nil
}

func (c *Client) Delete(ctx context.Context, req *model.CacheRequestKey) error {
	if len(req.Key) == 0 {
		return ErrInvalidKey
	}

	conn, err := c.conn()
	if err != nil {
		return err
	}
	defer conn.Close()

	m, ctx := c.cachingClient(conn, ctx)
	_, _ = m.Delete(ctx, req)
	return nil
}

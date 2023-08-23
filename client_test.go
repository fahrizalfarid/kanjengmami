// go test -v .
// go test -bench=. -benchmem
// go test -v -coverprofile cover.out -run=. -bench=. -benchmem
// go test -coverprofile cover.out -v .
// go tool cover -html cover.out -o cover.html
package kanjengmami

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fahrizalfarid/kanjengmami/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	*Client
}

const (
	addr = "172.28.182.149"
)

func (s *TestSuite) SetupTest() {
	// c := &Client{
	// 	Address:    []string{"172.28.182.234:80", "172.28.182.234:81", "172.28.182.234:82"},
	// 	SecretKey:  "sEcr3t",
	// 	PublicKey:  "./public_key.pem",
	// 	PrivateKey: "./private_key.pem",
	// }
	// c.loadPem()

	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")

	s.Client = c
}

func TestClient(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func TestConn(t *testing.T) {
	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")

	cn, err := c.conn()
	assert.Nil(t, err)
	assert.NotNil(t, cn)

	// set wrong address
	c = NewClient([]string{"192.168.1.1:8080", "192.168.1.2:8080", "192.168.1.3:8080"},
		"sEcr3t", "./private_key.pem", "./public_key.pem")

	cn, err = c.conn()
	assert.Nil(t, cn)
	assert.NotNil(t, err)

	panicF := func() {
		NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)}, "sEcr3t", "./wrong-path", "./wrong-path")
	}
	require.Panics(t, panicF)
	assert.Panics(t, panicF)
}

func (s *TestSuite) TestWithInsecure() {

	conn, err := s.conn()
	s.Assert().NotNil(conn)
	s.Assert().Nil(err)

	s.Client.SecretKey = ""
	conn, err = s.conn()
	s.Assert().NotNil(err)
	s.Assert().Nil(conn)
	s.Assert().Equal(err, ErrInvalidKeyAndServer)

	s.Client.Address = nil
	conn, err = s.conn()
	s.Assert().NotNil(err)
	s.Assert().Nil(conn)
	s.Assert().Equal(err, ErrInvalidKeyAndServer)

	te := NewClient([]string{"192.168.1.1:8080", "192.168.1.2:8080", "192.168.1.3:8080"},
		"sEcr3t", "./private_key.pem", "./public_key.pem")
	s.Client.SecretKey = "sEcr3t"

	conn, err = te.conn()
	s.Assert().NotNil(err)
	s.Assert().Nil(conn)
	s.Assert().Equal(err, ErrInvalidServers)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Run("set", func() {
		s.Client.Address = nil
		err := s.Client.Set(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})
		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKeyAndServer, err)

		s.Client.SecretKey = ""
		err = s.Client.Set(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})

		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKeyAndServer, err)

		c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
			"sEcr3t", "./private_key.pem", "./public_key.pem")
		s.Client = c

		err = s.Client.Set(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})
		s.Assert().Nil(err)

		err = s.Client.Set(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: -1,
		})
		s.Assert().Nil(err)
	})

	s.Run("get", func() {
		var expectedByte []byte
		var expectedRes *model.CacheResponse

		_ = s.Client.Set(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})

		data, err := s.Client.Get(ctx, &model.CacheRequestKey{
			Key: "coba",
		})

		s.Assert().Nil(err, err)
		s.Assert().IsType(expectedByte, data.Data)
		s.Assert().IsType(expectedRes, data)
		s.Assert().Equal("hallo world", string(data.Data))
		s.Assert().Equal("coba", data.Key)

		data, err = s.Client.Get(ctx, &model.CacheRequestKey{
			Key: "coba1",
		})
		s.Assert().Nil(data)
		s.Assert().NotNil(err)

		data, err = s.Client.Get(ctx, &model.CacheRequestKey{
			Key: "",
		})
		s.Assert().Nil(data)
		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKey, err)
	})

	s.Run("delete", func() {
		err := s.Client.Delete(ctx, &model.CacheRequestKey{
			Key: "coba",
		})
		s.Assert().Nil(err)

		err = s.Client.Delete(ctx, &model.CacheRequestKey{
			Key: "",
		})
		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKey, err)

		err = s.Client.Delete(ctx, &model.CacheRequestKey{
			Key: "coba1",
		})
		s.Assert().Nil(err)
	})
}

func BenchmarkSetInsecure(b *testing.B) {
	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")
	c.loadPem()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(context.TODO(), &model.CacheRequest{
			Key:         fmt.Sprintf("%d", i),
			Data:        []byte(fmt.Sprintf("%d", 1)),
			TtlInSecond: 5,
		})
	}
}

func BenchmarkSetSecure(b *testing.B) {
	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")
	c.loadPem()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.SetWithSecure(context.TODO(), &model.CacheRequest{
			Key:         fmt.Sprintf("%d", i),
			Data:        []byte(fmt.Sprintf("%d", 1)),
			TtlInSecond: 5,
		})
	}
}

func BenchmarkGetInsecure(b *testing.B) {
	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")
	c.loadPem()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Get(context.TODO(), &model.CacheRequestKey{
			Key: fmt.Sprintf("%d", i),
		})
	}
}

func BenchmarkGetSecure(b *testing.B) {
	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")
	c.loadPem()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.GetWithSecure(context.TODO(), &model.CacheRequestKey{
			Key: fmt.Sprintf("%d", i),
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	c := NewClient([]string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)},
		"sEcr3t", "./private_key.pem", "./public_key.pem")
	c.loadPem()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Delete(context.TODO(), &model.CacheRequestKey{
			Key: fmt.Sprintf("%d", i),
		})
	}
}

func (s *TestSuite) TestWithSecured() {

	conn, err := s.conn()
	s.Assert().NotNil(conn)
	s.Assert().Nil(err)

	s.Client.SecretKey = ""
	conn, err = s.conn()
	s.Assert().NotNil(err)
	s.Assert().Nil(conn)
	s.Assert().Equal(err, ErrInvalidKeyAndServer)

	s.Client.Address = nil
	conn, err = s.conn()
	s.Assert().NotNil(err)
	s.Assert().Nil(conn)
	s.Assert().Equal(err, ErrInvalidKeyAndServer)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Run("set", func() {
		s.Client.Address = nil
		err := s.Client.SetWithSecure(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})
		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKeyAndServer, err)

		s.Client.SecretKey = ""
		err = s.Client.SetWithSecure(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})

		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKeyAndServer, err)

		s.Client.Address = []string{fmt.Sprintf("%s:1081", addr), fmt.Sprintf("%s:1082", addr), fmt.Sprintf("%s:1083", addr)}
		s.Client.SecretKey = "sEcr3t"

		err = s.Client.SetWithSecure(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})
		s.Assert().Nil(err)

		err = s.Client.SetWithSecure(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: -1,
		})
		s.Assert().Nil(err)
	})

	s.Run("get", func() {
		var expectedByte []byte
		var expectedRes *model.CacheResponse

		_ = s.Client.SetWithSecure(ctx, &model.CacheRequest{
			Key:         "coba",
			Data:        []byte("hallo world"),
			TtlInSecond: 3 * 60 * 60,
		})

		data, err := s.Client.GetWithSecure(ctx, &model.CacheRequestKey{
			Key: "coba",
		})

		s.Assert().Nil(err, err)
		s.Assert().IsType(expectedByte, data.Data)
		s.Assert().IsType(expectedRes, data)
		s.Assert().Equal("hallo world", string(data.Data))
		s.Assert().Equal("coba", data.Key)

		data, err = s.Client.GetWithSecure(ctx, &model.CacheRequestKey{
			Key: "coba1",
		})
		s.Assert().Nil(data)
		s.Assert().NotNil(err)

		data, err = s.Client.GetWithSecure(ctx, &model.CacheRequestKey{
			Key: "",
		})
		s.Assert().Nil(data)
		s.Assert().NotNil(err)
		s.Assert().Equal(ErrInvalidKey, err)
	})
}

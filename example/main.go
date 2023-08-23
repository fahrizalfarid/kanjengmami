package main

import (
	"context"
	"fmt"
	"time"

	k "github.com/fahrizalfarid/kanjengmami"
	"github.com/fahrizalfarid/kanjengmami/model"
)

const timeout time.Duration = 1

func set(c *k.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	err := c.Set(ctx, &model.CacheRequest{
		Key:         "coba",
		Data:        []byte("hallo world"),
		TtlInSecond: -1,
	})
	fmt.Println("set", err)
}

func setWithSecure(c *k.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	err := c.SetWithSecure(ctx, &model.CacheRequest{
		Key:         "cobaSecure",
		Data:        []byte("hallo world"),
		TtlInSecond: -1,
	})
	fmt.Println("setWithSecure", err)
}

func get(c *k.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	data, err := c.Get(ctx, &model.CacheRequestKey{
		Key: "coba",
	})
	fmt.Printf("get %s %v\n", data, err)
}

func getWithSecure(c *k.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	data, err := c.GetWithSecure(ctx, &model.CacheRequestKey{
		Key: "cobaSecure",
	})
	fmt.Printf("getWithSecure %s %v\n", data, err)
}

func delete(c *k.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	err := c.Delete(ctx, &model.CacheRequestKey{
		Key: "coba",
	})
	fmt.Println("delete", err)
}

func main() {
	c := k.NewClient([]string{"172.28.182.234:1081", "172.28.182.234:1082", "172.28.182.234:1083"},
		"sEcr3t",
		"./private_key.pem",
		"./public_key.pem",
	)
	set(c)
	setWithSecure(c)
	get(c)
	getWithSecure(c)
	delete(c)
}

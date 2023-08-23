## pull your image
### single vm
```yml
version: "3.9"

services:
  kanjengmami-single:
    image: fahrizalfarid/kanjeng-mami
    container_name: kanjengmami-0
    command:
      - "--api-port=1081"
      - "-a=1@127.0.0.1:8081"
      - "--join=-1"
      - "--secret-key=sEcr3t"
    ports:
      - 1081:1081
      - 8081:8081
    volumes:
      - ${HOME}/kanjengmami:/data
    restart: always
```

### cluster
```yml
version: "3.9"

services:
  kanjengmami-leader:
    image: fahrizalfarid/kanjeng-mami
    container_name: kanjengmami-leader
    command:
      - "--api-port=1081"
      - "--cluster-id=128"
      - "--secret-key=sEcr3t"
      - "-a=1@127.0.0.1:8081"
      - "--join=0"
      - "-c=1@127.0.0.1:8081"
      - "-c=2@127.0.0.1:8082"
      - "-c=3@127.0.0.1:8083"
    ports:
      - 1081:1081
      - 127.0.0.1:8081:8081
    restart: always

  kanjengmami-0:
    image: fahrizalfarid/kanjeng-mami
    container_name: kanjengmami-0
    command:
      - "--api-port=1082"
      - "--cluster-id=128"
      - "--secret-key=sEcr3t"
      - "-a=1@127.0.0.1:8082"
      - "--join=1"  
    ports:
      - 1082:1082
      - 127.0.0.1:8082:8082
    
    restart: always

  kanjengmami-1:
    image: fahrizalfarid/kanjeng-mami
    container_name: kanjengmami-1
    command:
      - "--api-port=1083"
      - "--cluster-id=128"
      - "--secret-key=sEcr3t"
      - "-a=1@127.0.0.1:8083"
      - "--join=1"
    ports:
      - 1083:1083
      - 127.0.0.1:8083:8083
    restart: always
```

## regenerate proto

```$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative model/kanjengMami.proto```


## example

```go
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
```
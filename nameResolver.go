package kanjengmami

import "google.golang.org/grpc/resolver"

const (
	scheme      = "kanjengmami"
	serviceName = "caching.grpc.io"
	roundRobin  = `{"loadBalancingConfig": [{"round_robin":{}}]}`
)

var addrs []string

type resolverBuilder struct{}
type nameResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	addrs  map[string][]string
}

func (*resolverBuilder) Scheme() string { return scheme }
func (r *nameResolver) start() {
	addrStrs := r.addrs[r.target.Endpoint()]
	addrss := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrss[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{
		Addresses: addrss,
	})
}

func (*nameResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*nameResolver) Close()                                  {}

func (*resolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &nameResolver{
		target: target,
		cc:     cc,
		addrs: map[string][]string{
			serviceName: addrs,
		},
	}
	r.start()
	return r, nil
}

func initNameResolver() {
	resolver.Register(&resolverBuilder{})
}

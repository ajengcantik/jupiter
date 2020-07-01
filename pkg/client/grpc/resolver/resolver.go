package resolver

import (
	"context"

	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Register ...
func Register(name string, reg registry.Registry) {
	resolver.Register(&baseBuilder{
		name: name,
		reg:  reg,
	})
}

type baseBuilder struct {
	name string
	reg  registry.Registry
}

// Build ...
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	endpoints, err := b.reg.WatchServices(context.Background(), target.Endpoint, "grpc")
	if err != nil {
		return nil, err
	}

	var stop = make(chan struct{})
	xgo.Go(func() {
		for {
			select {
			case endpoint := <-endpoints:
				var state = resolver.State{
					Addresses:  make([]resolver.Address, 0),
					Attributes: attributes.New("config", endpoint.RouteConfigs),
				}
				for _, node := range endpoint.Nodes {
					var address resolver.Address
					address.Addr = node.Address
					address.ServerName = target.Endpoint
					state.Addresses = append(state.Addresses, address)
				}
				cc.UpdateState(state)
			case <-stop:
				return
			}
		}
	})

	return &baseResolver{
		stop: stop,
	}, nil
}

// Scheme ...
func (b baseBuilder) Scheme() string {
	return b.name
}

type baseResolver struct {
	stop chan struct{}
}

// ResolveNow ...
func (b *baseResolver) ResolveNow(options resolver.ResolveNowOptions) {}

// Close ...
func (b *baseResolver) Close() { b.stop <- struct{}{} }

package taskfilter

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

type ACWrapper struct {
	client.Client
}

func NewACWrapper(c client.Client) client.Client {
	return &ACWrapper{c}
}

func (p *ACWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)

	filter := func(services []*registry.Service) []*registry.Service {
		for _, service := range services {
			var nodes []*registry.Node
			for _, node := range service.Nodes {
				// some filter here to filter, centerid, memory, cpu, network, bandwidth, etc...
				if node.Metadata["datacenter"] == md["datacenter"] {
					nodes = append(nodes, node)
				}
			}
			service.Nodes = nodes
		}
		return services
	}

	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(filter),
	))

	fmt.Printf("[DC Wrapper] filtering for datacenter %s\n", md["datacenter"])
	return p.Client.Call(ctx, req, rsp, callOptions...)
}


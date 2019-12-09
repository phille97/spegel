package discovery

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/grandcat/zeroconf"
)

func Discover(service string) (*[]url.URL, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	allEntries := []url.URL{}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			allEntries = append(allEntries, url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", entry.HostName, entry.Port),
			})
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := resolver.Browse(ctx, service, "local.", entries); err != nil {
		return nil, err
	}

	<-ctx.Done()

	return &allEntries, nil
}

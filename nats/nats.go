// https://github.com/testcontainers/testcontainers-go/blob/v0.22.0/examples/nats/nats.go

package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NatsContainer represents the nats container type used in the module
type NatsContainer struct {
	testcontainers.Container
	URI string
}

// RunContainer creates an instance of the nats container type
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*NatsContainer, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "nats:latest",
		Hostname:     "127.0.0.1",
		ExposedPorts: []string{"4222/tcp", "6222/tcp", "8222/tcp"},
		Cmd:          []string{"-DV", "-js"},
		WaitingFor:   wait.ForLog("Listening for client connections on 0.0.0.0:4222"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "4222/tcp")
	if err != nil {
		return nil, nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	uri := fmt.Sprintf("nats://%s:%s", hostIP, mappedPort.Port())

	natsContainer := &NatsContainer{Container: container, URI: uri}

	cleanupFunc := func() {
		if err := natsContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	return natsContainer, cleanupFunc, nil
}

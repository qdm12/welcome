package docker

import (
	"context"
)

func (d *Docker) IsRunning(ctx context.Context) (running bool) {
	_, err := d.client.Ping(ctx)
	return err == nil
}

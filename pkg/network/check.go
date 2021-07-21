package network

import "context"

func (n *Network) Check(ctx context.Context, urls []string) (errors []error) {
	return n.connChecker.ParallelChecks(ctx, urls)
}

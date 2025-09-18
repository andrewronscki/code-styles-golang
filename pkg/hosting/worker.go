package hosting

import "context"

type Worker interface {
	Run(ctx context.Context, exit func())
}

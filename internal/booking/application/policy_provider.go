package application

import (
	"context"

	"github.com/spayder/bookery/internal/booking/domain"
)

type PolicyProvider interface {
	PolicyFor(ctx context.Context, resourceID domain.ResourceID) (domain.CancellationPolicy, error)
}

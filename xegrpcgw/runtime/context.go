package runtime

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type serverMetadataKey struct{}

// NewServerMetadataContext creates a new context with ServerMetadata
func NewServerMetadataContext(ctx context.Context, md runtime.ServerMetadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// ServerMetadataFromContext returns the ServerMetadata in ctx
func ServerMetadataFromContext(ctx context.Context) (md runtime.ServerMetadata, ok bool) {
	md, ok = ctx.Value(serverMetadataKey{}).(runtime.ServerMetadata)
	return
}

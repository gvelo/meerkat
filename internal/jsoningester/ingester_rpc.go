package jsoningester

import (
	"context"
	"meerkat/internal/cluster"
	"meerkat/internal/ingestion"
)

type IngesterRpc interface {
	SendRequest(ctx context.Context, member string, request *ingestion.IngestionRequest) error
}

func NewIngestRPC(cl cluster.Cluster) IngesterRpc {
	return &ingestRpc{
		cl: cl,
	}
}

type ingestRpc struct {
	cl cluster.Cluster
}

func (i ingestRpc) SendRequest(ctx context.Context, nodeId string, request *ingestion.IngestionRequest) error {

	node := i.cl.Node(nodeId)

	if node == nil {
		// TODO(gvelo): handle
		panic("member not found")
	}

	cl := ingestion.NewIngesterClient(node.ClientConn())
	_, err := cl.Ingest(ctx, request)

	return err

}

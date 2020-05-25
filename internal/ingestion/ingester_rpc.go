package ingestion

import (
	"context"
	"meerkat/internal/cluster"
	"meerkat/internal/ingestion/ingestionpb"
)

type IngesterRpc interface {
	SendRequest(ctx context.Context, member string, request *ingestionpb.IngestionRequest) error
}

func NewIngestRPC(connReg cluster.ConnRegistry) IngesterRpc {
	return &ingestRpc{
		connReg: connReg,
	}
}

type ingestRpc struct {
	connReg cluster.ConnRegistry
}

func (i ingestRpc) SendRequest(ctx context.Context, member string, request *ingestionpb.IngestionRequest) error {

	c := i.connReg.ClientConn(member)

	if c == nil {
		// TODO(gvelo): handle
		panic("member not found")
	}

	cl := ingestionpb.NewIngesterClient(c)
	_, err := cl.Ingest(ctx, request)

	return err

}

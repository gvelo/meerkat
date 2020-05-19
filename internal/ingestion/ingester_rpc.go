package ingestion

import "meerkat/internal/ingestion/ingestionpb"

type IngesterRpc interface {
	SendRequest(member string, request ingestionpb.IngestionRequest)
}



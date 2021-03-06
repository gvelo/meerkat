// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"meerkat/internal/build"
	"meerkat/internal/cluster"
	"meerkat/internal/config"
	"meerkat/internal/ingestion"
	"meerkat/internal/jsoningester"
	"meerkat/internal/query/exec"
	"meerkat/internal/query/execpb"
	"meerkat/internal/query/physical"
	"meerkat/internal/rest"
	"meerkat/internal/storage"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Meerkat struct {
	mu                sync.Mutex
	listener          net.Listener
	grpcServer        *grpc.Server
	clusterMgr        cluster.Manager
	apiServer         *rest.ApiServer
	catalog           cluster.Catalog
	Conf              config.Config
	log               zerolog.Logger
	segmentWriterPool *storage.SegmentWriterPool
	bufReg            ingestion.BufferRegistry
	segStorage        storage.SegmentStorage
	segRegistry       storage.SegmentRegistry
}

func (m *Meerkat) Start(ctx context.Context) {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.initLogger()

	m.log.Info().Msgf("starting meerkat %v (%v)", build.Version, build.Commit)

	// start components

	var err error

	// TODO(gvelo):make port configurable.
	m.listener, err = net.Listen("tcp", ":9191")

	if err != nil {
		m.log.Panic().Err(err).Msg("failed to listen: 9191")
	}

	m.grpcServer = grpc.NewServer()

	var seeds []string
	if m.Conf.Seeds == "" {
		seeds = make([]string, 0)
	} else {
		seeds = strings.Split(m.Conf.Seeds, ",")
	}

	m.clusterMgr, err = cluster.NewManager(m.Conf.GossipPort, seeds, m.Conf.DBPath)

	if err != nil {
		m.log.Panic().Err(err).Msg("cannot create cluster")
	}

	err = m.clusterMgr.Join()

	if err != nil {
		m.log.Panic().Err(err).Msg("cannot join cluster")
	}

	catalogRPC, err := cluster.NewCatalogRPC(m.clusterMgr)

	if err != nil {
		m.log.Panic().Err(err).Msg("cannot create catalogRPC")
		return
	}

	m.catalog, err = cluster.NewCatalog(m.grpcServer, m.Conf.DBPath, m.clusterMgr, catalogRPC)

	if err != nil {
		m.log.Panic().Err(err).Msg("cannot create catalog")
	}

	// m.schema, err = schema.NewSchema(m.catalog)
	//
	// if err != nil {
	//	m.log.Panic().Err(err).Msg("cannot create schema")
	// }

	m.segStorage = storage.NewStorage(m.Conf.DBPath)

	m.segRegistry = storage.NewSegmentRegistry(m.Conf.DBPath, m.segStorage)

	m.segRegistry.Start()

	m.segmentWriterPool = storage.NewSegmentWriterPool(1024,
		10,
		m.segStorage,
		m.segRegistry,
	)

	m.segmentWriterPool.Start()

	ingRcp := jsoningester.NewIngestRPC(m.clusterMgr)

	// TODO(gvelo): all this params should be externalized to conf.
	m.bufReg = ingestion.NewBufferRegistry(1024,
		10*time.Second,
		-1,
		1*time.Second,
		32,
		m.segmentWriterPool,
	)

	m.bufReg.Start()

	streamReg := physical.NewStreamRegistry()

	execServer := exec.NewServer(m.clusterMgr, m.segRegistry, streamReg)

	execpb.RegisterExecutorServer(m.grpcServer, execServer)

	executor := exec.NewExecutor(m.segRegistry, streamReg, m.clusterMgr)

	m.apiServer, err = rest.NewRestApi(m.clusterMgr, ingRcp, m.bufReg, executor)

	if err != nil {
		m.log.Panic().Err(err).Msg("cannot create rest server")
	}

	m.apiServer.Start()

	// ingestion server
	ingServer := ingestion.NewServer(m.bufReg)
	ingestion.RegisterIngesterServer(m.grpcServer, ingServer)

	go func() {
		err = m.grpcServer.Serve(m.listener)
		if err != nil {
			m.log.Error().Err(err).Msg("error serving grpc connections")
			return
		}
		m.log.Info().Msg("grpc server stopped")
	}()

	err = m.clusterMgr.Ready()

	m.log.Info().Msg("meerkat server started")

}

func (m *Meerkat) Shutdown(ctx context.Context) {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Info().Msg("stopping meerkar server")

	err := m.clusterMgr.Leaving()

	if err != nil {
		m.log.Error().Err(err).Msg("error leaving cluster")
	}

	err = m.apiServer.Shutdown(ctx)

	if err != nil {
		m.log.Error().Err(err).Msg("error stopping api server")
	}

	m.grpcServer.GracefulStop()

	m.bufReg.Stop()

	m.segmentWriterPool.Stop()

	err = m.catalog.Shutdown()

	if err != nil {
		m.log.Error().Err(err).Msg("error stopping catalog")
	}

	m.segRegistry.Stop()

	err = m.clusterMgr.Leave()

	if err != nil {
		m.log.Error().Err(err).Msg("error leaving cluster")
	}

	m.log.Info().Msg("meerkat server stopped")

}

func (m *Meerkat) initLogger() {

	// Default level is info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if m.Conf.LogLevel != "" {
		l, err := zerolog.ParseLevel(m.Conf.LogLevel)
		if err != nil {
			panic(err)
		}
		zerolog.SetGlobalLevel(l)
	}

	if m.Conf.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Logger = log.With().Caller().Logger()

	m.log = log.With().Str("src", "server").Logger()

}

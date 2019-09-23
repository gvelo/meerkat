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
	"meerkat/internal/rest"
	"meerkat/internal/schema"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	mu            sync.Mutex
	listener      net.Listener
	grpcServer    *grpc.Server
	cl            cluster.Cluster
	connRegistry  cluster.ConnRegistry
	schemaService schema.Schema
	apiServer     *rest.ApiServer
	catalog       cluster.Catalog
	l             zerolog.Logger
)

func Start(ctx context.Context, c config.Config) {

	mu.Lock()
	defer mu.Unlock()

	// Default level is info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if c.LogLevel != "" {
		l, err := zerolog.ParseLevel(c.LogLevel)
		if err != nil {
			panic(err)
		}
		zerolog.SetGlobalLevel(l)
	}

	if c.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Logger = log.With().Caller().Logger()

	l = log.With().Str("src", "server").Logger()

	l.Info().Msgf("starting meerkat %v (%v)", build.Version, build.Commit)

	// start components

	var err error

	listener, err = net.Listen("tcp", ":9191")

	if err != nil {
		l.Panic().Err(err).Msg("failed to listen: 9191")
	}

	grpcServer = grpc.NewServer()

	var seeds []string
	if c.Seeds == "" {
		seeds = make([]string, 0)
	} else {
		seeds = strings.Split(c.Seeds, ",")
	}

	cl, err = cluster.NewCluster(c.GossipPort, seeds, c.DBPath)
	if err != nil {
		l.Panic().Err(err).Msg("cannot create cluster")
	}

	connRegistry, err = cluster.NewConnRegistry(cl)

	if err != nil {
		l.Panic().Err(err).Msg("cannot create catalogRPC")
		return
	}

	catalogRPC, err := cluster.NewCatalogRPC(connRegistry)

	if err != nil {
		l.Panic().Err(err).Msg("cannot create catalogRPC")
		return
	}

	catalog, err = cluster.NewCatalog(grpcServer, c.DBPath, cl, catalogRPC)

	if err != nil {
		l.Panic().Err(err).Msg("cannot create catalog")
	}

	schemaService, err = schema.NewSchema(catalog)

	if err != nil {
		l.Panic().Err(err).Msg("cannot create schema")
	}

	apiServer, err = rest.NewRest(schemaService)

	if err != nil {
		l.Panic().Err(err).Msg("cannot create rest server")
	}

	apiServer.Start()

	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			l.Error().Err(err).Msg("error serving grpc connections")
			return
		}
		l.Info().Msg("grpc server stopped")
	}()

	cl.Join()

	l.Info().Msg("server started")

}

func Shutdown(ctx context.Context) {

	l.Info().Msg("stopping server")

	cl.Shutdown()

	err := apiServer.Shutdown(ctx)

	if err != nil {
		l.Error().Err(err).Msg("error stopping api server")
	}

	grpcServer.GracefulStop()

	schemaService.Shutdown()

	err = catalog.Shutdown()

	if err != nil {
		l.Error().Err(err).Msg("error stopping catalog")
	}

	connRegistry.Shutdown()

	l.Info().Msg("server stopped")

}

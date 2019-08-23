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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"meerkat/internal/build"
	"meerkat/internal/cluster"
	"meerkat/internal/config"
	"net"
	"os"
	"strings"
)

func Start(c config.Config) {

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

	log.Info().Msgf("starting meerkat %v (%v)", build.Version, build.Commit)
	// start components

	lis, err := net.Listen("tcp", ":9191")

	if err != nil {
		log.Panic().Err(err).Msg("failed to listen: 9191")
	}

	grpcServer := grpc.NewServer()

	var seeds []string
	if c.Seeds == "" {
		seeds = make([]string, 0)
	} else {
		seeds = strings.Split(c.Seeds, ",")
	}

	cl, err := cluster.NewCluster(c.GossipPort, seeds, c.DBPath)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create cluster")
		return
	}

	transport, err := cluster.NewTransport(cl)

	if err != nil {
		log.Panic().Err(err).Msg("cannot create transport")
		return
	}

	_, err = cluster.NewCatalog(grpcServer, c.DBPath, cl, transport)

	if err != nil {
		log.Panic().Err(err).Msg("cannot create catalog")
	}

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Panic().Err(err).Msg("cannot create grpc server")
			return
		}
	}()

	cl.Join()

}

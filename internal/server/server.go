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
	"meerkat/internal/build"
	"meerkat/internal/cluster"
	"meerkat/internal/config"
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

	seeds := strings.Split(c.Seeds, ",")

	cl := cluster.NewCluster(seeds, c.DBPath)

	err := cl.Start()

	if err != nil {
		log.Panic().Err(err).Msg("cannot start cluster")
		return
	}

	cl.Ready()

}

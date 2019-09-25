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

package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"meerkat/internal/build"
	"meerkat/internal/config"
	"meerkat/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rootCmd *cobra.Command
var logLevel string
var configFile string
var dbpath string
var seeds string
var gossipPort int

func init() {

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "print build information",
		Long:  "print build information",
		Run:   Version,
	}

	var cmdStart = &cobra.Command{
		Use:   "start",
		Short: "start the meerkat server",
		Long:  "start the meerkat server",
		Run:   Start,
	}

	rootCmd = &cobra.Command{Use: "meerkat"}
	rootCmd.AddCommand(cmdVersion, cmdStart)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "set the log level [panic,fatal,error,warn,info,debug] ")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&dbpath, "dbpath", "d", "", "database files path")
	rootCmd.PersistentFlags().StringVarP(&seeds, "seeds", "s", "", "the IP addresses of the clusters seed servers")
	rootCmd.PersistentFlags().IntVarP(&gossipPort, "gossip-port", "g", -1, "the gossip bind port")

	err := viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("log-level"))
	err = viper.BindPFlag("dbpath", rootCmd.PersistentFlags().Lookup("dbpath"))
	err = viper.BindPFlag("seeds", rootCmd.PersistentFlags().Lookup("seeds"))
	err = viper.BindPFlag("gossipPort", rootCmd.PersistentFlags().Lookup("gossip-port"))

	if err != nil {
		panic(err)
	}

}

func main() {

	fmt.Printf("meerkat %v ( %v ) \n\n", build.Version, build.Commit)

	err := rootCmd.Execute()

	if err != nil {
		panic(err)
	}

}

func Version(cmd *cobra.Command, args []string) {
	fmt.Printf("version: %v \n", build.Version)
	fmt.Printf("branch: %v \n", build.Branch)
	fmt.Printf("tag: %v \n", build.Tag)
	fmt.Printf("build user: %v \n", build.BuildUser)
	fmt.Printf("build date: %v \n", build.BuildDate)
	fmt.Printf("revision: %v \n", build.Commit)
	fmt.Printf("go version: %v \n", build.GoVersion)
	fmt.Printf("platform: %v \n", build.Platform)
}

func Start(cmd *cobra.Command, args []string) {

	if configFile != "" {

		viper.SetConfigFile(configFile)
		viper.AutomaticEnv()
		err := viper.ReadInConfig()

		if err != nil {
			panic(err)
		}

	}

	var conf config.Config

	err := viper.Unmarshal(&conf)

	if err != nil {
		panic(err)
	}

	meerkat := &server.Meerkat{Conf: conf}

	ctx := serverCtx()

	meerkat.Start(ctx)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	meerkat.Shutdown(ctx)

}

func serverCtx() context.Context {

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		defer cancel()

		select {
		case <-ctx.Done():
			return
		case <-signalCh:
			return
		}

	}()

	return ctx

}

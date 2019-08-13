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
	"flag"
	"fmt"
	"meerkat/internal/build"
	"meerkat/internal/query/rel"
	"os"
	"strings"
)

func main() {

	fmt.Printf("Buid version: %v \n ", build.Version)
	fmt.Printf("Buid branch: %v \n ", build.Branch)
	fmt.Printf("Buid tag: %v \n ", build.Tag)
	fmt.Printf("Buid commit: %v \n ", build.Commit)
	fmt.Printf("Buid build user: %v \n ", build.BuildUser)
	fmt.Printf("Buid build date: %v \n ", build.BuildDate)
	fmt.Printf("go version: %v \n ", build.GoVersion)
	fmt.Printf("platform: %v \n ", build.Platform)

	infoCmd := flag.NewFlagSet("-q", flag.ExitOnError)

	// The subcommand is expected as the first argument
	// to the program.
	if len(os.Args) < 1 {
		fmt.Println("expected '-q' for the query")
		os.Exit(1)
	}

	// Check which subcommand is invoked.
	switch os.Args[1] {

	// For every subcommand, we parse its own flags and
	// have access to trailing positional arguments.
	case "-q":
		infoCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'q'")
		fmt.Println("  tail:", infoCmd.Args())
		// ie: -q "stats avg(a) by field1"
		str := strings.Join(infoCmd.Args(), "")
		rel.ProcessQuery(str)
	default:
		fmt.Println("No command")
		os.Exit(1)
	}

}

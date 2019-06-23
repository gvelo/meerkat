package main

import (
	"flag"
	"fmt"
	"meerkat/internal/storage/readers"
	"os"
)

func main() {

	infoCmd := flag.NewFlagSet("stats", flag.ExitOnError)
	path := infoCmd.String("p", ".", "path")

	barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	barLevel := barCmd.Int("level", 0, "level")

	// The subcommand is expected as the first argument
	// to the program.
	if len(os.Args) < 2 {
		fmt.Println("expected 'stats' or 'bar' subcommands")
		os.Exit(1)
	}

	// Check which subcommand is invoked.
	switch os.Args[1] {

	// For every subcommand, we parse its own flags and
	// have access to trailing positional arguments.
	case "stats":
		infoCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'stats'")
		fmt.Println("  path:", *path)
		fmt.Println("  tail:", infoCmd.Args())
		printStats(*path)
	case "bar":
		barCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'bar'")
		fmt.Println("  level:", *barLevel)
		fmt.Println("  tail:", barCmd.Args())
	default:
		fmt.Println("Show command")
		os.Exit(1)
	}

}

func printStats(path string) {
	s, err := readers.ReadSegment(path)
	if err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
	fmt.Printf("===  Segment Name : %s  === \n", s.IndexInfo.Name)
	fmt.Printf("== %d Fields:\n", len(s.IndexInfo.Fields))
	for _, f := range s.IndexInfo.Fields {
		fmt.Printf("= [Name:%s], [Type:%d], [Id,%d], [indexed:%v] \n", f.Name, f.Type, f.ID, f.Index)
	}

}

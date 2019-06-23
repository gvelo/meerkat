package main

import (
	"fmt"
	"meerkat/internal/build"
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
}

package main

import (
	"fmt"

	"github.com/erikw/gomemo/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Starting Gomemo with\n%v\n", cfg)
}

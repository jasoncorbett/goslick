package main

import "github.com/serussell/logxi/v1"
import "github.com/jasoncorbett/goslick/slickconfig"

func main() {
	logger := log.New("slickcli")
	slickconfig.Configuration.LoadFromStandardLocations()
	logger.Info("Configuration parsed.")
}

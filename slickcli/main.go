package main

import "github.com/serussell/logxi/v1"

func main() {
	logger := log.New("slickcli")
	logger.Debug("Debug Message", "foo", "bar")
	logger.Info("Info Message", "foo", "bar")
	logger.Warn("Warn Message", "foo", "bar")
	logger.Error("Error Message", "foo", "bar")
}

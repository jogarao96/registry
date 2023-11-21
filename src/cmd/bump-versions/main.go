package main

import (
	"flag"
	"log/slog"
	"os"
	"registry-stable/internal/provider"
	"registry-stable/internal/repository-metadata-files/module"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // TODO REMOVE ME
	logger.Info("Starting version bump process for modules and providers")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")

	flag.Parse()

	modules, err := module.ListModules(*moduleDataDir)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	providers, err := provider.ListProviders(*providerDataDir, logger)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}

	for _, m := range modules {
		slog.Info("Beginning version bump process for module", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.TargetSystem))
		err = module.UpdateMetadataFile(m, *moduleDataDir)
		if err != nil {
			slog.Error("Failed to version bump module", slog.Any("err", err))
			os.Exit(1)
		}
	}

	for _, p := range providers {
		p.Logger.Info("Beginning version bump process")
		err = p.UpdateMetadataFile()
		if err != nil {
			p.Logger.Error("Failed to version bump provider", slog.Any("err", err))
			os.Exit(1)
		}
	}

	slog.Info("Completed version bump process for modules and providers")
}

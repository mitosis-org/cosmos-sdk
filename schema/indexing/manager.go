package indexing

import (
	"context"
	"fmt"

	"cosmossdk.io/schema/appdata"
	"cosmossdk.io/schema/decoding"
)

type Options struct {
	Context    context.Context
	Options    map[string]interface{}
	Resolver   decoding.DecoderResolver
	SyncSource decoding.SyncSource
}

func Start(opts Options) (appdata.Listener, error) {
	var indexers []appdata.Listener
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	for indexerName, factory := range indexerRegistry {
		indexerOpts, ok := opts.Options[indexerName]
		if !ok {
			continue
		}

		optsMap, ok := indexerOpts.(map[string]interface{})
		if !ok {
			return appdata.Listener{}, fmt.Errorf("invalid indexer options type %T for %s, expected a map", indexerOpts, indexerName)
		}

		indexer, err := factory(optsMap)
		if err != nil {
			return appdata.Listener{}, fmt.Errorf("failed to create indexer %s: %w", indexerName, err)
		}

		res, err := indexer.Initialize(ctx, InitializationData{})
		if err != nil {
			return appdata.Listener{}, fmt.Errorf("failed to initialize indexer %s: %w", indexerName, err)
		}

		indexers = append(indexers, res.Listener)

		// TODO handle last block persisted
	}

	return decoding.Middleware(appdata.AsyncListenerMux(indexers, 1024, ctx.Done()), decoding.Options{
		DecoderResolver: opts.Resolver,
		SyncSource:      opts.SyncSource,
	})
}

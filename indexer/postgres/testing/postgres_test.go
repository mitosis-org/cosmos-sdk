package testing

import (
	"context"
	"os"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/schema/appdata"
	indexertesting "cosmossdk.io/schema/testing"
	appdatatest "cosmossdk.io/schema/testing/appdata"

	_ "github.com/jackc/pgx/v5/stdlib"

	"cosmossdk.io/indexer/postgres"
)

func TestPostgresIndexer(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "postgres-indexer-test")
	require.NoError(t, err)

	dbPort := freeport.GetOne(t)
	pgConfig := embeddedpostgres.DefaultConfig().
		Port(uint32(dbPort)).
		DataPath(tempDir)

	dbUrl := pgConfig.GetConnectionURL()
	pg := embeddedpostgres.NewDatabase(pgConfig)
	require.NoError(t, pg.Start())

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
		require.NoError(t, pg.Stop())
		err := os.RemoveAll(tempDir)
		require.NoError(t, err)
	})

	indexer, err := postgres.NewIndexer(ctx, postgres.Options{
		Driver:        "pgx",
		ConnectionURL: dbUrl,
	})
	require.NoError(t, err)

	fixture := appdatatest.NewSimulator(appdatatest.SimulatorOptions{
		Listener: appdata.ListenerMux(
			appdata.DebugListener(os.Stdout),
			indexer.Listener(),
		),
		AppSchema: indexertesting.ExampleAppSchema,
	})

	require.NoError(t, fixture.Initialize())

	blockDataGen := fixture.SimulateBlockGenN(1000)
	for i := 0; i < 1000; i++ {
		blockDataGen.Example(i)
	}
}

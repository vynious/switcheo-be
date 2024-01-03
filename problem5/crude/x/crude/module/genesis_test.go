package crude_test

import (
	"testing"

	keepertest "crude/testutil/keeper"
	"crude/testutil/nullify"
	"crude/x/crude/module"
	"crude/x/crude/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		UserList: []types.User{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		UserCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CrudeKeeper(t)
	crude.InitGenesis(ctx, k, genesisState)
	got := crude.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.UserList, got.UserList)
	require.Equal(t, genesisState.UserCount, got.UserCount)
	// this line is used by starport scaffolding # genesis/test/assert
}

package crude

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"crude/x/crude/keeper"
	"crude/x/crude/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the user
	for _, elem := range genState.UserList {
		k.SetUser(ctx, elem)
	}

	// Set user count
	k.SetUserCount(ctx, genState.UserCount)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.UserList = k.GetAllUser(ctx)
	genesis.UserCount = k.GetUserCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}

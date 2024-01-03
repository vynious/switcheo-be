package keeper

import (
	"crude/x/crude/types"
)

var _ types.QueryServer = Keeper{}

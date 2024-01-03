package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"crude/x/crude/types"
)

func TestUserMsgServerCreate(t *testing.T) {
	_, srv, ctx := setupMsgServer(t)
	wctx := sdk.UnwrapSDKContext(ctx)

	creator := "A"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreateUser(wctx, &types.MsgCreateUser{Creator: creator})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}

func TestUserMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateUser
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateUser{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateUser{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateUser{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			_, err := srv.CreateUser(wctx, &types.MsgCreateUser{Creator: creator})
			require.NoError(t, err)

			_, err = srv.UpdateUser(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteUser
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteUser{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteUser{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeleteUser{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, srv, ctx := setupMsgServer(t)
			wctx := sdk.UnwrapSDKContext(ctx)

			_, err := srv.CreateUser(wctx, &types.MsgCreateUser{Creator: creator})
			require.NoError(t, err)
			_, err = srv.DeleteUser(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

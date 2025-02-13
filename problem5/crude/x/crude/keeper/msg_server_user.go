package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"crude/x/crude/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateUser(goCtx context.Context, msg *types.MsgCreateUser) (*types.MsgCreateUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var user = types.User{
		Creator:  msg.Creator,
		Name:     msg.Name,
		Email:    msg.Email,
		Username: msg.Username,
		Password: msg.Password,
		Address:  msg.Address,
	}
	id, err := k.AppendUser(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateUserResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateUser(goCtx context.Context, msg *types.MsgUpdateUser) (*types.MsgUpdateUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var user = types.User{
		Creator:  msg.Creator,
		Id:       msg.Id,
		Name:     msg.Name,
		Email:    msg.Email,
		Username: msg.Username,
		Password: msg.Password,
		Address:  msg.Address,
	}
	// Checks that the element exists
	val, found := k.GetUser(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetUser(ctx, user) // hashes the user password

	return &types.MsgUpdateUserResponse{}, nil
}

func (k msgServer) DeleteUser(goCtx context.Context, msg *types.MsgDeleteUser) (*types.MsgDeleteUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetUser(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveUser(ctx, msg.Id)

	return &types.MsgDeleteUserResponse{}, nil
}

func (k msgServer) UpdateUserPassword(goCtx context.Context, msg *types.MsgUpdateUserPassword) (*types.MsgUpdateUserPasswordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetUser(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	if err := k.ChangeUserPassword(ctx, msg.Id, msg.CurrentPassword, msg.NewPassword); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "current password incorrect")
	}
	return &types.MsgUpdateUserPasswordResponse{}, nil

}

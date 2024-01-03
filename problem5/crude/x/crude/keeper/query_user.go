package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/store/prefix"
	"crude/x/crude/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) UserAll(ctx context.Context, req *types.QueryAllUserRequest) (*types.QueryAllUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var users []types.User

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	userStore := prefix.NewStore(store, types.KeyPrefix(types.UserKey))

	pageRes, err := query.Paginate(userStore, req.Pagination, func(key []byte, value []byte) error {
		var user types.User
		if err := k.cdc.Unmarshal(value, &user); err != nil {
			return err
		}

		users = append(users, user)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fmt.Println("retrieving all users")
	return &types.QueryAllUserResponse{User: users, Pagination: pageRes}, nil
}

func (k Keeper) User(ctx context.Context, req *types.QueryGetUserRequest) (*types.QueryGetUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	user, found := k.GetUser(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetUserResponse{User: user}, nil
}

func (k Keeper) UserAllByAddress(ctx context.Context, req *types.QueryAllUserAddressRequest) (*types.QueryAllUserAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var users []types.User

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	userStore := prefix.NewStore(store, types.KeyPrefix(types.UserKey))
	address := req.Address
	_, err := query.Paginate(userStore, req.Pagination, func(key []byte, value []byte) error {
		var user types.User
		if err := k.cdc.Unmarshal(value, &user); err != nil {
			return err
		}
		if user.Address == address {
			users = append(users, user)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fmt.Println("retrieving all users")
	return &types.QueryAllUserAddressResponse{User: users}, nil
}

package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"crude/x/crude/types"
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	"golang.org/x/crypto/bcrypt"
)

// GetUserCount get the total number of user
func (k Keeper) GetUserCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.UserCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetUserCount set the total number of user
func (k Keeper) SetUserCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.UserCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendUser appends a user in the store with a new id and update the count
func (k Keeper) AppendUser(
	ctx context.Context,
	user types.User,
) (uint64, error) {
	// Create the user
	count := k.GetUserCount(ctx)

	// Set the ID of the appended value
	user.Id = count + 1

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.UserKey))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return count, fmt.Errorf("unable to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	appendedValue := k.cdc.MustMarshal(&user)
	store.Set(GetUserIDBytes(user.Id), appendedValue)

	// Update user count
	k.SetUserCount(ctx, count+1)

	return count, nil
}

// ChangeUserPassword takes in old password for verification and new password to updating user's password
func (k Keeper) ChangeUserPassword(ctx context.Context, id uint64, currentPassword string, newPassword string) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.UserKey))

	b := store.Get(GetUserIDBytes(id))
	if b == nil {
		return fmt.Errorf("id does not exist")
	}

	shell := types.User{}
	k.cdc.MustUnmarshal(b, &shell)

	// compare current password and inputted password if its the same
	if err := bcrypt.CompareHashAndPassword([]byte(shell.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}

	// hash the new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password: %w", err)
	}

	shell.Password = string(newHashedPassword)

	b = k.cdc.MustMarshal(&shell)
	store.Set(GetUserIDBytes(shell.Id), b)
	return nil
}

// SetUser set a specific user in the store
func (k Keeper) SetUser(ctx context.Context, user types.User) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.UserKey))

	// hashing of user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	b := k.cdc.MustMarshal(&user)
	store.Set(GetUserIDBytes(user.Id), b)
	return nil
}

// GetUser returns a user from its id
func (k Keeper) GetUser(ctx context.Context, id uint64) (val types.User, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.UserKey))
	b := store.Get(GetUserIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUser removes a user from the store
func (k Keeper) RemoveUser(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.UserKey))
	store.Delete(GetUserIDBytes(id))
}

// GetAllUser returns all user
func (k Keeper) GetAllUser(ctx context.Context) (list []types.User) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.UserKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.User
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetUserIDBytes returns the byte representation of the ID
func GetUserIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.UserKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}

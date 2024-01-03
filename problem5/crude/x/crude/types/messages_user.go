package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateUser{}

func NewMsgCreateUser(creator string, name string, email string, username string, password string, address string) *MsgCreateUser {
	return &MsgCreateUser{
		Creator:  creator,
		Name:     name,
		Email:    email,
		Username: username,
		Password: password,
		Address:  address,
	}
}

func (msg *MsgCreateUser) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateUser{}

func NewMsgUpdateUser(creator string, id uint64, name string, email string, username string, password string, address string) *MsgUpdateUser {
	return &MsgUpdateUser{
		Id:       id,
		Creator:  creator,
		Name:     name,
		Email:    email,
		Username: username,
		Password: password,
		Address:  address,
	}
}

func (msg *MsgUpdateUser) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteUser{}

func NewMsgDeleteUser(creator string, id uint64) *MsgDeleteUser {
	return &MsgDeleteUser{
		Id:      id,
		Creator: creator,
	}
}

func (msg *MsgDeleteUser) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

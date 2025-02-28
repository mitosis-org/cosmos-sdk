package keeper

import (
	"context"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// NOTE: Mitosis x/evmvalidator module doesn't support hooks.
// Instead, the x/evmvalidator calls the functions below directly.

// AfterValidatorBonded updates the signing info start height or create a new signing info
func (k Keeper) AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	signingInfo, err := k.GetValidatorSigningInfo(ctx, consAddr)
	if err == nil {
		signingInfo.StartHeight = sdkCtx.BlockHeight()
	} else {
		signingInfo = types.NewValidatorSigningInfo(
			consAddr,
			sdkCtx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
	}

	return k.SetValidatorSigningInfo(ctx, consAddr, signingInfo)
}

// AfterValidatorRemoved deletes the address-pubkey relation when a validator is removed,
func (k Keeper) AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress) error {
	return k.deleteAddrPubkeyRelation(ctx, cryptotypes.Address(consAddr))
}

// AfterValidatorCreated adds the address-pubkey relation when a validator is created.
func (k Keeper) AfterValidatorCreated(ctx context.Context, consPubKey cryptotypes.PubKey) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.AddPubkey(sdkCtx, consPubKey)
}

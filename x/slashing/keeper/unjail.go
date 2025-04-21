package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// Unjail calls the staking Unjail function to unjail a validator if the
// jailed period has concluded
func (k Keeper) Unjail(ctx context.Context, validatorAddr sdk.ValAddress) error {
	// NOTE: Use UnjailFromConsAddr instead.
	return sdkerrors.ErrNotSupported
}

// UnjailFromConsAddr calls the staking Unjail function to unjail a validator if the
// jailed period has concluded
func (k Keeper) UnjailFromConsAddr(ctx context.Context, consAddr sdk.ConsAddress) error {
	// If the validator has a ValidatorSigningInfo object that signals that the
	// validator was bonded and so we must check that the validator is not tombstoned
	// and can be unjailed at the current block.
	//
	// A validator that is jailed but has no ValidatorSigningInfo object signals
	// that the validator was never bonded and must've been jailed due to falling
	// below their minimum self-delegation. The validator can unjail at any point
	// assuming they've now bonded above their minimum self-delegation.
	info, err := k.GetValidatorSigningInfo(ctx, consAddr)
	if err == nil {
		// cannot be unjailed if tombstoned
		if info.Tombstoned {
			return types.ErrValidatorJailed
		}

		// cannot be unjailed until out of jail
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if sdkCtx.BlockHeader().Time.Before(info.JailedUntil) {
			return types.ErrValidatorJailed
		}
	}

	return k.sk.Unjail(ctx, consAddr)
}

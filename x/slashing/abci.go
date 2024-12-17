package slashing

import (
	"context"

	"cosmossdk.io/core/comet"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)

	k.Logger(sdk.UnwrapSDKContext(ctx)).Info("x/slashing: BeginBlocker")

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, voteInfo := range sdkCtx.VoteInfos() {
		err := k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, comet.BlockIDFlag(voteInfo.BlockIdFlag))

		if err == stakingtypes.ErrNoValidatorFound {
			k.Logger(sdkCtx).Error("@@@@@ VALIDATOR NOT FOUND BUT IGNORE IT.", "voteInfo.Validator.Address", voteInfo.Validator.Address, "error", err)
			continue
		}

		if err != nil {
			k.Logger(sdkCtx).Error("Error in HandleValidatorSignature", "error", err)
			return err
		}
	}

	k.Logger(sdkCtx).Info("x/slashing: BeginBlocker finished")
	return nil
}

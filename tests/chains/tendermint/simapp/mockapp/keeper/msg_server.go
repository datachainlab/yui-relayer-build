package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger-labs/yui-relayer/tests/chains/tendermint/simapp/mockapp/types"
)

var _ types.MsgServer = Keeper{}

// Transfer defines a rpc handler method for MsgTransfer.
func (k Keeper) SendPacket(goCtx context.Context, msg *types.MsgSendPacket) (*types.MsgSendPacketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.sendPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.Message, sender, msg.TimeoutHeight, msg.TimeoutTimestamp); err != nil {
		return nil, err
	}

	return &types.MsgSendPacketResponse{}, nil
}

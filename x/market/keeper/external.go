package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ovrclk/akash/x/escrow/types"
	etypes "github.com/ovrclk/akash/x/escrow/types"
)

type EscrowKeeper interface {
	AccountClose(ctx sdk.Context, id etypes.AccountID) error
	PaymentClose(ctx sdk.Context, id types.AccountID, pid string) error
}

package state

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/ovrclk/akash/app"
	"github.com/ovrclk/akash/testutil"
	dkeeper "github.com/ovrclk/akash/x/deployment/keeper"
	dtypes "github.com/ovrclk/akash/x/deployment/types"
	ekeeper "github.com/ovrclk/akash/x/escrow/keeper"
	etypes "github.com/ovrclk/akash/x/escrow/types"
	"github.com/ovrclk/akash/x/market/keeper"
	"github.com/ovrclk/akash/x/market/types"
	mtypes "github.com/ovrclk/akash/x/market/types"
	pkeeper "github.com/ovrclk/akash/x/provider/keeper"
	ptypes "github.com/ovrclk/akash/x/provider/types"
)

// TestSuite encapsulates a functional Akash nodes data stores for
// ephemeral testing.
type TestSuite struct {
	t       testing.TB
	ms      sdk.CommitMultiStore
	ctx     sdk.Context
	ekeeper ekeeper.Keeper
	mkeeper keeper.Keeper
	dkeeper dkeeper.Keeper
	pkeeper pkeeper.Keeper
	bkeeper bankkeeper.Keeper
}

// SetupTestSuite provides toolkit for accessing stores and keepers
// for complex data interactions.
func SetupTestSuite(t testing.TB, codec codec.Marshaler) *TestSuite {
	suite := &TestSuite{
		t: t,
	}

	eKey := sdk.NewKVStoreKey(etypes.StoreKey)
	mKey := sdk.NewKVStoreKey(types.StoreKey)
	dKey := sdk.NewKVStoreKey(dtypes.StoreKey)
	pKey := sdk.NewKVStoreKey(ptypes.StoreKey)

	db := dbm.NewMemDB()
	suite.ms = store.NewCommitMultiStore(db)
	suite.ms.MountStoreWithDB(eKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(mKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(dKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(pKey, sdk.StoreTypeIAVL, db)

	err := suite.ms.LoadLatestVersion()
	require.NoError(t, err)
	suite.ctx = sdk.NewContext(suite.ms, tmproto.Header{}, true, testutil.Logger(t))

	newapp := app.Setup(false)

	// XXX: unused?
	suite.ekeeper = ekeeper.NewKeeper(codec, eKey, suite.bkeeper)
	suite.mkeeper = keeper.NewKeeper(codec, mKey, newapp.GetSubspace(mtypes.ModuleName), suite.ekeeper)
	suite.dkeeper = dkeeper.NewKeeper(codec, dKey, newapp.GetSubspace(dtypes.ModuleName))
	suite.pkeeper = pkeeper.NewKeeper(codec, pKey)

	return suite
}

// SetBlockHeight provides arbitrarily setting the chain's block height.
func (ts *TestSuite) SetBlockHeight(height int64) {
	ts.ctx = ts.ctx.WithBlockHeight(height)
}

// Store provides access to the underlying KVStore
func (ts *TestSuite) Store() sdk.CommitMultiStore {
	return ts.ms
}

// Context of the current mempool
func (ts *TestSuite) Context() sdk.Context {
	return ts.ctx
}

// MarketKeeper key store
func (ts *TestSuite) MarketKeeper() keeper.Keeper {
	return ts.mkeeper
}

// DeploymentKeeper key store
func (ts *TestSuite) DeploymentKeeper() dkeeper.Keeper {
	return ts.dkeeper
}

// ProviderKeeper key store
func (ts *TestSuite) ProviderKeeper() pkeeper.Keeper {
	return ts.pkeeper
}

// BankKeeper key store
func (ts *TestSuite) BankKeeper() bankkeeper.Keeper {
	return ts.bkeeper
}

package wallet_hd

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"os"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func inmemStorage() *in_memory.InMemStore {
	return in_memory.NewInMemStore()
}

func key(seed []byte, relativePath string, storage core.Storage) (*core.DerivableKey,error) {
	if err := e2types.InitBLS(); err != nil {
		os.Exit(1)
	}

	key,err := core.BaseKeyFromSeed(seed, storage)
	if err != nil {
		return nil,err
	}

	if len(relativePath) >0 {
		return key.Derive(relativePath)
	} else {
		return key,nil
	}
}

func TestCreateAccounts(t *testing.T) {
	tests := []struct{
		testName string
		newAccounttName string
		expectedErr string
	} {
		{
			testName: "Add new account",
			newAccounttName: "account1",
			expectedErr:"",
		},
		{
			testName: "Add duplicate account, should error",
			newAccounttName: "account1",
			expectedErr:"account account1 already exists",
		},
		{
			testName: "Add account with no name, should error",
			newAccounttName: "account1",
			expectedErr:"account name is empty",
		},
	}

	// create key and wallet
	storage := inmemStorage()
	seed :=  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	err := storage.SecurelySavePortfolioSeed(seed)
	require.NoError(t,err)
	key,err := key(seed,"",storage)
	if err != nil {
		t.Error(err)
		return
	}
	w := &HDWallet{
		name:"wallet",
		id:uuid.New(),
		indexMapper:make(map[string]uuid.UUID),
		key:key,
		context:&core.PortfolioContext{
			Storage:     storage,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_,err := w.CreateValidatorAccount(test.newAccounttName)
			if test.expectedErr != "" {
				require.Errorf(t,err,test.expectedErr)
			} else {
				require.NoError(t,err)
			}
		})
	}
}

func TestWalletMarshaling(t *testing.T) {
	tests := []struct{
		id uuid.UUID
		testName string
		walletType core.WalletType
		name string
		indexMapper map[string]uuid.UUID
		seed []byte
		path string
	}{
		{
			testName:"simple wallet, no account",
			id:uuid.New(),
			walletType:core.HDWallet,
			name: "wallet1",
			indexMapper:map[string]uuid.UUID{},
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
		{
			testName:"simple wallet with accounts",
			id:uuid.New(),
			walletType:core.HDWallet,
			name: "wallet1",
			indexMapper:map[string]uuid.UUID{
				"account1" : uuid.New(),
			},
			seed:  _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"),
			path: "/0/0",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// setup storage
			storage := inmemStorage()
			err := storage.SecurelySavePortfolioSeed(test.seed)
			if err != nil {
				t.Error(err)
				return
			}

			// create key and wallet
			key,err := key(test.seed,test.path,storage)
			if err != nil {
				t.Error(err)
				return
			}
			w := &HDWallet{
				walletType:test.walletType,
				name:test.name,
				id:test.id,
				indexMapper:test.indexMapper,
				key:key,
			}

			// marshal
			byts,err := json.Marshal(w)
			if err != nil {
				t.Error(err)
				return
			}
			//unmarshal
			w1 := &HDWallet{context:&core.PortfolioContext{Storage:storage}}
			err = json.Unmarshal(byts,w1)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t,w.id,w1.id)
			require.Equal(t,w.name,w1.name)
			require.Equal(t,w.walletType,w1.walletType)
			for k := range w.indexMapper {
				v := w.indexMapper[k]
				require.Equal(t,v,w1.indexMapper[k])
			}
			require.Equal(t,w.key.PublicKey().Marshal(),w1.key.PublicKey().Marshal())
		})
	}

}
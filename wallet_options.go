package KeyVault

import (
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

type WalletOptions struct {
	encryptor wtypes.Encryptor
	password []byte
	name string
	store wtypes.Store
	enableSimpleSigner bool
	seed []byte
}

func (options *WalletOptions)SetEncryptor(encryptor wtypes.Encryptor) *WalletOptions {
	options.encryptor = encryptor
	return options
}

func (options *WalletOptions)SetStore(store wtypes.Store) *WalletOptions {
	options.store = store
	return options
}

func (options *WalletOptions)SetWalletName(name string) *WalletOptions {
	options.name = name
	return options
}

func (options *WalletOptions)SetWalletPassword(password string) *WalletOptions {
	options.password = []byte(password)
	return options
}

func (options *WalletOptions)EnableSimpleSigner(val bool) *WalletOptions {
	options.enableSimpleSigner = true
	return options
}

func (options *WalletOptions)SetSeed(seed []byte) *WalletOptions {
	options.seed = seed
	return options
}

//func (options *WalletOptions) GenerateSeed() error {
//	seed := make([]byte, 32)
//	_, err := rand.Read(seed)
//
//	options.SetSeed(seed)
//
//	return err
//}
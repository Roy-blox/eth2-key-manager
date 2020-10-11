package handler

import (
	"encoding/hex"
	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	rootcmd "github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/bloxapp/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

// Account DepositData generates account deposit-data and prints it.
func (h *Account) DepositData(cmd *cobra.Command, args []string) error {
	err := types.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	// Get index flag.
	indexFlagValue, err := flag.GetIndexFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the index flag value")
	}

	// Get seed flag.
	seedFlagValue, err := flag.GetSeedFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the seed flag value")
	}

	seedBytes, err := hex.DecodeString(seedFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to HEX decode seed")
	}

	// Get network flag
	network, err := rootcmd.GetNetworkFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the network flag value")
	}

	// Get public key flag.
	publicKeyFlagValue, err := flag.GetPublicKeyFlagValue(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve the public key flag value")
	}

	// TODO get rid of network
	store := in_memory.NewInMemStore(network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	_, err = eth2keymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	_, err = wallet.CreateValidatorAccount(seedBytes, &indexFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to create validator account")
	}

	account, err := wallet.AccountByPublicKey(publicKeyFlagValue)
	if err != nil {
		return errors.Wrap(err, "failed to get account by public key")
	}

	depositData, err := account.GetDepositData()
	if err != nil {
		return errors.Wrap(err, "failed to get deposit data")
	}

	err = h.printer.JSON(depositData)
	if err != nil {
		return errors.Wrap(err, "failed to print deposit-data JSON")
	}
	return nil
}

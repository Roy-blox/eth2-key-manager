package flag

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bloxapp/eth2-key-manager/cli/util/cliflag"
)

// ResponseType represents the network.
type ResponseType string

// Available response types.
const (
	// StorageResponseType represents the storage response type.
	StorageResponseType ResponseType = "storage"

	// ObjectResponseType represents the storage response type.
	ObjectResponseType ResponseType = "object"
)

// ResponseTypeFromString returns response type from the given string value
func ResponseTypeFromString(n string) ResponseType {
	switch n {
	case string(StorageResponseType):
		return StorageResponseType
	case string(ObjectResponseType):
		return ObjectResponseType
	default:
		panic(fmt.Sprintf("undefined response type %s", n))
	}
}

// Flag names.
const (
	indexFlag        = "index"
	seedFlag         = "seed"
	accumulateFlag   = "accumulate"
	responseTypeFlag = "response-type"
)

// AddIndexFlag adds the index flag to the command
func AddIndexFlag(c *cobra.Command) {
	cliflag.AddPersistentIntFlag(c, indexFlag, 0, "account index", true)
}

// GetIndexFlagValue gets the index flag from the command
func GetIndexFlagValue(c *cobra.Command) (int, error) {
	return c.Flags().GetInt(indexFlag)
}

// AddAccumulateFlag adds the accumulate flag to the command
func AddAccumulateFlag(c *cobra.Command) {
	cliflag.AddPersistentBoolFlag(c, accumulateFlag, false, "accumulate accounts", false)
}

// GetAccumulateFlagValue gets the accumulate flag from the command
func GetAccumulateFlagValue(c *cobra.Command) (bool, error) {
	return c.Flags().GetBool(accumulateFlag)
}

// AddResponseTypeFlag adds the response-type flag to the command
func AddResponseTypeFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, responseTypeFlag, string(StorageResponseType), "response type", false)
}

// GetResponseTypeFlagValue gets the response-type flag from the command
func GetResponseTypeFlagValue(c *cobra.Command) (ResponseType, error) {
	responseTypeValue, err := c.Flags().GetString(responseTypeFlag)
	if err != nil {
		return "", err
	}

	return ResponseTypeFromString(responseTypeValue), nil
}

// AddSeedFlag adds the seed flag to the command
func AddSeedFlag(c *cobra.Command) {
	cliflag.AddPersistentStringFlag(c, seedFlag, "", "key-vault seed", true)
}

// GetSeedFlagValue gets the seed flag from the command
func GetSeedFlagValue(c *cobra.Command) (string, error) {
	return c.Flags().GetString(seedFlag)
}

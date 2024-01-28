package types

const (
	// module name
	ModuleName = "reward"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName

	// TransientKey is the key to access the reward transient store, that is reset
	// during the Commit phase.
	TransientKey = "transient_" + ModuleName
)

var (
	KeyPrefixTransientGasReward = []byte{}
)

package v3

const (
	Version              uint64 = 3
	SquareSizeUpperBound int    = 128
	SubtreeRootThreshold int    = 64
	TxSizeCostPerByte    uint64 = 10
	GasPerBlobByte       uint32 = 8
	MaxTxSize            int    = 2097152 // 2 MiB in bytes
	// SdkMsgTransactionCap is the maximum number of SDK messages, aside from PFBs, that a block can contain.
	SdkMsgTransactionCap = 200

	// PFBTransactionCap is maximum number of PFB messages a block can contain.
	PFBTransactionCap = 600
)

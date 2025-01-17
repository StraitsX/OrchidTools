package constants

import (
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/common"
)

var ContractAddress = struct {
	TeleporterHome common.Address
	RemoteHome     common.Address

	// XSGD
	CchainXSGD common.Address
}{
	TeleporterHome: common.HexToAddress("0x48216a3597a19c4903f8933c711237e386daf088"),
	RemoteHome:     common.HexToAddress("0xf0f57f63a964423d3cf3840bcb2b2889aae8d7a7"),
	CchainXSGD:     common.HexToAddress("0xb2f85b7ab3c2b6f62df06de6ae7d09c010a5096e"),
}

var (
	DefaultERC20RequiredGas = big.NewInt(500000)

	// blockchainID is obtained from chain info section of https://subnets.avax.network/straitsx
	STRAITSX_SUBNET_BLOCKCHAIN_ID = stringToByte32LeftPadded("EJ4DyXHe4ydhsLLMiDPsHtoq5RDqgyao6Lwb9znKhs59q4NQx")

	C_CHAIN_EVM_CHAIN_ID = big.NewInt(43114)

	BLOCKCHAIN_EXPLORER_URL = "https://snowtrace.io/tx/"
)

func stringToByte32LeftPadded(s string) [32]byte {
	id := ids.FromStringOrPanic(s)
	return [32]byte(id)
}

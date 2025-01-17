package contractcalls

import (
	"math/big"
	"teleporter/constants"
	"teleporter/contract_bindings/erc20tokenhome"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/tozd/go/errors"
)

// BridgeFromCChainToSTXSubnet function bridges XSGD from the C-Chain to the StraitsX Subnet
// amount is the amount of XSGD to be sent. Add in the 6 decimals to account for XSGD decimal places, e.g. 1 XSGD = 1000000
// recipientAddress is the address of the recipient receiving the XSGD on the STX subnet
// client is the ethclient connection to the Avalanche C-Chain RPC node
// senderPrivateKey is the private key of the account that will bridge the XSGD to STX subnet
func BridgeFromCChainToSTXSubnet(amount *big.Int, recipientAddress common.Address, client *ethclient.Client, senderPrivateKey string) (string, error) {
	// signerPrivateKey is the private key of the account that will send the transaction
	signerPrivateKeyECDSA, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		return "", errors.Errorf("Failed to create deployer: %w", err)
	}

	// initialise teleporter contract
	teleporterContract, err := erc20tokenhome.NewERC20TokenHome(constants.ContractAddress.TeleporterHome, client)
	if err != nil {
		return "", errors.Errorf("Failed to create teleporter contract: %w", err)
	}

	// initialise the sender's transaction options
	transactionOpts, err := bind.NewKeyedTransactorWithChainID(signerPrivateKeyECDSA, constants.C_CHAIN_EVM_CHAIN_ID)
	if err != nil {
		return "", errors.Errorf("Failed to create transactionOpts: %w", err)
	}

	// send the XSGD to the STX subnet
	input := erc20tokenhome.SendTokensInput{
		DestinationBlockchainID:            constants.STRAITSX_SUBNET_BLOCKCHAIN_ID, // StraitsX Subnet blockchain ID obtained from chain info section of https://subnets.avax.network/straitsx
		DestinationTokenTransferrerAddress: constants.ContractAddress.RemoteHome,    // RemoteHome contract address, also known as the Bridge XSGD contract address
		Recipient:                          recipientAddress,                        // Recipient address on the STX subnet
		PrimaryFeeTokenAddress:             constants.ContractAddress.CchainXSGD,    // C-Chain XSGD contract address
		PrimaryFee:                         big.NewInt(0),                           // No primary fee is required
		SecondaryFee:                       big.NewInt(0),                           // No secondary fee is required
		RequiredGasLimit:                   constants.DefaultERC20RequiredGas,       // Default ERC20 required gas limit
		MultiHopFallback:                   common.HexToAddress("0x0"),              // Multi-hop fallback address, set to 0x0 as it is not required
	}
	tx, err := teleporterContract.Send(transactionOpts, input, amount)
	if err != nil {
		return "", errors.Errorf("Failed to send: %v", err)
	}

	return tx.Hash().Hex(), nil
}

package contractcalls

import (
	"math/big"
	"teleporter/constants"
	"teleporter/contract_bindings/erc20token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/tozd/go/errors"
)

// ApproveTeleport function approves the teleporter contract to spend XSGD on behalf of the approver
// amount is the amount of XSGD to be sent. Add in the 6 decimals to account for XSGD decimal places, e.g. 1 XSGD = 1000000
// client is the ethclient connection to the Avalanche C-Chain RPC node
// approverPrivateKey is the private key of the account that will approve the teleporter contract to spend XSGD on behalf of the approver
func ApproveTeleporter(amount *big.Int, client *ethclient.Client, approverPrivateKey string) (string, error) {

	// convert the private key to an ECDSA private key
	approverPrivateKeyECDSA, err := crypto.HexToECDSA(approverPrivateKey)
	if err != nil {
		return "", errors.Errorf("Failed to create approver: %w", err)
	}

	// initialise xsgd contract
	xsgdContract, err := erc20token.NewErc20token(constants.ContractAddress.CchainXSGD, client)
	if err != nil {
		return "", errors.Errorf("Failed to create xsgd contract: %w", err)
	}

	// initialise the approver's transaction options
	transactionOpts, err := bind.NewKeyedTransactorWithChainID(approverPrivateKeyECDSA, constants.C_CHAIN_EVM_CHAIN_ID)
	if err != nil {
		return "", errors.Errorf("Failed to create transactionOpts: %w", err)
	}

	// approve the teleporter contract to spend XSGD on behalf of the approver
	tx, err := xsgdContract.Approve(transactionOpts, constants.ContractAddress.TeleporterHome, amount)
	if err != nil {
		return "", errors.Errorf("Failed to approve: %w", err)
	}

	return tx.Hash().Hex(), nil
}

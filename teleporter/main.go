package main

// Usage:
//   Navigate to the root directory of the project and run the following command:
//     go run main.go -n <script_name> -args <args>
//
//   To bridge XSGD:
//     1. Ensure that the sender has sufficient XSGD on Avax C-Chain.
//     2. Run the approveTeleporter script to approve the teleporter contract to spend XSGD on behalf of the approver.
//     3. Ensure that the approveTeleporter script is confirmed on Avax C-Chain.
//     4. Run the bridgeXSGD script.

//   Arguments:
//   The httpRPCNodeURL should be the URL of the Avalanche C-Chain RPC node that you are connecting to. You can get one from a node provider e.g. Alchemy, Infura etc.
//   The private key should be the private key of the account that XSGD will be drawn from. Do not include 0x prefix for the private key hex string.
//   The amount of XSGD to be sent should be in the smallest unit of XSGD (e.g. 1 XSGD = 1000000).
//   The recipient address should be the address of the recipient receiving the XSGD on the STX subnet (e.g. 0xa5fb83CEb5252187ADE7c928d08A5fc215Ec4226).

//   Example:
//    go run main.go -n approveTeleporter -args <httpRPCNodeURL>,<ApproverPrivateKey>,<AmountOfXSGD>
//    go run main.go -n bridgeXSGD -args <httpRPCNodeURL>,<SenderPrivateKey>,<AmountOfXSGD>,<RecipientAddressHex>

import (
	"flag"
	"fmt"
	"math/big"
	"strings"
	"teleporter/constants"
	"teleporter/contractcalls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	scriptName = flag.String("n", "", "Script name")
	args       = flag.String("args", "", "Arguments")
)

func main() {
	flag.Parse()

	if *scriptName == "" {
		panic("Script name is required")
	}

	// Split the arguments by comma delimiter
	argv := strings.Split(*args, ",")

	// Connect to the Avalanche C-Chain RPC node
	httpRPCNodeURL := argv[0]
	client, err := ethclient.Dial(httpRPCNodeURL)
	handleErr(err)

	switch *scriptName {
	case "approveTeleporter":
		// Check if the number of arguments is correct
		if len(argv) != 3 {
			panic("Invalid number of arguments")
		}

		// Get the private key of the account that will approve the teleporter contract to spend XSGD on behalf of the approver
		approverPrivateKey := argv[1]

		// Get the amount of XSGD to be sent
		amountStr := argv[2]

		// Convert the amount to big.Int
		amountBigInt, ok := big.NewInt(0).SetString(amountStr, 10)
		if !ok {
			panic("unable to convert amount to big.Int")
		}

		// Approve the teleporter contract to spend XSGD on behalf of the approver
		txHash, err := contractcalls.ApproveTeleporter(amountBigInt, client, approverPrivateKey)
		handleErr(err)
		fmt.Printf("Approve tx sent, tx: %s\n", getBlockchainExplorerURL(txHash))
	case "bridgeXSGD":
		// Check if the number of arguments is correct
		if len(argv) != 4 {
			panic("Invalid number of arguments")
		}

		// Get the private key of the account that XSGD will be drawn from
		senderPrivateKey := argv[1]

		// Get the amount of XSGD to be sent. This should be lower than the amount approved in the approveTeleporter script
		amountStr := argv[2]

		// Get the recipient address
		recipientAddressStr := argv[3]

		// Convert the amount to big.Int
		amountBigInt, ok := big.NewInt(0).SetString(amountStr, 10)
		if !ok {
			panic("unable to convert amount to big.Int")
		}

		// Convert the recipient address to common.Address
		recipientAddress := common.HexToAddress(recipientAddressStr)

		// Bridge XSGD from Avalanche C-Chain to STX Subnet
		txHash, err := contractcalls.BridgeFromCChainToSTXSubnet(amountBigInt, recipientAddress, client, senderPrivateKey)
		handleErr(err)
		fmt.Printf("Bridge tx sent, tx: %s\n", getBlockchainExplorerURL(txHash))
	default:
		panic("Invalid script name")
	}

	fmt.Printf("Script %s executed successfully. Please check the blockchain explorer link above to confirm that transaction has been confirmed.\n", *scriptName)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getBlockchainExplorerURL(txHash string) string {
	// Print the blockchain explorer URL for the transaction
	return fmt.Sprintf("%s%s", constants.BLOCKCHAIN_EXPLORER_URL, txHash)
}

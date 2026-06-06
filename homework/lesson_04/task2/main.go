package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	counter "task2/contract"
)

func main() {
	//  Parse flags
	contractAddr := flag.String("contract", "", "counter contract address")
	numSet := flag.Int64("number", 0, "set number")
	flag.Parse()

	if *contractAddr == "" {
		log.Fatal("contract address required: --contract <address>")
	}

	//  Load environment variables
	rpcURL := os.Getenv("ETH_RPC_URL")
	privateKeyHex := os.Getenv("SENDER_PRIVATE_KEY")

	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL env variable is required")
	}
	if privateKeyHex == "" {
		log.Fatal("SENDER_PRIVATE_KEY env variable is required")
	}

	//  Connect to Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()
	fmt.Printf("connected to: %s\n", rpcURL)

	ctx := context.Background()

	//  Load private key
	privateKey, err := crypto.HexToECDSA(trim0x(privateKeyHex))
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddr := crypto.PubkeyToAddress(*publicKey)
	fmt.Printf("sender address  : %s\n", fromAddr.Hex())
	fmt.Printf("contract address: %s\n\n", *contractAddr)

	//  Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to get chain ID: %v", err)
	}
	fmt.Printf("chain ID: %d\n\n", chainID)

	//  Bind contract
	addr := common.HexToAddress(*contractAddr)
	c, err := counter.NewCounter(addr, client)
	if err != nil {
		log.Fatalf("failed to bind contract: %v", err)
	}

	//  Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create transactor: %v", err)
	}
	auth.GasLimit = 0

	//  Get initial number
	num, err := c.Number(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("failed to get number: %v", err)
	}
	fmt.Printf("initial number: %s\n\n", num.String())

	//Call increment
	fmt.Println("calling increment()...")
	auth.Nonce, err = getPendingNonce(ctx, client, fromAddr)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}

	tx, err := c.Increment(auth)
	if err != nil {
		log.Fatalf("failed to call increment: %v", err)
	}
	fmt.Printf("tx hash: %s\n", tx.Hash().Hex())

	if err := waitForTx(ctx, client, tx); err != nil {
		log.Fatalf("increment tx failed: %v", err)
	}

	num, err = c.Number(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("failed to get number after increment: %v", err)
	}
	fmt.Printf("number after increment: %s\n\n", num.String())

	// Call setNumber
	fmt.Printf("calling setNumber(%d)...\n", *numSet)
	auth.Nonce, err = getPendingNonce(ctx, client, fromAddr)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}

	tx, err = c.SetNumber(auth, big.NewInt(*numSet))
	if err != nil {
		log.Fatalf("failed to call setNumber: %v", err)
	}
	fmt.Printf("tx hash: %s\n", tx.Hash().Hex())

	if err := waitForTx(ctx, client, tx); err != nil {
		log.Fatalf("setNumber tx failed: %v", err)
	}

	num, err = c.Number(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("failed to get number after setNumber: %v", err)
	}
	fmt.Printf("number after setNumber: %s\n", num.String())
}

// getPendingNonce returns the next pending nonce for the given address
func getPendingNonce(
	ctx context.Context,
	client *ethclient.Client,
	addr common.Address,
) (*big.Int, error) {
	nonce, err := client.PendingNonceAt(ctx, addr)
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(nonce)), nil
}

func trim0x(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}

// waitForTx polls until the transaction is confirmed or context is cancelled
func waitForTx(ctx context.Context, client *ethclient.Client, tx *types.Transaction) error {
	fmt.Printf("waiting for tx %s to be confirmed...\n", tx.Hash().Hex())

	for {
		// poll every 3 seconds
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
		}

		// fetch receipt
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			// receipt not available yet, keep polling
			fmt.Println("  pending, retrying...")
			continue
		}

		// receipt found, check status
		if receipt.Status == types.ReceiptStatusSuccessful {
			fmt.Printf("  confirmed in block %d\n", receipt.BlockNumber.Uint64())
			return nil
		}

		// receipt found but tx failed on chain
		return fmt.Errorf("tx reverted in block %d", receipt.BlockNumber.Uint64())
	}
}

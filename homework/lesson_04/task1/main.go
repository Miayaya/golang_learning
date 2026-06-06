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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//查询最新区块：go run main.go --mode query
//查询指定区块：go run main.go --mode query --number 10345678
//发送交易：go run main.go --mode transfer --to 0x2782F39e807cd7dB0f0892Eb28C93193BD3f8173 --amount 0.002

func main() {
	mode := flag.String("mode", "query", "operation mode: query or transfer")
	blockNumber := flag.Int64("number", 0, "block number to query (0 = latest)")
	toAdress := flag.String("to", "", "recipient address")
	amount := flag.Float64("amount", 0.001, "ETH amount to send")
	flag.Parse()

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()
	fmt.Printf("Connect to %s\n\n", rpcURL)

	switch *mode {
	case "query":
		err := queryBlock(ctx, client, *blockNumber)
		if err != nil {
			log.Fatalf("query block error: %v", err)
		}

	case "transfer":
		err := sendTransfer(ctx, client, *toAdress, *amount)
		if err != nil {
			log.Fatalf("send tx error: %v", err)
		}

	default:
		log.Fatalf("unknown mode: %s (use: query or transfer)", *mode)
	}

}

func queryBlock(ctx context.Context, client *ethclient.Client, blockNumber int64) error {
	var num *big.Int
	if blockNumber > 0 {
		num = big.NewInt(blockNumber)
	}
	block, err := client.BlockByNumber(ctx, num)
	if err != nil {
		return fmt.Errorf("Failed to get block: %w", &err)
	}
	//output block information
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Bolck Number    : %d \n", block.Number().Uint64())
	fmt.Printf("Block Hash      : %s \n", block.Hash().Hex())
	fmt.Printf("Parent Hash      : %s \n", block.ParentHash().Hex())
	fmt.Printf("TX Number       : %d \n", len(block.Transactions()))
	fmt.Printf("Gas Limit      : %d \n", block.GasLimit())
	fmt.Printf("Gas Used      : %d \n", block.GasUsed())
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	//Print top 5 transactions

	txs := block.Transactions()
	if len(txs) > 0 {
		fmt.Println("Transaction list (Most 5):")
		limit := 5
		if len(txs) < limit {
			limit = len(txs)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("  [%d] %s\n", i+1, txs[i].Hash().Hex())
		}
	}

	return nil

}

func sendTransfer(ctx context.Context, client *ethclient.Client, toAddress string, amountEth float64) error {
	if toAddress == "" {
		log.Fatal("missing  --to flag for transfer mode")
	}

	// 检查私钥环境变量
	privKeyHex := os.Getenv("SENDER_PRIVATE_KEY")
	if privKeyHex == "" {
		log.Fatal("SENDER_PRIVATE_KEY is not set (required for transfer mode)")
	}

	// 解析私钥
	privKey, err := crypto.HexToECDSA(trim0x(privKeyHex))
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}

	// 获取发送方地址
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
	toAddr := common.HexToAddress(toAddress)

	// 获取链 ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain id: %v", err)
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	//获取 amount
	amount := new(big.Float).SetFloat64(amountEth)
	weiFloat := new(big.Float).Mul(amount, big.NewFloat(1e18))
	weiInt, _ := weiFloat.Int(nil)

	//获取 Gas 参数
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas tip: %w", err)
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get header: %w", err)
	}

	var gasFeeCap *big.Int
	if header.BaseFee != nil {
		// EIP-1559：MaxFeePerGas = 2 * BaseFee + GasTipCap
		gasFeeCap = new(big.Int).Add(
			new(big.Int).Mul(header.BaseFee, big.NewInt(2)),
			gasTipCap,
		)
		fmt.Printf("Base Fee    : %s wei\n", header.BaseFee.String())
	} else {
		// 非 EIP-1559：用传统 gas price
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to get gas price: %w", err)
		}
		gasFeeCap = gasPrice
	}

	//检查余额是否足够
	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}
	fmt.Printf("The account balance is   : %s wei\n", balance.String())

	gasLimit, err := estimateGas(ctx, client, fromAddr, toAddr, weiInt, nil)
	if err != nil {
		// 估算失败时 fallback 到 21000（纯转账的安全值）
		fmt.Printf("[WARN] gas 估算失败，使用默认值 21000: %v\n", err)
		gasLimit = 21000
	}

	// 第九步：验证余额是否足够支付转账 + gas
	gasCost := new(big.Int).Mul(gasFeeCap, new(big.Int).SetUint64(gasLimit))
	totalCost := new(big.Int).Add(weiInt, gasCost)
	if balance.Cmp(totalCost) < 0 {
		return fmt.Errorf(
			"insufficient ETH balance for gas:need %s wei(transfer %s + gas %s),balance is %s",
			totalCost.String(),
			weiInt.String(),
			gasCost.String(),
			balance.String(),
		)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &toAddr,
		Value:     weiInt,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
	})

	//签名交易
	signer := types.NewLondonSigner(chainID)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction has been sent!\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("From          : %s\n", fromAddr.Hex())
	fmt.Printf("To            : %s\n", toAddr.Hex())
	fmt.Printf("Amount    : %f ETH (%s wei)\n", amountEth, weiInt.String())
	fmt.Printf("Gas Limit   : %d\n", gasLimit)
	fmt.Printf("Gas Tip Cap : %s wei\n", gasTipCap.String())
	fmt.Printf("Gas Fee Cap : %s wei\n", gasFeeCap.String())
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Transaction Hash    : %s\n", signedTx.Hash().Hex())
	fmt.Printf("Etherscan   : https://sepolia.etherscan.io/tx/%s\n",
		signedTx.Hash().Hex())

	return nil

}

func trim0x(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}

// estimateGas 动态估算 gas，并加上缓冲
func estimateGas(
	ctx context.Context,
	client *ethclient.Client,
	from common.Address,
	to common.Address,
	value *big.Int,
	data []byte,
) (uint64, error) {

	// 构造估算请求
	msg := ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: value,
		Data:  data,
	}

	// 向节点请求 gas 估算
	estimated, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// 加 20% 缓冲，防止边界情况失败
	// estimated * 120 / 100
	buffered := estimated * 120 / 100

	fmt.Printf("Gas estimated    : %d\n", estimated)
	fmt.Printf("Gas include buffered  : %d(+20%%)\n", buffered)

	return buffered, nil
}

package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func transferETH(client *ethclient.Client, fromPrivKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	ctx := context.Background()

	publicKey := fromPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	tx := types.NewTx(&types.LegacyTx{Nonce: nonce, GasPrice: gasPrice, Gas: gasLimit, To: &to, Value: amount})

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivKey)
	if err != nil {
		log.Fatal(err)
	}

	return client.SendTransaction(ctx, signedTx)
}

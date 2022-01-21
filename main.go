package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	to        string
	proxy     string
	clash     string
	rpcUrl    string
	socketUrl string
)

type run_single func() error

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func parseArgs() error {
	flag.StringVar(&to, "to", "0x14c75FC7aE1e566f57893435F34c7A488CBEf2e1", "your address (required)")
	flag.StringVar(&proxy, "proxy", "socks5://127.0.0.1:7890", "proxy url")
	flag.StringVar(&clash, "clash", "✈️ 手动切换", "clash api ")
	flag.StringVar(&socketUrl, "socket-url", "wss://testnet.binance.org/faucet-smart/api", "bsc faucet socket url")
	flag.StringVar(&rpcUrl, "rpc-url", "https://data-seed-prebsc-1-s1.binance.org:8545", "bsc testnet rpc url")
	flag.Parse()

	if to == "" {
		flag.Usage()
		return errors.New("`to` cannot be empty")
	}

	return nil
}

func run() error {
	if err := parseArgs(); err != nil {
		return err
	}
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return err
	}
	defer client.Close()

	if clash == "" {
		err = run_single_get()
		if err != nil {
			return err
			// fmt.Print(err)
		}
	} else {
		err = run_clash(run_single_get, clash)
		if err != nil {
			return err
			// fmt.Print(err)
		}
	}

	return nil
}

func run_single_get() error {
	ctx := context.Background()

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return err
	}
	defer client.Close()

	captcha, err := NewCaptcha(ctx)
	if err != nil {
		return fmt.Errorf("NewCaptcha: %w", err)
	}

	account, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	socketTx, err := NewSocketTransaction(
		ctx,
		socketUrl,
		captcha,
		crypto.PubkeyToAddress(account.PublicKey).Hex(),
		proxy,
	)
	if err != nil {
		return fmt.Errorf("NewSocketTransaction: %w", err)
	}

	waitCtx, waitCancel := context.WithTimeout(ctx, time.Minute)
	defer waitCancel()
	if err := WaitTx(waitCtx, client, socketTx); err != nil {
		return fmt.Errorf("WaitTx: %w", err)
	}

	sendBnbCtx, sendBnbCancel := context.WithTimeout(ctx, time.Minute)
	defer sendBnbCancel()
	tx, err := SendBnb(sendBnbCtx, client, account, to)
	if err != nil {
		return fmt.Errorf("SendBnb: %w", err)
	}
	fmt.Println("success", tx)
	return nil
}

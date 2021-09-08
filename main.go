package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/openware/binance-cli/pkg/binance"
	"github.com/openware/binance-cli/pkg/helpers"
	"github.com/openware/binance-cli/pkg/opendax"

	"github.com/openware/pkg/kli"
)

// version of the command line
var version = "SNAPSHOT"

// AutoEnabled defines whether auto mode should be used for the markets cmd
var AutoEnabled = false

func main() {
	cli := kli.NewCli("binance-cli", "Binance cli", version)

	feesCommand := kli.NewCommand("fees", "Compare fees").Action(compareFees)
	cli.DefaultCommand(feesCommand)
	cli.AddCommand(feesCommand)

	marketsCommand := kli.NewCommand("markets", "Compare markets").Action(compareMarkets)
	cli.AddCommand(marketsCommand)

	marketsCommand.BoolFlag("auto", "Automatically update every market and save the output", &AutoEnabled)

	if err := cli.Run(); err != nil {
		fmt.Printf("Error encountered: %v\n", err)
		os.Exit(1)
	}
}

func compareMarkets() error {
	config := readConfig()
	binanceClient := binance.NewBinanceClient("", "")
	binanceInfo, err := binanceClient.ExchangeInfo()
	if err != nil {
		return err
	}

	opendaxClient := opendax.NewOpendaxClient(config.PlatformBaseUrl)
	opendaxClient.Authorize(config.OpendaxApiKey, config.OpendaxApiSecret)
	opendaxMarkets, err := opendaxClient.FetchOpendaxMarkets()
	if err != nil {
		return err
	}

	var updatedMarkets []string

	for _, opendaxMarket := range opendaxMarkets {
		binanceMarket, ok := binanceInfo.MarketRegistry[opendaxMarket.ToBinanceMarketName()]
		if ok {
			tickerPrice, err := binanceClient.TickerPriceInfo(binanceMarket.Symbol)
			if err != nil {
				fmt.Printf("ERR: compareMarkets: ticker price fetch for %s failed: %s\n", binanceMarket.Symbol, err)
				continue
			}

			minAmount, err := binanceMarket.CalculateMinAmount(tickerPrice.Price)
			if err != nil {
				fmt.Printf("ERR: compareMarkets: min amount calculation for %s failed: %s\n", binanceMarket.Symbol, err)
				continue
			}

			convertedBinanceMarket := binanceMarket.ToOpendaxMarket(minAmount)
			fmt.Println("Comparing", opendaxMarket.Symbol)
			fmt.Println("Equal:", opendax.CompareOpendaxMarkets(opendaxMarket, convertedBinanceMarket))
			fmt.Println("Binance:")
			convertedBinanceMarket.Print()
			fmt.Println("Opendax:")
			opendaxMarket.Print()
			fmt.Println("")

			if opendax.CompareOpendaxMarkets(opendaxMarket, convertedBinanceMarket) {
				fmt.Println("Skipping")
				continue
			}

			var input string
			if AutoEnabled {
				fmt.Println("Skipping market update prompt due to auto mode")
			} else {
				fmt.Print("Update this market?")
				fmt.Scanln(&input)
			}

			if AutoEnabled || input == "y" {
				updatedMarket, err := opendaxClient.UpdateOpendaxMarket(opendax.UpdateMarketRequest{
					Symbol:          opendaxMarket.Symbol,
					MinPrice:        convertedBinanceMarket.MinPrice,
					MaxPrice:        convertedBinanceMarket.MaxPrice,
					MinAmount:       convertedBinanceMarket.MinAmount,
					AmountPrecision: convertedBinanceMarket.AmountPrecision,
					PricePrecision:  convertedBinanceMarket.PricePrecision,
				})

				if err != nil {
					panic(err)
				}

				fmt.Println("New market:")
				updatedMarket.Print()

				updatedMarkets = append(updatedMarkets, updatedMarket.Name)
			} else if input == "n" {
				continue
			} else {
				panic("Wrong input")
			}
		} else {
			fmt.Println(opendaxMarket.Symbol, "is missing on Binance")
		}
	}

	if AutoEnabled {
		err = helpers.WriteToFile("updated-markets.txt", fmt.Sprintf("%v", updatedMarkets))
		if err != nil {
			fmt.Printf("Error saving updated markets: %s\nUpdated markets: %v", err, updatedMarkets)
		}
	}

	fmt.Println("Total OpenDAX markets:", len(opendaxMarkets))

	return nil
}

func compareFees() error {
	config := readConfig()

	opendaxClient := opendax.NewOpendaxClient(config.PlatformBaseUrl)
	opendaxCurrencies, err := opendaxClient.FetchOpendaxCurrencies()
	if err != nil {
		return err
	}

	binanceClient := binance.NewBinanceClient(config.BinanceApiKey, config.BinanceSecret)
	binanceCurrencies, err := binanceClient.CoinsInfo()
	if err != nil {
		return err
	}

	// Save Binance Currencies info as Map to optimize search
	binanceCoinsRegistry := make(map[string]*binance.BinanceCurrency)
	for _, coin := range binanceCurrencies {
		binanceCoinsRegistry[coin.Code] = coin
	}

	for _, opendaxCurrency := range opendaxCurrencies {
		binanceCurrency := binanceCoinsRegistry[opendaxCurrency.ToBinanceCoinName()]
		if binanceCurrency == nil {
			color.Yellow(fmt.Sprintf("\n%s cannot be found on Binance, skipping ...\n", opendaxCurrency.ToBinanceCoinName()))
			continue
		}

		for _, network := range binanceCurrency.Networks {
			fmt.Printf("\n%s coin on %s network:\n", opendaxCurrency.ToBinanceCoinName(), network.Name)

			opendaxMinWithdraw, err := opendaxCurrency.MinWithdrawAmount.Float64()
			if err != nil {
				color.Magenta(fmt.Sprintf("\nERROR: %s\n%s cannot convert Opendax Min Withdraw to Float64, skipping ...\n", err, opendaxCurrency.ToBinanceCoinName()))
			}

			binanceMinWithdraw, err := network.WithdrawMin.Float64()
			if err != nil {
				color.Magenta(fmt.Sprintf("\nERROR: %s\n%s cannot convert Binance Min Withdraw to Float64, skipping ...\n", err, opendaxCurrency.ToBinanceCoinName()))
			}

			if opendaxMinWithdraw >= binanceMinWithdraw {
				color.Green(fmt.Sprintf("MinWithdraw amount satisfy condition\nOpendax: %f; Binance: %f;\n", opendaxMinWithdraw, binanceMinWithdraw))
			} else {
				color.Red(fmt.Sprintf("MinWithdraw amount DOES NOT satisfy condition!\nOpendax: %f; Binance: %f;\n", opendaxMinWithdraw, binanceMinWithdraw))
			}

			opendaxWithdrawFee, err := opendaxCurrency.WithdrawFee.Float64()
			if err != nil {
				color.Magenta(fmt.Sprintf("\nERROR: %s\n%s cannot convert Opendax Withdraw Fee to Float64, skipping ...\n", err, opendaxCurrency.ToBinanceCoinName()))
			}

			binanceWithdrawFee, err := network.WithdrawFee.Float64()
			if err != nil {
				color.Magenta(fmt.Sprintf("\nERROR: %s\n%s cannot convert Binance Withdraw Fee to Float64, skipping ...\n", err, opendaxCurrency.ToBinanceCoinName()))
			}

			if opendaxWithdrawFee >= binanceWithdrawFee {
				color.Green(fmt.Sprintf("WithdrawFee amount satisfy condition\nOpendax: %f; Binance: %f;\n", opendaxWithdrawFee, binanceWithdrawFee))
			} else {
				color.Red(fmt.Sprintf("WithdrawFee amount DOES NOT satisfy condition!\nOpendax: %f; Binance: %f;\n", opendaxWithdrawFee, binanceWithdrawFee))
			}
		}
	}

	return nil
}

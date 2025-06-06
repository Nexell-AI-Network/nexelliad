package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Nexell-AI-Network/nexelliad/v2/cmd/nexelliawallet/libkaspawallet/serialization"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/consensushashing"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/constants"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/txscript"
	"github.com/pkg/errors"
)

func parse(conf *parseConfig) error {
	if conf.Transaction == "" && conf.TransactionFile == "" {
		return errors.Errorf("Either --transaction or --transaction-file is required")
	}
	if conf.Transaction != "" && conf.TransactionFile != "" {
		return errors.Errorf("Both --transaction and --transaction-file cannot be passed at the same time")
	}

	transactionHex := conf.Transaction
	if conf.TransactionFile != "" {
		transactionHexBytes, err := ioutil.ReadFile(conf.TransactionFile)
		if err != nil {
			return errors.Wrapf(err, "Could not read hex from %s", conf.TransactionFile)
		}
		transactionHex = strings.TrimSpace(string(transactionHexBytes))
	}

	transactions, err := decodeTransactionsFromHex(transactionHex)
	if err != nil {
		return err
	}
	for i, transaction := range transactions {

		partiallySignedTransaction, err := serialization.DeserializePartiallySignedTransaction(transaction)
		if err != nil {
			return err
		}

		fmt.Printf("Transaction #%d ID: \t%s\n", i+1, consensushashing.TransactionID(partiallySignedTransaction.Tx))
		fmt.Println()

		allInputSompi := uint64(0)
		for index, input := range partiallySignedTransaction.Tx.Inputs {
			partiallySignedInput := partiallySignedTransaction.PartiallySignedInputs[index]

			if conf.Verbose {
				fmt.Printf("Input %d: \tOutpoint: %s:%d \tAmount: %.2f Kaspa\n", index, input.PreviousOutpoint.TransactionID,
					input.PreviousOutpoint.Index, float64(partiallySignedInput.PrevOutput.Value)/float64(constants.SompiPerKaspa))
			}

			allInputSompi += partiallySignedInput.PrevOutput.Value
		}
		if conf.Verbose {
			fmt.Println()
		}

		allOutputSompi := uint64(0)
		for index, output := range partiallySignedTransaction.Tx.Outputs {
			scriptPublicKeyType, scriptPublicKeyAddress, err := txscript.ExtractScriptPubKeyAddress(output.ScriptPublicKey, conf.ActiveNetParams)
			if err != nil {
				return err
			}

			addressString := scriptPublicKeyAddress.EncodeAddress()
			if scriptPublicKeyType == txscript.NonStandardTy {
				scriptPublicKeyHex := hex.EncodeToString(output.ScriptPublicKey.Script)
				addressString = fmt.Sprintf("<Non-standard transaction script public key: %s>", scriptPublicKeyHex)
			}

			fmt.Printf("Output %d: \tRecipient: %s \tAmount: %.2f Kaspa\n",
				index, addressString, float64(output.Value)/float64(constants.SompiPerKaspa))

			allOutputSompi += output.Value
		}
		fmt.Println()

		fmt.Printf("Fee:\t%d Sompi\n\n", allInputSompi-allOutputSompi)
	}

	return nil
}

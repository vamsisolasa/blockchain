package main

import (
    	"github.com/stellar/go/build"
    	"github.com/stellar/go/clients/horizon"
    	"log"
	"github.com/stellar/go/keypair"
)

func main () {
issuerSeed := "SC43AJMJIVKPVDIYPAQAIYMLD3RN4L76ZOPWQE4FURFJISNWU7J6J2MU"
recipientSeed := "SCWAO5M2HJOZKJCNYLHMDLQVPKZAUYOBUVN3DVG6NCPLTHAC42AMHBBZ"

// Keys for accounts to issue and receive the new asset
issuer, err := keypair.Parse(issuerSeed)
if err != nil { log.Fatal(err) }
recipient, err := keypair.Parse(recipientSeed)
if err != nil { log.Fatal(err) }

// Create an object to represent the new asset
astroDollar := build.CreditAsset("AstroDollar", issuer.Address())

// First, the receiving account must trust the asset
trustTx, err := build.Transaction(
    build.SourceAccount{recipient.Address()},
    build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
    build.TestNetwork,
    build.Trust(astroDollar.Code, astroDollar.Issuer, build.Limit("100.25")),
)
if err != nil { log.Fatal(err) }
trustTxe, err := trustTx.Sign(recipientSeed)
if err != nil { log.Fatal(err) }
trustTxeB64, err := trustTxe.Base64()
if err != nil { log.Fatal(err) }
_, err = horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
if err != nil { log.Fatal(err) }

// Second, the issuing account actually sends a payment using the asset
paymentTx, err := build.Transaction(
    build.SourceAccount{issuer.Address()},
    build.TestNetwork,
    build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
    build.Payment(
        build.Destination{AddressOrSeed: recipient.Address()},
        build.CreditAmount{"AstroDollar", issuer.Address(), "10"},
    ),
)
if err != nil { log.Fatal(err) }
paymentTxe, err := paymentTx.Sign(issuerSeed)
if err != nil { log.Fatal(err) }
paymentTxeB64, err := paymentTxe.Base64()
if err != nil { log.Fatal(err) }
_, err = horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
if err != nil { log.Fatal(err) }
}

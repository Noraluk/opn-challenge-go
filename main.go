package main

import (
	"challenge-go/pkg/config"
	"challenge-go/pkg/logger"
	myOmise "challenge-go/pkg/omise"
	"challenge-go/services"
	"os"

	"github.com/omise/omise-go"
)

func main() {
	if len(os.Args) <= 1 || (len(os.Args) > 1 && len(os.Args[1]) == 0) {
		panic("please enter file path")
	}
	filePath := os.Args[1]

	err := config.Init()
	if err != nil {
		panic(err)
	}

	conf := config.GetConfig()
	err = logger.Init(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	client, err := omise.NewClient(conf.Omise.PublicKey, conf.Omise.SecretKey)
	if err != nil {
		panic(err)
	}

	omiseClient := myOmise.NewOmiseClient(client)

	donationService := services.NewDonationService(omiseClient)
	err = donationService.MakeDonations(filePath, conf.Currency)
	if err != nil {
		panic(err)
	}
}

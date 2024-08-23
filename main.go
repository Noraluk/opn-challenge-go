package main

import (
	"challenge-go/pkg/config"
	"challenge-go/pkg/logger"
	myOmise "challenge-go/pkg/omise"
	"challenge-go/services"
	"flag"

	"github.com/omise/omise-go"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file_path", "", "path of input csv file")
	var currency string
	flag.StringVar(&currency, "currency", "thb", "currency unit")
	var logLevel string
	flag.StringVar(&logLevel, "log_level", "info", "level of log")
	flag.Parse()

	err := config.Init()
	if err != nil {
		panic(err)
	}

	err = logger.Init(logLevel)
	if err != nil {
		panic(err)
	}

	conf := config.GetConfig()
	client, err := omise.NewClient(conf.Omise.PublicKey, conf.Omise.SecretKey)
	if err != nil {
		panic(err)
	}

	omiseClient := myOmise.NewOmiseClient(client)

	donationService := services.NewDonationService(omiseClient)
	err = donationService.MakeDonations(filePath, currency)
	if err != nil {
		panic(err)
	}
}

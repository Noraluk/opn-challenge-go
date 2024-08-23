package services

import (
	"challenge-go/cipher"
	"challenge-go/models"
	"challenge-go/pkg/logger"
	myOmise "challenge-go/pkg/omise"
	"container/heap"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type DonationService interface {
	MakeDonations(filePath string, currency string) error
}

type donationService struct {
	omiseClient myOmise.OmiseClient
	logger      logger.Logger
}

func NewDonationService(omiseClient myOmise.OmiseClient) DonationService {
	return &donationService{
		omiseClient: omiseClient,
		logger:      logger.WithPrefix("service/donation"),
	}
}

func (s donationService) MakeDonations(filePath string, currency string) error {
	fmt.Println("performing donations...")

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	rot128Reader, err := cipher.NewRot128Reader(f)
	if err != nil {
		return err
	}

	csvReader := csv.NewReader(rot128Reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	records = records[1:]

	successfulCustomerCh := make(chan models.Customer)
	failedCustomerCh := make(chan models.Customer)
	quit := make(chan struct{})

	readWg := &sync.WaitGroup{}
	readWg.Add(1)
	go func() {
		defer readWg.Done()

		var successfulDonation int64
		var faultyDonation int64

		customerHeap := &models.CustomerHeap{}
		heap.Init(customerHeap)

		for {
			select {
			case cust := <-successfulCustomerCh:
				s.logger.Wrap("successful customer : %v", cust).Debug()
				successfulDonation += cust.Amount
				if customerHeap.Len() < 3 {
					heap.Push(customerHeap, cust)
				} else if cust.Amount > (*customerHeap)[0].Amount {
					heap.Pop(customerHeap)
					heap.Push(customerHeap, cust)
				}
			case cust := <-failedCustomerCh:
				s.logger.Wrap("failed customer : %v", cust).Debug()
				faultyDonation += cust.Amount
			case <-quit:
				fmt.Printf("done.\n\n")
				s.showSummary(successfulDonation, faultyDonation, customerHeap, len(records))
				return
			}
		}
	}()

	wg := new(sync.WaitGroup)
	wgCh := make(chan struct{}, 5)
	for _, record := range records {
		wg.Add(1)
		wgCh <- struct{}{}

		go func() {
			var err error
			customer := models.Customer{}
			err = customer.SetCustomer(record)
			if err != nil {
				s.logger.Wrap("set customer, got err: %v", err).Debug()
				return
			}

			defer func() {
				if err != nil {
					failedCustomerCh <- customer
				}
				<-wgCh
				wg.Done()
			}()

			card, err := s.createCard(customer)
			if err != nil {
				s.logger.Wrap("create card, got err: %v", err).Debug()
				return
			}

			_, err = s.createCharge(customer.Amount, currency, card.ID)
			if err != nil {
				s.logger.Wrap("create charge, got err: %v", err).Debug()
				return
			}

			successfulCustomerCh <- customer
		}()
	}

	wg.Wait()

	close(quit)
	readWg.Wait()

	return nil
}

func (s donationService) createCard(customer models.Customer) (*omise.Card, error) {
	card := &omise.Card{}
	err := s.omiseClient.CreateToken(card, &operations.CreateToken{
		Name:            customer.Name,
		Number:          customer.CCNumber,
		ExpirationMonth: time.Month(customer.ExpiredMonth),
		ExpirationYear:  customer.ExpiredYear + 10, // it needs to be plus 10 because input year is older than the current year
		SecurityCode:    customer.CVV,
	})
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (s donationService) createCharge(amount int64, currency string, cardToken string) (*omise.Charge, error) {
	charge := &omise.Charge{}
	err := s.omiseClient.CreateCharge(charge, &operations.CreateCharge{
		Amount:   amount,
		Currency: currency,
		Card:     cardToken,
	})
	if err != nil {
		return nil, err
	}
	return charge, nil
}

func (s donationService) showSummary(successfulDonation int64, faultyDonation int64, customerHeap *models.CustomerHeap, customerCount int) {
	total := int(successfulDonation + faultyDonation)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRows([]table.Row{
		{"total received: THB", humanize.FormatInteger("#,###.##", total)},
		{"successfully donated: THB", humanize.FormatInteger("#,###.##", int(successfulDonation))},
		{"faulty donation: THB", humanize.FormatInteger("#,###.##", int(faultyDonation))},
		{""},
		{"average per person: THB", humanize.FormatFloat("#,###.##", float64(total)/float64(customerCount))},
	})

	sortedCustomers := make([]models.Customer, customerHeap.Len())
	for i := range sortedCustomers {
		sortedCustomers[i] = heap.Pop(customerHeap).(models.Customer)
	}

	for i := len(sortedCustomers) - 1; i >= 0; i-- {
		if i == len(sortedCustomers)-1 {
			t.AppendRow(table.Row{"top donors:    ", sortedCustomers[i].Name})
		} else {
			t.AppendRow(table.Row{"", sortedCustomers[i].Name})
		}
	}

	t.SetStyle(table.Style{
		Box: table.StyleBoxDefault,
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 1,
			Align:  text.AlignRight,
		},
		{
			Number: 3,
			Align:  text.AlignRight,
		},
	})
	t.Render()
}

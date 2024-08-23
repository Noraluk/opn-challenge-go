package models

import (
	"strconv"
)

type Customer struct {
	Name         string
	Amount       int64
	CCNumber     string
	CVV          string
	ExpiredMonth int
	ExpiredYear  int
}

func (m *Customer) SetCustomer(record []string) error {
	var err error

	m.Name = record[0]
	m.Amount, err = strconv.ParseInt(record[1], 10, 64)
	if err != nil {
		return err
	}

	m.CCNumber = record[2]
	m.CVV = record[3]
	m.ExpiredMonth, err = strconv.Atoi(record[4])
	if err != nil {
		return err
	}
	m.ExpiredYear, err = strconv.Atoi(record[5])
	if err != nil {
		return err
	}

	return nil
}

type CustomerHeap []Customer

func (h CustomerHeap) Len() int           { return len(h) }
func (h CustomerHeap) Less(i, j int) bool { return h[i].Amount < h[j].Amount } // Min-heap
func (h CustomerHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CustomerHeap) Push(x interface{}) {
	*h = append(*h, x.(Customer))
}

func (h *CustomerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

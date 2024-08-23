package omise

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmiseClient interface {
	CreateToken(result interface{}, operation *operations.CreateToken) error
	CreateCharge(result interface{}, operation *operations.CreateCharge) error
}

type omiseClient struct {
	*omise.Client
}

func NewOmiseClient(client *omise.Client) OmiseClient {
	return &omiseClient{
		Client: client,
	}
}

func (o omiseClient) CreateToken(result interface{}, operation *operations.CreateToken) error {
	return o.Client.Do(result, operation)
}

func (o omiseClient) CreateCharge(result interface{}, operation *operations.CreateCharge) error {
	return o.Client.Do(result, operation)
}

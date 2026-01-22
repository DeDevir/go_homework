package v1

import paymentV1 "github.com/DeDevir/go_homework/shared/pkg/proto/payment/v1"

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}

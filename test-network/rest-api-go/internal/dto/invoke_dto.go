package dto

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type InvokeRequest struct {
	ChaincodeID string   `json:"chaincodeid"`
	ChannelID   string   `json:"channelid"`
	Function    string   `json:"function"`
	Args        []string `json:"args"`
}

type InvokeSuccessResponse = SuccessResponse[*client.Status]

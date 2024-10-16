package domain

type InvokeRequest struct {
	ChaincodeID string   `json:"chaincodeid"`
	ChannelID   string   `json:"channelid"`
	Function    string   `json:"function"`
	Args        []string `json:"args"`
}

type InvokeSuccessResponse = SuccessResponse[interface{}]

type InvokeInterface interface {
	ExecuteInvoke(channelID, chaincodeName, function string, args []string) (*TransactionProposalStatus, error)
}

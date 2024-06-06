package dto

type InvokeRequest struct {
	ChaincodeID string   `json:"chaincodeid"`
	ChannelID   string   `json:"channelid"`
	Function    string   `json:"function"`
	Args        []string `json:"args"`
}

type InvokeResponse struct {
	TransactionID string `json:"transaction_id"`
	Result        string `json:"result"`
}

package dto

type InvokeRequest struct {
	ChaincodeID string   `json:"chaincodeid"`
	ChannelID   string   `json:"channelid"`
	Function    string   `json:"function"`
	Args        []string `json:"args"`
}

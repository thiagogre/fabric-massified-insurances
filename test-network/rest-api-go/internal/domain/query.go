package domain

type QueryInterface interface {
	ExecuteQuery(channelID, chainCodeName, function string, args []string) (interface{}, error)
}

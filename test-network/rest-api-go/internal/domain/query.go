package domain

type QuerySuccessResponse = SuccessResponse[interface{}]

type QueryInterface interface {
	ExecuteQuery(channelID, chainCodeName, function string, args []string) (interface{}, error)
}

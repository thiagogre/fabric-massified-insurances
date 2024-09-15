package domain

type InvokeInterface interface {
	ExecuteInvoke(channelID, chaincodeName, function string, args []string) (*TransactionProposalStatus, error)
}

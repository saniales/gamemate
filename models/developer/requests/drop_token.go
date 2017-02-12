package developerResponses

//DropToken represents a POSITIVE response to a developerRequests.DropToken.
//
//If the response is NEGATIVE, please refer to error
type DropToken struct {
	result bool
}

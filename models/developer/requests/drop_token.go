package developerResponses

//DropToken represents a POSITIVE response to a developerRequests.DropToken.
//
//If the response is NEGATIVE, please refer to error
type DropToken struct {
	result string
}

//CreateResponse adds the "OK" message string as result, implying this is
//a POSITIVE response.
func (receiver *DropToken) CreateResponse() {
	result = "OK"
}

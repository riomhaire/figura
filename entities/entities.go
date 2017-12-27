package entities

const NoError = 0
const UnknownApplication = 1
const AuthenticationError = 2
const NotImplimentedError = 3

// ApplicationConfiguration - contains application information and any errors
type ApplicationConfiguration struct {
	ResultType int
	Message    string
	Data       []byte
}

package restik

type Reply interface {
	New() Reply
	SetResponse(interface{})
	SetError(error)
	GetError() error
}

type serveReply struct {
	Response interface{} `json:"response,omitempty"`
	Error    Error       `json:"error,omitempty"`
}

func (sr *serveReply) New() Reply {
	return &serveReply{}
}

func (sr *serveReply) SetResponse(resp interface{}) {
	sr.Response = resp
}

func (sr *serveReply) SetError(err error) {
	sr.Error = FromAnotherError(err)
}

func (sr *serveReply) GetError() error {
	return sr.Error
}

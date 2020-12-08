package restik

type Reply interface {
	New() Reply
	SetResponse(interface{})
	SetError(error)
	GetError() error
	HasError() bool
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

func (sr *serveReply) HasError() bool {
	return sr.Error != nil
}

// func writeReply(w http.ResponseWriter, sr Reply) (int, error) {
// 	b, err := json.Marshal(sr)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return 0, err
// 	}
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	if sr.HasError() {
// 		err := FromAnotherError(sr.GetError())
// 		w.WriteHeader(err.GetStatus())
// 	} else {
// 		w.WriteHeader(http.StatusOK)
// 	}
// 	return w.Write(b)
// }

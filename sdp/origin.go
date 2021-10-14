package sdp

import (
	"fmt"
)

type SDPOrigin struct {
	Username       string
	SessionID      string
	SessionVersion string
	Address        string
}

func (origin *SDPOrigin) Init() {
	origin.Username = "-"
	origin.SessionID = "0"
	origin.SessionVersion = "0"
	origin.Address = "0.0.0.0"
}

func (origin *SDPOrigin) Println() string {
	return fmt.Sprintf("o=%s %s %s IN IP4 %s", origin.Username, origin.SessionID, origin.SessionVersion, origin.Address)
}

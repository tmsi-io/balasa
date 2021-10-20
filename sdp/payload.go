package sdp

import (
	"bytes"
	"fmt"
)

type Payload struct {
	Version     string
	Origin      *Origin
	SessionName string
	Info        string
	URI         string
	Email       string
	Phone       string
	Connection  string //IN IP4 224.2.17.12/127
	Medias      []*Media
	Attrs       []string
	Times       []Time
}

func (sdp *Payload) AddAttribute(attr string) {
	sdp.Attrs = append(sdp.Attrs, attr)
}

func (sdp *Payload) AddMedia(media *Media) {
	sdp.Medias = append(sdp.Medias, media)
}

func (sdp *Payload) Encode() []byte {
	buff := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buff, "v=%s\r\n", sdp.Version)
	_, _ = fmt.Fprintf(buff, "%s\r\n", sdp.Origin.Println())
	_, _ = fmt.Fprintf(buff, "s=%s\r\n", sdp.SessionName)
	_, _ = fmt.Fprintf(buff, "c=%s\r\n", sdp.Connection)
	_, _ = fmt.Fprintf(buff, "t=0 0\r\n")
	_, _ = fmt.Fprintf(buff, "a=control:*\r\n")
	_, _ = fmt.Fprintf(buff, "a=range:npt=now-\r\n")
	for _, media := range sdp.Medias {
		_, _ = fmt.Fprintf(buff, "%s\r\n", media.Println())
	}
	buff.WriteByte('\n')
	return buff.Bytes()
}

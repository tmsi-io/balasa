package sdp

import (
	"bytes"
	"fmt"
	"github.com/tmsi-io/balasa/rtp"
)

type MediaDescription struct {
	Type     string // "Audio" / "Video"
	Port     int
	Formats  string // "0"
	Protocol string // "TCP/RTP/AVP",
}

func (desc MediaDescription) Println() string {
	return fmt.Sprintf("%s %d %s %s", desc.Type, desc.Port, desc.Protocol, desc.Formats)
}

type Media struct {
	Title     string
	ConnInfo  []*SDPConnection
	Bandwidth map[string]int
	Desc      MediaDescription
	Encrypt   string
	Attrs     []string
}

func (media *Media) Init(mediaType string, mPayloadType, bandwidth int, trackID int, ClockRate int, fmtpCode int, sps string) {
	media.Bandwidth = make(map[string]int)
	media.Bandwidth["AS"] = bandwidth // 预留
	des := MediaDescription{}         //
	des.Protocol = ProtocolRTPAVP     // RTP/AVP
	des.Type = mediaType              //
	des.Port = 0
	media.AddAttribute(fmt.Sprintf("%s%d", TrackControl, trackID))
	if mediaType == MediaTypeVideoStr {
		if mPayloadType == rtp.PayloadH265 {
			media.AddAttribute(fmt.Sprintf("%s%d", AttrRTPMapH265, ClockRate*1000))
		} else if mPayloadType == rtp.PayloadH264 {
			if fmtpCode != 0 && sps != "" {
				media.AddAttribute(fmt.Sprintf("%s%d %s%s", AttrFmtp, fmtpCode, AttrSps, sps))
			}
			media.AddAttribute(fmt.Sprintf("%s%d", AttrRTPMapH264, ClockRate*1000))
		} else if mPayloadType == rtp.PayloadTS {
			media.AddAttribute(fmt.Sprintf("%s%d", AttrRTPMapTS, ClockRate*1000))
		} else {
			fmt.Println("No map Payload :", mPayloadType)
		}
	} else if mediaType == MediaTypeAudioStr {
		if mPayloadType == rtp.PayloadPCMA {
			media.AddAttribute(fmt.Sprintf("%s%d", AttrRTPMapPCMA, ClockRate*1000))
		} else if mPayloadType == rtp.PayloadG729 {
			media.AddAttribute(fmt.Sprintf("%s%d", AttrRTPMapG729, ClockRate*1000))
		} else if mPayloadType == rtp.PayloadG723 {
			media.AddAttribute(fmt.Sprintf("%s%d", AttrRTPMapG723, ClockRate*1000))
		} else {
			fmt.Println("No map Payload :", mPayloadType)
		}

	}
	media.AddAttribute(AttrRecvOnly)
	des.Formats = fmt.Sprintf("%d", mPayloadType)
	media.Desc = des
}

func (media *Media) AddAttribute(attr string) {
	media.Attrs = append(media.Attrs, attr)
}

func (media *Media) Println() *bytes.Buffer {
	buff := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buff, "m=%s\r\n", media.Desc.Println())
	for _, attr := range media.Attrs {
		_, _ = fmt.Fprintf(buff, "a=%s\r\n", attr)
	}
	//_, _ = fmt.Fprintf(buff, "k=%s", media.Encrypt)
	return buff
}

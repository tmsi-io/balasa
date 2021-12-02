package net

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tmsi-io/balasa/pes"
	"github.com/tmsi-io/balasa/rtp"
	"net"
	"sync"
	"time"
)

type TsFrameItem struct {
	Frame *bytes.Buffer
	Ts    uint32
}

type Uploader struct {
	BadCount         uint32
	SSRC             int
	sSSRC            string
	Channel          string
	ClientAddr       string
	Conn             *RichConn         //
	connRW           *bufio.ReadWriter //
	InBytes          int               // stats info
	Stoped           bool
	FrameTimeStamp   uint32
	LossPack         uint32
	DecodeFailed     uint32
	FirstFrame       bool
	JumpSeq          uint16
	VPayloadType     int
	APayloadType     int
	StartAt          time.Time
	StopDesc         string
	StopCode         int
	AudioTs          uint32
	BAudioFirst      bool
	AudioFirstTs     int64
	VideoTs          uint32
	BVideoFirst      bool
	VideoFirstTs     int64
	onceGet          sync.Once
	PsDecode         pes.DecPSPackage
	LastFrameInput   time.Time
	MaxFrameInterval int64
	pool             *Pool
}

func NewUploader(conn net.Conn) {
	networkBuffer := 1024 * 1024 * 1
	timeoutTCPConn := &RichConn{
		Conn:         conn,
		WriteTimeout: 4 * time.Second,
		ReadTimeout:  4 * time.Second,
	}
	var upload = Uploader{
		Conn: timeoutTCPConn,
		connRW: bufio.NewReadWriter(
			bufio.NewReaderSize(timeoutTCPConn, networkBuffer),
			bufio.NewWriterSize(timeoutTCPConn, networkBuffer)),
		SSRC:       0,
		FirstFrame: true,
		StartAt:    time.Now(),
		ClientAddr: conn.RemoteAddr().String(),
		PsDecode:   pes.DecPSPackage{},
		pool:       NewPool(128),
	}
	go upload.Start()
}

func (upload *Uploader) Start() {
	upload.BVideoFirst = true
	upload.BAudioFirst = true
	upload.DealReceivedGBStream()
	upload.Stop()
}

func (upload *Uploader) Stop() {
	defer func() {
		fmt.Println("Uploader Stop! ")
	}()
	if upload.Stoped {
		return
	} else {
		if len(upload.StopDesc) == 0 {
			if len(upload.Channel) > 0 {
				upload.StopDesc = fmt.Sprintf("%s-Normal", upload.Channel[0:5])
			} else {
				upload.StopDesc = fmt.Sprintf("%s-Normal", "Unknown")
			}
		}
		if upload.Conn != nil {
			_ = upload.connRW.Flush()
			_ = upload.Conn.Close()
			upload.Conn = nil
			upload.Stoped = true
		}
	}
}

// SetStreamInfo reset rate clock with receive rate
func (upload *Uploader) SetStreamInfo(VideoType uint32, AudioType uint32) {
	Now := time.Now()
	if upload.BVideoFirst {
		if VideoType == pes.PStreamTypeH264 {
			upload.VPayloadType = rtp.PayloadH264
		} else if VideoType == pes.PStreamTypeH265 {
			upload.VPayloadType = rtp.PayloadH265
		}
		upload.VideoFirstTs = Now.UnixNano() / int64(time.Millisecond)
		upload.BVideoFirst = false
	}
	v := Now.UnixNano()/int64(time.Millisecond) - upload.VideoFirstTs
	if v < 0 || v*rtp.DefaultVideoRate > 4294967295 {
		upload.VideoFirstTs = Now.UnixNano() / int64(time.Millisecond)
		v = 0
	}
	upload.VideoTs = uint32(v * rtp.DefaultVideoRate)
	if upload.BAudioFirst {
		if AudioType == pes.PStreamTypeG711 {
			upload.APayloadType = rtp.PayloadPCMA
		} else if AudioType == pes.PStreamTypeG722 {
			upload.APayloadType = rtp.PayloadG722
		}
		upload.AudioFirstTs = Now.UnixNano() / int64(time.Millisecond)
		upload.BAudioFirst = false
	}
	a := Now.UnixNano()/int64(time.Millisecond) - upload.AudioFirstTs
	if a < 0 || a*rtp.DefaultAudioRate > 4294967295 {
		upload.AudioFirstTs = Now.UnixNano() / int64(time.Millisecond)
		a = 0
	}
	upload.AudioTs = uint32(a * rtp.DefaultAudioRate)
}

package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tmsi-io/balasa/buff"
	"github.com/tmsi-io/balasa/rtp"
	"io"
	"strconv"
	"sync/atomic"
	"time"
)

/*
	TCP Header 0x24
*/

//1 2 3 4 5 6 7 8 1 2 3 4 5 6 7 8   16
//+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//|V=2|P|X|CC =4|M| PT =7       |

func (upload *Uploader) DealReceivedGBStream() {
	/*
		接收国标推流连接，拆分RTP包
	*/
	_buff := buff.BufferPool.Get().(*bytes.Buffer)
	defer func() {
		_buff.Reset()
		buff.BufferPool.Put(_buff)
	}()

	var count int
	for !upload.Stoped {
		upload.pool.ReInit()
		if upload.BadCount >= 20 { // 设备存在问题时，会有发空包的情况, 防止空跑
			break
		}
		rtpHeaderLen := upload.pool.Get(2)
		if _, err := io.ReadFull(upload.connRW, rtpHeaderLen); err == nil { // 读取RTP包头长度信息
			rtpLen := int(binary.BigEndian.Uint16(rtpHeaderLen))
			if rtpLen >= 65535 {
				atomic.AddUint32(&upload.BadCount, 1)
				continue
			}
			rtpBytes := upload.pool.Get(rtpLen)
			if _, err := io.ReadFull(upload.connRW, rtpBytes); err == nil {
				upload.InBytes += rtpLen + 2
				if ssrc, Timestamp, Seq, _, PRTPData, err := rtp.DecodeRTPInfoFromByte(rtpBytes); err == nil {
					if upload.SSRC == 0 {
						strSSRC := strconv.Itoa(ssrc)
						fmt.Println(strSSRC)
					}
					if upload.FirstFrame {
						upload.FrameTimeStamp = Timestamp // 设置初始值
						upload.JumpSeq = Seq
						upload.FirstFrame = false
					}
					if Seq-upload.JumpSeq < 0 {
						continue
					}
					if Timestamp == upload.FrameTimeStamp {
						count++
						if PRTPData != nil {
							_buff.Write(PRTPData)
						}
					} else {
						// count lost package
						loss := int(Seq-1) - int(upload.JumpSeq) - count
						if loss > 0 { // count lost
							atomic.AddUint32(&upload.LossPack, uint32(loss))
						}
						if _buff.Len() > 0 {
							tCost := time.Now().Sub(upload.LastFrameInput).Milliseconds()
							upload.LastFrameInput = time.Now()
							if tCost > upload.MaxFrameInterval {
								upload.MaxFrameInterval = tCost
							}
							upload.PushDataWithChannel(_buff, upload.FrameTimeStamp)
						}
						//重置缓冲区
						count = 0
						_buff.Reset()
						_buff.Write(PRTPData)
						upload.FrameTimeStamp = Timestamp
						upload.JumpSeq = Seq // set frame jump seq
					}
				} else {
					atomic.AddUint32(&upload.BadCount, 1)
				}
			} else {
				break
			}
		} else {
			break
		}
	}
}

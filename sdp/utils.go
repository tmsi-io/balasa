package sdp

import "github.com/tmsi-io/balasa/rtp"

func GetSdpInfo(PayloadVideo int, PayloadAudio int, VideoClock int, AudioClock int, VideoTrackID int, AudioTrackID int, StreamType int, FrameRate int, fmtpCode int, sps string) *SDPPayload {
	var sdp SDPPayload
	sdp.Version = SDPVersion
	sdp.SessionName = SessionNameStr
	sdp.Connection = DefaultConnection
	sdp.Medias = []*Media{}
	switch StreamType {
	case StreamType_ES:
		mediaVideo := Media{}
		mediaVideo.Init(MediaTypeVideoStr, PayloadVideo, 0, VideoTrackID, VideoClock, fmtpCode, sps)
		sdp.AddMedia(&mediaVideo)
		if PayloadAudio != 0 && EnableAudio {
			mediaAudio := Media{}
			mediaAudio.Init(MediaTypeAudioStr, PayloadAudio, 0, AudioTrackID, AudioClock, fmtpCode, sps)
			sdp.AddMedia(&mediaAudio)
		}
	case StreamType_TS:
		mediaVideo := Media{}
		mediaVideo.Init(MediaTypeVideoStr, rtp.PayloadTS, 0, VideoTrackID, VideoClock, fmtpCode, sps)
		sdp.AddMedia(&mediaVideo)
	default:
	}
	sdp.Origin = new(SDPOrigin)
	sdp.Origin.Init()
	return &sdp
}

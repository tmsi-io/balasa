package sdp

const SDPVersion = "0"
const MediaTypeAudioStr = "audio"
const MediaTypeVideoStr = "video"
const AttrRTPMapH264 = "rtpmap:96 H264/" //
const AttrRTPMapTS = "rtpmap:33 MP2T/"   //
const AttrRTPMapH265 = "rtpmap:98 H265/"
const AttrRTPMapPCMA = "rtpmap:8 PCMA/"
const AttrRTPMapG729 = "rtpmap:18 G729/"
const AttrRTPMapG723 = "rtpmap:4 G723/"
const AttrRecvOnly = "recvonly"
const AttrFrameRate = "framerate:"
const AttrFmtp = "fmtp:"
const AttrSps = "sprop-parameter-sets="
const SessionNameStr = "Sctel MoJing MediaStream Server"
const DefaultConnection = "IN IP4 127.0.0.1"
const ProtocolRTPAVP = "RTP/AVP"
const TrackControl = "control:trackID="

const (
	StreamType_ES = 0
	StreamType_TS = 1
)

var EnableAudio bool

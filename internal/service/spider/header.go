package spider

import "whisper/pkg/http"

var (
	Cookie []http.Header = []http.Header{
		http.Header{
			Key:   "Cookie",
			Value: "buvid3=24B36EDD-1DA7-9AF2-8B64-F75104E0C30C46082infoc; b_nut=1689649646; i-wanna-go-back=-1; _uuid=5C9D10A6D-86C6-D310B-CC39-127E10BEA22F1044547infoc; FEED_LIVE_VERSION=V8; DedeUserID=32479325; DedeUserID__ckMd5=812a87eb40af15a8; CURRENT_FNVAL=4048; rpdid=0zbfVCLOOB|C5F9dZ7x|3vU|3w1QlBbS; buvid_fp_plain=undefined; hit-new-style-dyn=1; hit-dyn-v2=1; b_ut=5; LIVE_BUVID=AUTO3316914987245464; buvid4=705F6682-26A3-0364-1BFE-F533506DCF4447711-023071811-f6C3yQ8RJpZ8qk%2FMIitUYawpVKA3GyrjOtP51e3DbcblQvUx83H%2F1A%3D%3D; iflogin_when_web_push=1; CURRENT_QUALITY=80; enable_web_push=ENABLE; header_theme_version=CLOSE; home_feed_column=5; SESSDATA=ccf90982%2C1717206468%2Cadc9a%2Ac1CjAizeXm0Onk_fRxjs5O6_pzOtTLTkqKZJiv3sh5OTpHvwlwk46Uq-9mFIC0L6z55P4SVlowSWZQSTU0aFhra2xLWjJzclR3b1Q1MkY4RFYtd1I1azJVYldIdG9QRHVVUWtSdloxZmVrQUdZZXJLTmRzRFo1UjI0RGdhTVY5RF9Nb2FJdS02Qml3IIEC; bili_jct=515eb286de67ea315b5da661ce3759b7; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE5MTM2ODAsImlhdCI6MTcwMTY1NDQyMCwicGx0IjotMX0.0aFnJAtGOCjeYCAdTbkjOSmiVK9qjbVFNRDGu11MvDQ; bili_ticket_expires=1701913620; sid=6f3kkh9l; fingerprint=a68faa493e76410fc6d949665a3ead50; buvid_fp=f7e5828d6e30a99ab5793a2e7eeab0a9; browser_resolution=1671-734; bp_video_offset_32479325=871270511565340689; b_lsid=FCBDF5EC_18C3570C644; PVID=1",
		},
	}

	Referer []http.Header = []http.Header{
		http.Header{
			Key:   "Referer",
			Value: "https://space.bilibili.com",
		},
	}

	UA []http.Header = []http.Header{
		http.Header{
			Key:   "user-agent",
			Value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		},
	}
)

func CommonHeaders() []http.Header {
	var h []http.Header
	h = append(h, Referer...)
	h = append(h, UA...)
	return h
}

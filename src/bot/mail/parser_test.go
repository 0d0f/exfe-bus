package mail

import (
	"bytes"
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"model"
	"net/mail"
	"testing"
)

var config model.Config

func init() {
	config.Email.Prefix = "x"
	config.Email.Domain = "0d0f.com"
}

func TestParsePlain(t *testing.T) {
	type Test struct {
		input  string
		output string
	}
	var tests = []Test{
		{"测试普通邮件\n\n在 2013-03-19 13:10:34，\"Googol Lee\" <googollee@gmail.com> 写道：\n\n测试普通邮件回复", "测试普通邮件"},
	}
	for i, test := range tests {
		output := parsePlain(test.input)
		assert.Equal(t, output, test.output, fmt.Sprintf("test %d", i))
	}
}

func TestParseContentType(t *testing.T) {
	type Test struct {
		contentType string
		mime        string
		pairs       string
	}
	var tests = []Test{
		{"multipart/alternative; boundary=\"----=_Part_134558_1130028024.1363669790299", "multipart/alternative", "map[boundary:----=_Part_134558_1130028024.1363669790299]"},
	}
	for i, test := range tests {
		mime, pairs := parseContentType(test.contentType)
		assert.Equal(t, mime, test.mime, fmt.Sprintf("test %d", i))
		assert.Equal(t, fmt.Sprintf("%v", pairs), test.pairs, fmt.Sprintf("test %d", i))
	}
}

func TestNormalAttachment(t *testing.T) {
	var str = `Received: from googollee$163.com ( [116.237.198.97] ) by
 ajax-webmail-wmsvr64 (Coremail) ; Tue, 19 Mar 2013 16:05:42 +0800 (CST)
X-Originating-IP: [116.237.198.97]
Date: Tue, 19 Mar 2013 16:05:42 +0800 (CST)
From: googollee  <googollee@163.com>
To: "Googol Lee" <googollee@gmail.com>
Cc: =?GBK?Q?=5BDEV=5D_EXFE_=A1=A4X=A1=A4?= <x@0d0f.com>
Subject: =?utf-8?Q?=E9=82=AE=E4=BB=B6=E5=88=9B=E5=BB=BAX=E6=B5=8B?=
 =?utf-8?Q?=E8=AF=95_take2?=
X-Priority: 3
X-Mailer: Coremail Webmail Server Version SP_ntes V3.5 build
 20130201(21528.5249.5248) Copyright (c) 2002-2013 www.mailtech.cn 163com
In-Reply-To: <CAOf82vOyhNiqKTqw79nEQePxNP2r+RgZj3UujbBBaEsQjAGrKg@mail.gmail.com>
References: <CAOf82vPe2Nt1hs6VVeDcczxvrqvvzCYdzUjv3vCWT1XK1MkFyQ@mail.gmail.com>
 <365e2f32.8ad2.13d810a1b28.Coremail.googollee@163.com>
 <CAOf82vOyhNiqKTqw79nEQePxNP2r+RgZj3UujbBBaEsQjAGrKg@mail.gmail.com>
X-CM-CTRLDATA: WWdja2Zvb3Rlcl9odG09MTI3ODo4MQ==
Content-Type: multipart/mixed; 
	boundary="----=_Part_207651_1797332469.1363680342939"
MIME-Version: 1.0
Message-ID: <231eec8.da02.13d81aeb39b.Coremail.googollee@163.com>

------=_Part_207651_1797332469.1363680342939
Content-Type: multipart/alternative; 
	boundary="----=_Part_207653_1138501142.1363680342939"

------=_Part_207653_1138501142.1363680342939
Content-Type: text/plain; charset=GBK
Content-Transfer-Encoding: base64

CtTZtM64/NDCtdi14woKCgoKCtTaIDIwMTMtMDMtMTkgMTM6MTA6MzSjrCJHb29nb2wgTGVlIiA8
Z29vZ29sbGVlQGdtYWlsLmNvbT4g0LS1wKO6Cgqy4srUxtXNqNPKvP672Li0CgoKCjIwMTMvMy8x
OSBnb29nb2xsZWUgPGdvb2dvbGxlZUAxNjMuY29tPgqy4srUxtXNqNPKvP672Li0CgrU2iAyMDEz
LTAzLTE5IDEzOjA0OjQzo6wiR29vZ29sIExlZSIgPGdvb2dvbGxlZUBnbWFpbC5jb20+INC0tcCj
ugoKCrS0vaiwoaOho6GjoQoKCgotLQrQwrXEwO3C27TTydnK/cjLtcTW99XFtb3Su82zzOzPwqOs
sqKyu8rH0vLOqtXiuPbA7cLby7W3/sHLsfDIy8XXxvq+ybnbteOjrLb4ysfS8s6q0ru0+sjLtcTK
xciloaMKCgoKCgoKCgotLQrQwrXEwO3C27TTydnK/cjLtcTW99XFtb3Su82zzOzPwqOssqKyu8rH
0vLOqtXiuPbA7cLby7W3/sHLsfDIy8XXxvq+ybnbteOjrLb4ysfS8s6q0ru0+sjLtcTKxciloaM=

------=_Part_207653_1138501142.1363680342939
Content-Type: text/html; charset=GBK
Content-Transfer-Encoding: base64

PGRpdiBzdHlsZT0ibGluZS1oZWlnaHQ6MS43O2NvbG9yOiMwMDAwMDA7Zm9udC1zaXplOjE0cHg7
Zm9udC1mYW1pbHk6YXJpYWwiPjxicj7U2bTOuPzQwrXYteM8YnI+PGJyPjxicj48YnI+PGRpdj48
L2Rpdj48ZGl2IGlkPSJkaXZOZXRlYXNlTWFpbENhcmQiPjwvZGl2Pjxicj7U2iAyMDEzLTAzLTE5
IDEzOjEwOjM0o6wiR29vZ29sJm5ic3A7TGVlIiZuYnNwOyZsdDtnb29nb2xsZWVAZ21haWwuY29t
Jmd0OyDQtLXAo7o8YnI+IDxibG9ja3F1b3RlIGlkPSJpc1JlcGx5Q29udGVudCIgc3R5bGU9IlBB
RERJTkctTEVGVDogMWV4OyBNQVJHSU46IDBweCAwcHggMHB4IDAuOGV4OyBCT1JERVItTEVGVDog
I2NjYyAxcHggc29saWQiPjxkaXYgZGlyPSJsdHIiPrLiytTG1c2o08q8/rvYuLQ8L2Rpdj48ZGl2
IGNsYXNzPSJnbWFpbF9leHRyYSI+PGJyPjxicj48ZGl2IGNsYXNzPSJnbWFpbF9xdW90ZSI+MjAx
My8zLzE5IGdvb2dvbGxlZSA8c3BhbiBkaXI9Imx0ciI+Jmx0OzxhIGhyZWY9Im1haWx0bzpnb29n
b2xsZWVAMTYzLmNvbSIgdGFyZ2V0PSJfYmxhbmsiPmdvb2dvbGxlZUAxNjMuY29tPC9hPiZndDs8
L3NwYW4+PGJyPjxibG9ja3F1b3RlIGNsYXNzPSJnbWFpbF9xdW90ZSIgc3R5bGU9Im1hcmdpbjow
IDAgMCAuOGV4O2JvcmRlci1sZWZ0OjFweCAjY2NjIHNvbGlkO3BhZGRpbmctbGVmdDoxZXgiPgoK
suLK1MbVzajTyrz+u9i4tDxicj48YnI+1NogMjAxMy0wMy0xOSAxMzowNDo0M6OsIkdvb2dvbCZu
YnNwO0xlZSImbmJzcDsmbHQ7PGEgaHJlZj0ibWFpbHRvOmdvb2dvbGxlZUBnbWFpbC5jb20iIHRh
cmdldD0iX2JsYW5rIj5nb29nb2xsZWVAZ21haWwuY29tPC9hPiZndDsg0LS1wKO6PGRpdiBjbGFz
cz0iSE9FblpiIj48ZGl2IGNsYXNzPSJoNSI+PGJyPiA8YmxvY2txdW90ZSBzdHlsZT0iUEFERElO
Ry1MRUZUOjFleDtNQVJHSU46MHB4IDBweCAwcHggMC44ZXg7Qk9SREVSLUxFRlQ6I2NjYyAxcHgg
c29saWQiPgoKPGRpdiBkaXI9Imx0ciI+tLS9qLCho6GjoaOhPGJyIGNsZWFyPSJhbGwiPjxkaXY+
PGJyPjwvZGl2Pi0tIDxicj7QwrXEwO3C27TTydnK/cjLtcTW99XFtb3Su82zzOzPwqOssqKyu8rH
0vLOqtXiuPbA7cLby7W3/sHLsfDIy8XXxvq+ybnbteOjrLb4ysfS8s6q0ru0+sjLtcTKxciloaMK
PC9kaXY+CjwvYmxvY2txdW90ZT48YnI+PGJyPjxzcGFuIHRpdGxlPSJuZXRlYXNlZm9vdGVyIj48
c3Bhbj48L3NwYW4+PC9zcGFuPjwvZGl2PjwvZGl2PjwvYmxvY2txdW90ZT48L2Rpdj48YnI+PGJy
IGNsZWFyPSJhbGwiPjxkaXY+PGJyPjwvZGl2Pi0tIDxicj7QwrXEwO3C27TTydnK/cjLtcTW99XF
tb3Su82zzOzPwqOssqKyu8rH0vLOqtXiuPbA7cLby7W3/sHLsfDIy8XXxvq+ybnbteOjrLb4ysfS
8s6q0ru0+sjLtcTKxciloaMKPC9kaXY+CjwvYmxvY2txdW90ZT48L2Rpdj48YnI+PGJyPjxzcGFu
IHRpdGxlPSJuZXRlYXNlZm9vdGVyIj48c3BhbiBpZD0ibmV0ZWFzZV9tYWlsX2Zvb3RlciI+PC9z
cGFuPjwvc3Bhbj4=
------=_Part_207653_1138501142.1363680342939--

------=_Part_207651_1797332469.1363680342939
Content-Type: application/octet-stream; name="ics.ics"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="ics.ics"

QkVHSU46VkNBTEVOREFSDQpNRVRIT0Q6UkVRVUVTVA0KVkVSU0lPTjoyLjANClBST0RJRDotLy9N
aWNyb3NvZnQgQ29ycG9yYXRpb24vL1dpbmRvd3MgTGl2ZSBDYWxlbmRhci8vRU4NCkJFR0lOOlZU
SU1FWk9ORQ0KVFpJRDpDaGluYSBTdGFuZGFyZCBUaW1lDQpCRUdJTjpTVEFOREFSRA0KRFRTVEFS
VDoyMDA4MDEwMVQwMDAwMDANClRaT0ZGU0VUVE86KzA4MDANClRaT0ZGU0VURlJPTTorMDgwMA0K
RU5EOlNUQU5EQVJEDQpFTkQ6VlRJTUVaT05FDQpCRUdJTjpWRVZFTlQNClVJRDphN2JiMWQyZi0z
ZWYxLTQyNTEtOWJjMi1kMGMyNGRlZTI0ZDUNCkRUU1RBTVA6MjAxMzAzMTlUMDczODE5Wg0KQ0xB
U1M6UFVCTElDDQpYLU1JQ1JPU09GVC1DRE8tQlVTWVNUQVRVUzpCVVNZDQpUUkFOU1A6T1BBUVVF
DQpTRVFVRU5DRTowDQpEVFNUQVJUO1RaSUQ9Q2hpbmEgU3RhbmRhcmQgVGltZToyMDEzMDMyMFQw
OTAwMDANCkRURU5EO1RaSUQ9Q2hpbmEgU3RhbmRhcmQgVGltZToyMDEzMDMyMFQxMDAwMDANClNV
TU1BUlk65p2l6Ieqb3V0bG9vayBjYWxlbmRhcg0KTE9DQVRJT0465Zyw54K5MjMzDQpQUklPUklU
WTowDQpBVFRFTkRFRTtDVVRZUEU9SU5ESVZJRFVBTDtST0xFPVJFUS1QQVJUSUNJUEFOVDtQQVJU
U1RBVD1ORUVEUy1BQ1RJT047UlNWUD0NCiBUUlVFOk1BSUxUTzp4QDBkMGYuY29tDQpBVFRFTkRF
RTtDVVRZUEU9SU5ESVZJRFVBTDtST0xFPVJFUS1QQVJUSUNJUEFOVDtQQVJUU1RBVD1ORUVEUy1B
Q1RJT047UlNWUD0NCiBUUlVFOk1BSUxUTzpnb29nb2xsZWVAMTYzLmNvbQ0KT1JHQU5JWkVSO0NO
PUxlZSBHb29nb2wgTGVlOk1BSUxUTzpnb29nb2xsZWVAaG90bWFpbC5jb20NCkJFR0lOOlZBTEFS
TQ0KQUNUSU9OOkRJU1BMQVkNClRSSUdHRVI6LVBUMTVNDQpFTkQ6VkFMQVJNDQpCRUdJTjpWQUxB
Uk0NCkFDVElPTjpESVNQTEFZDQpUUklHR0VSOi1QVDE1TQ0KRU5EOlZBTEFSTQ0KRU5EOlZFVkVO
VA0KRU5EOlZDQUxFTkRBUg0K
------=_Part_207651_1797332469.1363680342939--`

	buf := bytes.NewBufferString(str)
	msg, err := mail.ReadMessage(buf)
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewParser(msg, &config)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, parser.messageID, "231eec8.da02.13d81aeb39b.Coremail.googollee@163.com")
	assert.Equal(t, parser.content, "再次更新地点")
	assert.Equal(t, parser.subject, "邮件创建X测试 take2")
	assert.Equal(t, parser.event.ID, "a7bb1d2f-3ef1-4251-9bc2-d0c24dee24d5")
	assert.Equal(t, parser.HasICS(), true)
	assert.Equal(t, parser.GetCross().Place.Title, "地点233")
}

func Test163(t *testing.T) {
	var str = `Received: from googollee$163.com ( [116.237.198.97] ) by
 ajax-webmail-wmsvr64 (Coremail) ; Tue, 19 Mar 2013 13:09:50 +0800 (CST)
X-Originating-IP: [116.237.198.97]
Date: Tue, 19 Mar 2013 13:09:50 +0800 (CST)
From: googollee  <googollee@163.com>
To: x+c236@0d0f.com
Subject: =?GBK?B?UmU6suLK1NPKvP60tL2oNQ==?=
X-Priority: 3
X-Mailer: Coremail Webmail Server Version SP_ntes V3.5 build
 20130201(21528.5249.5248) Copyright (c) 2002-2013 www.mailtech.cn 163com
In-Reply-To: <5147f215.87f1420a.5619.ffff95d4@mx.google.com>
References: <x+c236@exfe.com>
 <5147f215.87f1420a.5619.ffff95d4@mx.google.com>
X-CM-CTRLDATA: yAgEyWZvb3Rlcl9odG09ODA1Nzo4MQ==
Content-Type: multipart/alternative; 
	boundary="----=_Part_134558_1130028024.1363669790299"
MIME-Version: 1.0
Message-ID: <688b5ae2.8bc9.13d810dae5b.Coremail.googollee@163.com>

------=_Part_134558_1130028024.1363669790299
Content-Type: text/plain; charset=GBK
Content-Transfer-Encoding: base64

suLK1GV4ZmXTyrz+u9i4tAoKCgoKCgrU2iAyMDEzLTAzLTE5IDEzOjA1OjI1o6wiW0RFVl1FWEZF
IKGkWKGkIiA8eEAwZDBmLmNvbT4g0LS1wKO6Cgp8ILLiytTTyrz+tLS9qDUgfAp8CnwgSW52aXRh
dGlvbiBmcm9tIEdvb2dvbCBMZWUuIHwKfAp8CgpUaW1lCgpUbyBiZSBkZWNpZGVkCgp8CgpQbGFj
ZQoKVG8gYmUgZGVjaWRlZAoKfAp8CnwgSSdtIGluQ2hlY2sgaXQgb3V0Li4uIHwKfCC0tL2osKGj
oaOho6EgfAp8IHwKfCB8CnwKfCB8IEdvb2dvbCBMZWVob3N0IHwKfCB8IGdvb2dvbGxlZSB8CnwK
fAp8IFJlcGx5IHRoaXMgZW1haWwgZGlyZWN0bHkgYXMgY29udmVyc2F0aW9uLCBvciBUcnkgRVhG
RSBhcHAuClRoaXMgoaRYoaQgaW52aXRhdGlvbiBpcyBzZW50IGJ5IEdvb2dvbCBMZWUgZnJvbSBF
WEZFLiB8
------=_Part_134558_1130028024.1363669790299
Content-Type: text/html; charset=GBK
Content-Transfer-Encoding: base64

PGRpdiBzdHlsZT0ibGluZS1oZWlnaHQ6MS43O2NvbG9yOiMwMDAwMDA7Zm9udC1zaXplOjE0cHg7
Zm9udC1mYW1pbHk6YXJpYWwiPrLiytRleGZl08q8/rvYuLQ8YnI+PGJyPjxicj48YnI+PGJyPjxk
aXY+PC9kaXY+PGRpdiBpZD0iZGl2TmV0ZWFzZU1haWxDYXJkIj48L2Rpdj48YnI+1NogMjAxMy0w
My0xOSAxMzowNToyNaOsIltERVZdRVhGRSZuYnNwO6GkWKGkIiZuYnNwOyZsdDt4QDBkMGYuY29t
Jmd0OyDQtLXAo7o8YnI+IDxibG9ja3F1b3RlIGlkPSJpc1JlcGx5Q29udGVudCIgc3R5bGU9IlBB
RERJTkctTEVGVDogMWV4OyBNQVJHSU46IDBweCAwcHggMHB4IDAuOGV4OyBCT1JERVItTEVGVDog
I2NjYyAxcHggc29saWQiPgogICAgCiAgICAgICAgCiAgICAgICAgCiAgICAgICAgPHN0eWxlPgog
ICAgICAgICAgICAuZXhmZV9tYWlsX2xhYmVsIHsKICAgICAgICAgICAgICAgIGJhY2tncm91bmQt
Y29sb3I6ICNENUU4RjI7CiAgICAgICAgICAgICAgICBjb2xvcjogIzNhNmVhNTsKICAgICAgICAg
ICAgICAgIGZvbnQtc2l6ZTogMTFweDsKICAgICAgICAgICAgICAgIHBhZGRpbmc6IDAgMnB4IDAg
MnB4OwogICAgICAgICAgICB9CiAgICAgICAgICAgIC5leGZlX21haWxfbWF0ZXMgewogICAgICAg
ICAgICAgICAgY29sb3I6ICMzYTZlYTU7CiAgICAgICAgICAgICAgICBmb250LXNpemU6IDEycHg7
CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0eSB7CiAgICAgICAg
ICAgICAgICBmb250LXN0eWxlOiBpdGFsaWM7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4
ZmVfbWFpbF9pZGVudGl0eV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjMTkxOTE5Owog
ICAgICAgICAgICB9CiAgICAgICAgPC9zdHlsZT4KICAgIAogICAgCiAgICAgICAgPHRhYmxlIHdp
ZHRoPSI2NDAiIGJvcmRlcj0iMCIgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHls
ZT0iZm9udC1mYW1pbHk6IEhlbHZldGljYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDog
MTlweDsgY29sb3I6ICMxOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHBhZGRpbmc6IDMwcHgg
NDBweCAzMHB4IDQwcHg7IGJhY2tncm91bmQtY29sb3I6ICNmYmZiZmI7IG1pbi1oZWlnaHQ6IDU2
MnB4OyI+CiAgICAgICAgICAgIDx0Ym9keT48dHI+CiAgICAgICAgICAgICAgICA8dGQgY29sc3Bh
bj0iMyIgdmFsaWduPSJ0b3AiIHN0eWxlPSJmb250LXNpemU6IDMycHg7IGxpbmUtaGVpZ2h0OiAz
OHB4OyBwYWRkaW5nLWJvdHRvbTogMThweDsiPgogICAgICAgICAgICAgICAgICAgIDxhIGhyZWY9
Imh0dHA6Ly8wZDBmLmNvbS8jIXRva2VuPTkzZTA2ZjkzNDUyY2FlMzA5NmM3YzMzODg3YjI1MGRl
IiBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsgZm9udC13ZWln
aHQ6IDMwMDsiPgogICAgICAgICAgICAgICAgICAgICAgICCy4srU08q8/rS0vag1CiAgICAgICAg
ICAgICAgICAgICAgPC9hPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4K
ICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIzNDAiIHN0eWxlPSJ2
ZXJ0aWNhbC1hbGlnbjogYmFzZWxpbmU7IGZvbnQtd2VpZ2h0OiAzMDA7Ij4KICAgICAgICAgICAg
ICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAi
PgogICAgICAgICAgICAgICAgICAgICAgICA8dGJvZHk+PHRyPgogICAgICAgICAgICAgICAgICAg
ICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiBzdHlsZT0icGFkZGluZy1ib3R0b206IDIwcHg7IGZv
bnQtc2l6ZTogMjBweDsgdmVydGljYWwtYWxpZ246IGJhc2VsaW5lOyI+CiAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgSW52aXRh
dGlvbiBmcm9tIDxzcGFuIHN0eWxlPSJmb250LXdlaWdodDo4MDA7Ij5Hb29nb2wgTGVlPC9zcGFu
Pi4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAg
ICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAg
ICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIg
Y2VsbHNwYWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGJv
ZHk+PHRyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGln
bj0idG9wIiB3aWR0aD0iMTYwIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICA8YSBocmVmPSJodHRwOi8vMGQwZi5jb20vIyF0b2tlbj05M2UwNmY5MzQ1MmNhZTMw
OTZjN2MzMzg4N2IyNTBkZSIgc3R5bGU9InRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAkKICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxpbmUtaGVp
Z2h0OiAyNnB4OyBtYXJnaW46IDA7IGNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBUaW1lCiAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9ImNvbG9yOiAjMTkxOTE5OyBtYXJn
aW46IDA7Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgIFRvIGJlIGRlY2lkZWQKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L2E+CiAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiBzdHlsZT0icGFkZGlu
Zy1sZWZ0OiAxMHB4OyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgPGEgaHJlZj0iaHR0cDovLzBkMGYuY29tLyMhdG9rZW49OTNlMDZmOTM0NTJjYWUzMDk2Yzdj
MzM4ODdiMjUwZGUiIHN0eWxlPSJ0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij4KICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDIwcHg7IGxp
bmUtaGVpZ2h0OiAyNnB4OyBtYXJnaW46IDA7IGNvbG9yOiAjMzMzMzMzOyI+CiAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICBQbGFjZQogICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJjb2xvcjogIzE5MTkx
OTsgbWFyZ2luOiAwOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICBUbyBiZSBkZWNpZGVkCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgIDwvcD4gCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
IDwvYT4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgPC90Ym9keT48L3RhYmxlPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg
PC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAg
ICAgCiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAg
ICAgIDx0ZCB2YWxpZ249InRvcCIgc3R5bGU9InBhZGRpbmctdG9wOiAzMHB4OyBwYWRkaW5nLWJv
dHRvbTogMzBweDsiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxhIHN0eWxlPSJm
bG9hdDogbGVmdDsgZGlzcGxheTogYmxvY2s7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsgYm9yZGVy
OiAxcHggc29saWQgI2JlYmViZTsgYmFja2dyb3VuZC1jb2xvcjogIzNBNkVBNTsgY29sb3I6ICNG
RkZGRkY7IHBhZGRpbmc6IDVweCAzMHB4IDVweCAzMHB4OyBtYXJnaW4tbGVmdDogMjVweDsiIGFs
dD0iQWNjZXB0IiBocmVmPSJodHRwOi8vMGQwZi5jb20vP3Rva2VuPTkzZTA2ZjkzNDUyY2FlMzA5
NmM3YzMzODg3YjI1MGRlJmFtcDtyc3ZwPWFjY2VwdCI+SSdtIGluPC9hPgogICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgIDxhIHN0eWxlPSJmbG9hdDogbGVmdDsgZGlzcGxheTogYmxvY2s7
IHRleHQtZGVjb3JhdGlvbjogbm9uZTsgYm9yZGVyOiAxcHggc29saWQgI2JlYmViZTsgYmFja2dy
b3VuZC1jb2xvcjogI0U2RTZFNjsgY29sb3I6ICMxOTE5MTk7IHBhZGRpbmc6IDVweCAyNXB4IDVw
eCAyNXB4OyBtYXJnaW4tbGVmdDogMTVweDsiIGFsdD0iQ2hlY2sgaXQgb3V0IiBocmVmPSJodHRw
Oi8vMGQwZi5jb20vIyF0b2tlbj05M2UwNmY5MzQ1MmNhZTMwOTZjN2MzMzg4N2IyNTBkZSI+Q2hl
Y2sgaXQgb3V0Li4uPC9hPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAg
ICAgICAgICAgICAgICAgICAgPC90cj4KICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAg
ICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgIDx0ZCB2YWxp
Z249InRvcCI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgtLS9qLCho6GjoaOhCiAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAg
ICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAg
IDwvdGJvZHk+PC90YWJsZT4KICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICA8
dGQgd2lkdGg9IjMwIj48L3RkPgogICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIj4KICAg
ICAgICAgICAgICAgICAgICA8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNw
YWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICA8dGJvZHk+PHRyPgogICAgICAgICAg
ICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGQgaGVpZ2h0
PSI2OCIgdmFsaWduPSJ0b3AiIGFsaWduPSJyaWdodCI+CiAgICAgICAgICAgICAgICAgICAgICAg
ICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAg
ICAgICAgIDwvdHI+CiAgICAgICAgICAgICAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAg
ICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8dGFibGUg
Ym9yZGVyPSIwIiBzdHlsZT0iY29sb3I6ICMzMzMzMzM7IiBjZWxscGFkZGluZz0iMCIgY2VsbHNw
YWNpbmc9IjAiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRib2R5Pjx0cj4KICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgIDx0ZCB3aWR0aD0iMjUiIGhlaWdodD0iMjUiIGFsaWduPSJs
ZWZ0IiB2YWxpZ249InRvcCI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgPGltZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHRpdGxlPSJHb29nb2wgTGVlIiBhbHQ9
Ikdvb2dvbCBMZWUiIHNyYz0iaHR0cDovL3d3dy5ncmF2YXRhci5jb20vYXZhdGFyLzE1YjdmYzFi
MTAxZWUyODliODE2Nzg4MTI3ODFhZWE2Ij4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8
dGQ+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxzcGFuPkdvb2dv
bCBMZWU8L3NwYW4+IDxzcGFuIGNsYXNzPSJleGZlX21haWxfbGFiZWwiPmhvc3Q8L3NwYW4+CiAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRkIHdpZHRoPSIyNSIgaGVpZ2h0PSIy
NSIgYWxpZ249ImxlZnQiIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCIgdGl0bGU9Imdvb2dv
bGxlZSIgYWx0PSJnb29nb2xsZWUiIHNyYz0iaHR0cDovL2FwaS4wZDBmLmNvbS92Mi9hdmF0YXIv
ZGVmYXVsdD9uYW1lPWdvb2dvbGxlZSI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPHRk
PgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8c3Bhbj5nb29nb2xs
ZWU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgog
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3Ri
b2R5PjwvdGFibGU+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAg
ICAgICAgICAgICAgICA8L3RyPgogICAgICAgICAgICAgICAgICAgIDwvdGJvZHk+PC90YWJsZT4K
ICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj4K
ICAgICAgICAgICAgICAgIDx0ZCBjb2xzcGFuPSIzIiBzdHlsZT0iZm9udC1zaXplOiAxMXB4OyBs
aW5lLWhlaWdodDogMTVweDsgY29sb3I6ICM3RjdGN0Y7IHBhZGRpbmctdG9wOiA0MHB4OyI+CiAg
ICAgICAgICAgICAgICAgICAgUmVwbHkgdGhpcyBlbWFpbCBkaXJlY3RseSBhcyBjb252ZXJzYXRp
b24sIG9yIFRyeSA8YSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9u
ZTsiIGhyZWY9Imh0dHA6Ly9pdHVuZXMuYXBwbGUuY29tL3VhL2FwcC9leGZlL2lkNTE0MDI2NjA0
Ij5FWEZFPC9hPiBhcHAuCiAgICAgICAgICAgICAgICAgICAgPGJyPgogICAgICAgICAgICAgICAg
ICAgIFRoaXMgPGEgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7
IiBocmVmPSJodHRwOi8vMGQwZi5jb20vIyF0b2tlbj05M2UwNmY5MzQ1MmNhZTMwOTZjN2MzMzg4
N2IyNTBkZSI+oaRYoaQ8L2E+IGludml0YXRpb24gaXMgc2VudCBieSA8c3BhbiBjbGFzcz0iZXhm
ZV9tYWlsX2lkZW50aXR5X25hbWUiPkdvb2dvbCBMZWU8L3NwYW4+IGZyb20gPGEgc3R5bGU9ImNv
bG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vMGQwZi5j
b20iPkVYRkU8L2E+LgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAg
ICAgICA8L3Rib2R5PjwvdGFibGU+CiAgICAKCjwvYmxvY2txdW90ZT48L2Rpdj48YnI+PGJyPjxz
cGFuIHRpdGxlPSJuZXRlYXNlZm9vdGVyIj48c3BhbiBpZD0ibmV0ZWFzZV9tYWlsX2Zvb3RlciI+
PC9zcGFuPjwvc3Bhbj4=
------=_Part_134558_1130028024.1363669790299--`

	buf := bytes.NewBufferString(str)
	msg, err := mail.ReadMessage(buf)
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewParser(msg, &config)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, parser.messageID, "688b5ae2.8bc9.13d810dae5b.Coremail.googollee@163.com")
	assert.Equal(t, parser.content, "测试exfe邮件回复")
}

func TestInit(t *testing.T) {
	var str = `Delivered-To: panda@0d0f.com
Received: by 10.70.0.234 with SMTP id 10csp193484pdh;
        Sun, 17 Mar 2013 23:15:37 -0700 (PDT)
X-Received: by 10.152.113.34 with SMTP id iv2mr13078933lab.20.1363587335984;
        Sun, 17 Mar 2013 23:15:35 -0700 (PDT)
Return-Path: <googollee@hotmail.com>
Received: from bay0-omc2-s22.bay0.hotmail.com (bay0-omc2-s22.bay0.hotmail.com. [65.54.190.97])
        by mx.google.com with ESMTP id ia4si6157124lab.180.2013.03.17.23.15.34;
        Sun, 17 Mar 2013 23:15:35 -0700 (PDT)
Received-SPF: pass (google.com: domain of googollee@hotmail.com designates 65.54.190.97 as permitted sender) client-ip=65.54.190.97;
Authentication-Results: mx.google.com;
       spf=pass (google.com: domain of googollee@hotmail.com designates 65.54.190.97 as permitted sender) smtp.mail=googollee@hotmail.com
Received: from BAY152-DS11 ([65.54.190.124]) by bay0-omc2-s22.bay0.hotmail.com with Microsoft SMTPSVC(6.0.3790.4675);
	 Sun, 17 Mar 2013 23:15:34 -0700
X-EIP: [cqX6eeEhbwCCOhFZWRs2ZHcUKc1BKGRH]
X-Originating-Email: [googollee@hotmail.com]
Message-ID: <BAY152-ds11D2686FF50F86678FFD1AA0E80@phx.gbl>
Return-Path: googollee@hotmail.com
From: "Lee Googol Lee" <googollee@hotmail.com>
To: googollee@gmail.com;panda@0d0f.com
Date: Sun, 17 Mar 2013 23:15:33 -0700
Subject: =?utf-8?B?5rS75Yqo?=
Content-Type: multipart/alternative;
	boundary="_32936406-5bb3-4fc0-b17a-11cc15875002_"
MIME-Version: 1.0
X-OriginalArrivalTime: 18 Mar 2013 06:15:34.0264 (UTC) FILETIME=[0008CF80:01CE23A0]

--_32936406-5bb3-4fc0-b17a-11cc15875002_
Content-Type: text/plain; charset="utf-8"
Content-Transfer-Encoding: base64


--_32936406-5bb3-4fc0-b17a-11cc15875002_
Content-Type: text/html; charset="utf-8"
Content-Transfer-Encoding: base64

DQoNCjxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+DQogICAgDQogICAgDQogICAgLk1haW5Db250YWlu
ZXINCiAgICB7DQogICAgICAgIGZvbnQtZmFtaWx5OiA7IC8qIFRoZSB0b3RhbCB3aWR0aCBpcyA3
MDBweCwgd2hpY2ggd2UgZ2V0IGJ5IGhhdmluZyBwYWRkaW5nIG9mIDUwcHggb24gZWFjaCBsZWZ0
IGFuZCByaWdodCBzaWRlIGFuZCBhIDYwMCB3aWR0aCBzZXR0aW5nICovDQogICAgICAgIGJhY2tn
cm91bmQtY29sb3I6ICNFNkU2RTY7DQogICAgfQ0KDQogICAgLk1haW5WZXJ0aWNhbEJvcmRlcg0K
ICAgIHsNCiAgICAgICAgd2lkdGg6IDUwcHg7DQogICAgfQ0KICAgIC5NYWluSG9yaXpvbnRhbEJv
cmRlcg0KICAgIHsNCiAgICAgICAgaGVpZ2h0OiAyMHB4Ow0KICAgIH0NCiAgICANCiAgICAuSW5u
ZXJIb3Jpem9udGFsQm9yZGVyDQogICAgew0KICAgICAgICBoZWlnaHQ6IDIwcHg7DQogICAgfQ0K
ICAgIC5Jbm5lclZlcnRpY2FsQm9yZGVyDQogICAgew0KICAgICAgICB3aWR0aDogNDBweDsNCiAg
ICAgICAgYmFja2dyb3VuZC1jb2xvcjogI0ZGRkZGRjsgICAgICAgIA0KICAgIH0NCiAgICAuQ29u
dGVudENvbnRhaW5lcg0KICAgIHsNCiAgICAgICAgd2lkdGg6IDUyMHB4Ow0KICAgICAgICBiYWNr
Z3JvdW5kLWNvbG9yOiAjRkZGRkZGOyAgICAgICAgDQogICAgfQ0KICAgIA0KICAgIC5Jbm5lclZl
cnRpY2FsQm9yZGVyV2lkdGgNCiAgICB7DQogICAgICAgIHdpZHRoOiA0MHB4Ow0KICAgIH0NCiAg
ICAuQ29udGVudFdpZHRoDQogICAgew0KICAgICAgICB3aWR0aDogNTIwcHg7DQogICAgfQ0KICAg
IA0KICAgIC5Gb290ZXJDb250YWluZXINCiAgICB7DQogICAgICAgIHRleHQtYWxpZ246IHJpZ2h0
Ow0KICAgIH0NCjwvc3R5bGU+DQoNCg0KPHRhYmxlIGNlbGxwYWRkaW5nPSIwIiBjZWxsc3BhY2lu
Zz0iMCIgY2xhc3M9Ik1haW5Db250YWluZXIiPg0KICAgIDx0ciBjbGFzcz0iTWFpbkhvcml6b250
YWxCb3JkZXIiPg0KICAgICAgICA8dGQgY2xhc3M9Ik1haW5WZXJ0aWNhbEJvcmRlciI+PC90ZD4g
ICAgDQogICAgICAgIDx0ZCBjbGFzcz0iSW5uZXJWZXJ0aWNhbEJvcmRlcldpZHRoIj48L3RkPg0K
ICAgICAgICA8dGQgY2xhc3M9IkNvbnRlbnRXaWR0aCI+PC90ZD4NCiAgICAgICAgPHRkIGNsYXNz
PSJJbm5lclZlcnRpY2FsQm9yZGVyV2lkdGgiPjwvdGQ+DQogICAgICAgIDx0ZCBjbGFzcz0iTWFp
blZlcnRpY2FsQm9yZGVyIj48L3RkPiAgICAgICAgDQogICAgPC90cj4NCiAgICA8dHIgY2xhc3M9
IklubmVySG9yaXpvbnRhbEJvcmRlciI+DQogICAgICAgIDx0ZCBjbGFzcz0iTWFpblZlcnRpY2Fs
Qm9yZGVyIj48L3RkPg0KICAgICAgICA8dGQgY2xhc3M9IklubmVyVmVydGljYWxCb3JkZXIiPjwv
dGQ+DQogICAgICAgIDx0ZCBjbGFzcz0iQ29udGVudENvbnRhaW5lciI+PC90ZD4NCiAgICAgICAg
PHRkIGNsYXNzPSJJbm5lclZlcnRpY2FsQm9yZGVyIj48L3RkPg0KICAgICAgICA8dGQgY2xhc3M9
Ik1haW5WZXJ0aWNhbEJvcmRlciI+PC90ZD4NCiAgICA8L3RyPg0KICAgIDx0cj4NCiAgICAgICAg
PHRkIGNsYXNzPSJNYWluVmVydGljYWxCb3JkZXIiPjwvdGQ+DQogICAgICAgIDx0ZCBjbGFzcz0i
SW5uZXJWZXJ0aWNhbEJvcmRlciI+PC90ZD4NCiAgICAgICAgPHRkIGNsYXNzPSJDb250ZW50Q29u
dGFpbmVyIj4NCiAgICAgICAgICAgIDxkaXYgY2xhc3M9IkNvbnRlbnRDb250YWluZXIiPg0KICAg
ICAgICAgICAgICAgIA0KDQo8c3R5bGUgdHlwZT0idGV4dC9jc3MiPg0KICAgIA0KICAgIC5NZWV0
aW5nUmVxdWVzdEhlYWRlcg0KICAgIHsNCiAgICAgICAgZm9udC1zaXplOiAyMnB4Ow0KICAgICAg
ICBmb250LXdlaWdodDpib2xkOw0KICAgICAgICBjb2xvcjojNDQ0NDQ0Ow0KICAgIH0NCiAgICAu
TWVldGluZ1JlcXVlc3RNZXNzYWdlQ29udGFpbmVyDQogICAgew0KICAgICAgICBwYWRkaW5nLXRv
cDoxNnB4Ow0KICAgIH0NCiAgICAuTWVldGluZ1JlcXVlc3RNZXNzYWdlDQogICAgew0KICAgICAg
ICBmb250LXNpemU6IDIycHg7DQogICAgICAgIGNvbG9yOiNGNDc5M0E7DQogICAgfQ0KICAgIC5N
ZWV0aW5nUmVxdWVzdFF1b3RlDQogICAgew0KICAgICAgICBmb250LWZhbWlseTogOw0KICAgICAg
ICBmb250LXdlaWdodDpib2xkOw0KICAgICAgICBmb250LXNpemU6MjRwdDsNCiAgICAgICAgY29s
b3I6Izg4ODg4ODsNCiAgICB9DQogICAgLk1lZXRpbmdSZXF1ZXN0RGVzY3JpcHRpb24NCiAgICB7
DQogICAgICAgIGZvbnQtZmFtaWx5OiA7DQogICAgICAgIGNvbG9yOiM0NDQ0NDQ7DQogICAgICAg
IGZvbnQtc2l6ZToxM3B4Ow0KICAgIH0NCiAgICAuTWVldGluZ1JlcXVlc3RIUnVsZQ0KICAgIHsN
CiAgICAgICAgYmFja2dyb3VuZC1jb2xvcjogI0VCRUJFQjsNCiAgICAgICAgZm9udC1zaXplOiAx
cHg7DQogICAgICAgIGhlaWdodDoxcHg7DQogICAgICAgIHdpZHRoOjEwMCU7DQogICAgICAgIG1h
cmdpbjoxNnB4IDBweDsNCiAgICB9DQogICAgLk1lZXRpbmdSZXF1ZXN0VGFibGUNCiAgICB7DQog
ICAgICAgIHdpZHRoOjEwMCU7DQogICAgICAgIGJvcmRlci1jb2xsYXBzZTpjb2xsYXBzZTsNCiAg
ICB9DQogICAgLk1lZXRpbmdSZXF1ZXN0VGFibGUgVEQNCiAgICB7DQogICAgICAgIHBhZGRpbmct
dG9wOjE2cHg7DQogICAgfQ0KICAgIC5NZWV0aW5nUmVxdWVzdFRpbWVMb2NhdGlvbkNvbnRhaW5l
cg0KICAgIHsNCiAgICAgICAgZm9udC1zaXplOjE2cHg7DQogICAgICAgIGNvbG9yOiM4ODg4ODg7
DQogICAgICAgIHdpZHRoOjEwMCU7DQogICAgfQ0KICAgIC5NZWV0aW5nUmVxdWVzdENhbmNlbA0K
ICAgIHsNCiAgICAgICAgaGVpZ2h0OjQ4cHg7DQogICAgICAgIHdpZHRoOjQ4cHg7DQogICAgICAg
IG1hcmdpbi1yaWdodDoxMnB4Ow0KICAgIH0NCiAgICANCjwvc3R5bGU+DQoNCjxkaXYgY2xhc3M9
Ik1lZXRpbmdSZXF1ZXN0SGVhZGVyIj5MZWUgR29vZ29sIExlZSDlkJHkvaDlj5HpgIHkuobigJwm
IzI3OTYzOyYjMjExNjA74oCd55qE6YKA6K+3PC9kaXY+DQo8dGFibGUgY2xhc3M9Ik1lZXRpbmdS
ZXF1ZXN0VGFibGUiPg0KICAgIDx0cj4NCiAgICA8dGQgY2xhc3M9Ik1lZXRpbmdSZXF1ZXN0VGlt
ZUxvY2F0aW9uQ29udGFpbmVyIj4NCiAgICAgICAgPGRpdj4yMDEz5bm0M+aciDE55pelPC9kaXY+
PGRpdj7kuIrljYg5OjAwIC0g5LiK5Y2IMTA6MDA8L2Rpdj4NCiAgICAgICAgPGRpdj4mIzIyMzIw
OyYjMjg4NTc7PC9kaXY+DQogICAgPC90ZD4NCiAgICA8L3RyPg0KPC90YWJsZT4NCg0KDQo8dGFi
bGUgY2xhc3M9Ik1lZXRpbmdSZXF1ZXN0VGFibGUiPg0KICAgIDx0cj4NCiAgICA8dGQgY2xhc3M9
Ik1lZXRpbmdSZXF1ZXN0VGltZUxvY2F0aW9uQ29udGFpbmVyIj48ZGl2PuatpOa0u+WKqOWPkeeU
n+S6jiAoR01UKzA4OjAwKSDljJfkuqzvvIzph43luobvvIzpppnmuK/nibnliKvooYzmlL/ljLrv
vIzkuYzpsoHmnKjpvZA8L2Rpdj48L3RkPg0KICAgIDwvdHI+DQo8L3RhYmxlPg0KPGRpdiBjbGFz
cz0iTWVldGluZ1JlcXVlc3RIUnVsZSI+PC9kaXY+DQoNCiAgICAgICAgICAgIDwvZGl2Pg0KICAg
ICAgICAgICAgPGJyIC8+DQogICAgICAgICAgICA8ZGl2IGNsYXNzPSJGb290ZXJDb250YWluZXIi
Pg0KICAgICAgICAgICAgICAgIDxpbWcgc3JjPSJodHRwczovL2dmeDUuaG90bWFpbC5jb20vY2Fs
LzExLjAwL3VwZGF0ZWJldGEvbHRyL2xvZ29fd2xfaG90bWFpbF8xMjAuZ2lmIiBhbHQ9IldpbmRv
d3MgTGl2ZSI+PC9pbWc+IA0KICAgICAgICAgICAgPC9kaXY+DQogICAgICAgIDwvdGQ+DQogICAg
ICAgIDx0ZCBjbGFzcz0iSW5uZXJWZXJ0aWNhbEJvcmRlciI+PC90ZD4NCiAgICAgICAgPHRkIGNs
YXNzPSJNYWluVmVydGljYWxCb3JkZXIiPjwvdGQ+DQogICAgPC90cj4NCiAgICA8dHIgY2xhc3M9
IklubmVySG9yaXpvbnRhbEJvcmRlciI+DQogICAgICAgIDx0ZCBjbGFzcz0iTWFpblZlcnRpY2Fs
Qm9yZGVyIj48L3RkPg0KICAgICAgICA8dGQgY2xhc3M9IklubmVyVmVydGljYWxCb3JkZXIiPjwv
dGQ+DQogICAgICAgIDx0ZCBjbGFzcz0iQ29udGVudENvbnRhaW5lciI+PC90ZD4NCiAgICAgICAg
PHRkIGNsYXNzPSJJbm5lclZlcnRpY2FsQm9yZGVyIj48L3RkPg0KICAgICAgICA8dGQgY2xhc3M9
Ik1haW5WZXJ0aWNhbEJvcmRlciI+PC90ZD4NCiAgICA8L3RyPg0KICAgIDx0ciBjbGFzcz0iTWFp
bkhvcml6b250YWxCb3JkZXIiPg0KICAgICAgICA8dGQgY2xhc3M9Ik1haW5WZXJ0aWNhbEJvcmRl
ciI+PC90ZD4gICAgDQogICAgICAgIDx0ZCBjbGFzcz0iSW5uZXJWZXJ0aWNhbEJvcmRlcldpZHRo
Ij48L3RkPg0KICAgICAgICA8dGQgY2xhc3M9IkNvbnRlbnRXaWR0aCI+PC90ZD4NCiAgICAgICAg
PHRkIGNsYXNzPSJJbm5lclZlcnRpY2FsQm9yZGVyV2lkdGgiPjwvdGQ+DQogICAgICAgIDx0ZCBj
bGFzcz0iTWFpblZlcnRpY2FsQm9yZGVyIj48L3RkPiAgICAgICAgDQogICAgPC90cj4NCjwvdGFi
bGU+DQo=

--_32936406-5bb3-4fc0-b17a-11cc15875002_
Content-Type: text/calendar; charset="utf-8"; method=REQUEST
Content-Transfer-Encoding: base64

QkVHSU46VkNBTEVOREFSDQpNRVRIT0Q6UkVRVUVTVA0KVkVSU0lPTjoyLjANClBST0RJRDotLy9N
aWNyb3NvZnQgQ29ycG9yYXRpb24vL1dpbmRvd3MgTGl2ZSBDYWxlbmRhci8vRU4NCkJFR0lOOlZU
SU1FWk9ORQ0KVFpJRDpDaGluYSBTdGFuZGFyZCBUaW1lDQpCRUdJTjpTVEFOREFSRA0KRFRTVEFS
VDoyMDA4MDEwMVQwMDAwMDANClRaT0ZGU0VUVE86KzA4MDANClRaT0ZGU0VURlJPTTorMDgwMA0K
RU5EOlNUQU5EQVJEDQpFTkQ6VlRJTUVaT05FDQpCRUdJTjpWRVZFTlQNClVJRDpiZTcxYWUxNS01
NjY0LTQ1MWEtODRiMy1hMzFjM2Y1NWFjZmMNCkRUU1RBTVA6MjAxMzAzMThUMDYxNTMzWg0KQ0xB
U1M6UFVCTElDDQpYLU1JQ1JPU09GVC1DRE8tQlVTWVNUQVRVUzpCVVNZDQpUUkFOU1A6T1BBUVVF
DQpTRVFVRU5DRTowDQpEVFNUQVJUO1RaSUQ9Q2hpbmEgU3RhbmRhcmQgVGltZToyMDEzMDMxOVQw
OTAwMDANCkRURU5EO1RaSUQ9Q2hpbmEgU3RhbmRhcmQgVGltZToyMDEzMDMxOVQxMDAwMDANClNV
TU1BUlk65rS75YqoDQpMT0NBVElPTjrlnLDngrkNClBSSU9SSVRZOjANCkFUVEVOREVFO0NVVFlQ
RT1JTkRJVklEVUFMO1JPTEU9UkVRLVBBUlRJQ0lQQU5UO1BBUlRTVEFUPU5FRURTLUFDVElPTjtS
U1ZQPQ0KIFRSVUU6TUFJTFRPOmdvb2dvbGxlZUBnbWFpbC5jb20NCkFUVEVOREVFO0NVVFlQRT1J
TkRJVklEVUFMO1JPTEU9UkVRLVBBUlRJQ0lQQU5UO1BBUlRTVEFUPU5FRURTLUFDVElPTjtSU1ZQ
PQ0KIFRSVUU6TUFJTFRPOnBhbmRhQDBkMGYuY29tDQpPUkdBTklaRVI7Q049TGVlIEdvb2dvbCBM
ZWU6TUFJTFRPOmdvb2dvbGxlZUBob3RtYWlsLmNvbQ0KQkVHSU46VkFMQVJNDQpBQ1RJT046RElT
UExBWQ0KVFJJR0dFUjotUFQxNU0NCkVORDpWQUxBUk0NCkJFR0lOOlZBTEFSTQ0KQUNUSU9OOkRJ
U1BMQVkNClRSSUdHRVI6LVBUMTVNDQpFTkQ6VkFMQVJNDQpFTkQ6VkVWRU5UDQpFTkQ6VkNBTEVO
REFSDQo=

--_32936406-5bb3-4fc0-b17a-11cc15875002_--`

	buf := bytes.NewBufferString(str)
	msg, err := mail.ReadMessage(buf)
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewParser(msg, &config)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, parser.messageID, "BAY152-ds11D2686FF50F86678FFD1AA0E80@phx.gbl")
}

func TestIcsCreate(t *testing.T) {
	var str = `Received: from googollee$163.com ( [114.92.187.87] ) by
 ajax-webmail-wmsvr129 (Coremail) ; Wed, 20 Mar 2013 16:35:41 +0800 (CST)
X-Originating-IP: [114.92.187.87]
Date: Wed, 20 Mar 2013 16:35:41 +0800 (CST)
From: googollee <googollee@163.com>
To: "x@0d0f.com" <x@0d0f.com>, =?GBK?B?1de6oyDA7g==?= <googollee@hotmail.com>
Subject: =?GBK?Q?=D0=C2event?=
X-Priority: 3
X-Mailer: Coremail Webmail Server Version SP_ntes V3.5 build
 20130201(21528.5249.5248) Copyright (c) 2002-2013 www.mailtech.cn 163com
X-CM-CTRLDATA: kqKFxGZvb3Rlcl9odG09OTE6ODE=
Content-Type: multipart/mixed; 
	boundary="----=_Part_534281_720845632.1363768541461"
MIME-Version: 1.0
Message-ID: <1b4b8508.23889.13d86f08116.Coremail.googollee@163.com>

------=_Part_534281_720845632.1363768541461
Content-Type: multipart/alternative; 
	boundary="----=_Part_534283_55711688.1363768541461"

------=_Part_534283_55711688.1363768541461
Content-Type: text/plain; charset=GBK
Content-Transfer-Encoding: base64

suLK1NPKvP5pY3O0tL2o
------=_Part_534283_55711688.1363768541461
Content-Type: text/html; charset=GBK
Content-Transfer-Encoding: base64

PGRpdiBzdHlsZT0ibGluZS1oZWlnaHQ6MS43O2NvbG9yOiMwMDAwMDA7Zm9udC1zaXplOjE0cHg7
Zm9udC1mYW1pbHk6YXJpYWwiPrLiytTTyrz+aWNztLS9qDwvZGl2Pjxicj48YnI+PHNwYW4gdGl0
bGU9Im5ldGVhc2Vmb290ZXIiPjxzcGFuIGlkPSJuZXRlYXNlX21haWxfZm9vdGVyIj48L3NwYW4+
PC9zcGFuPg==
------=_Part_534283_55711688.1363768541461--

------=_Part_534281_720845632.1363768541461
Content-Type: application/octet-stream; name="=?GBK?Q?=D0=C2event.ics?="
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="=?GBK?Q?=D0=C2event.ics?="

QkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL0FwcGxlIEluYy4vL01hYyBP
UyBYIDEwLjguMy8vRU4NCkNBTFNDQUxFOkdSRUdPUklBTg0KQkVHSU46VkVWRU5UDQpDUkVBVEVE
OjIwMTMwMzIwVDA4MjMyNVoNClVJRDoyODJGMDgxNi1ENzdGLTQ5OUUtQTVFMi1BNDAxM0QxRDEz
REMNCkRURU5EO1ZBTFVFPURBVEU6MjAxMzAzMjINClRSQU5TUDpUUkFOU1BBUkVOVA0KU1VNTUFS
WTrmlrBldmVudA0KRFRTVEFSVDtWQUxVRT1EQVRFOjIwMTMwMzIxDQpEVFNUQU1QOjIwMTMwMzIw
VDA4MjM0M1oNCkxPQ0FUSU9OOuWcsOeCuQ0KU0VRVUVOQ0U6Mw0KQkVHSU46VkFMQVJNDQpYLVdS
LUFMQVJNVUlEOjk3QUJEREVELTA5NDktNDZBNy04NDkzLUQ3OENEMDA4NjEyOA0KVUlEOjk3QUJE
REVELTA5NDktNDZBNy04NDkzLUQ3OENEMDA4NjEyOA0KVFJJR0dFUjotUFQxNUgNClgtQVBQTEUt
REVGQVVMVC1BTEFSTTpUUlVFDQpBVFRBQ0g7VkFMVUU9VVJJOkJhc3NvDQpBQ1RJT046QVVESU8N
CkVORDpWQUxBUk0NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==
------=_Part_534281_720845632.1363768541461--`

	buf := bytes.NewBufferString(str)
	msg, err := mail.ReadMessage(buf)
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewParser(msg, &config)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, parser.GetCross().Exfee.Invitations[0].Identity.ExternalID, "googollee@163.com")
	assert.Equal(t, parser.GetCross().Exfee.Invitations[1].Identity.ExternalID, "googollee@hotmail.com")
}

func TestCrossIDType(t *testing.T) {
	var str = `Delivered-To: x+259@0d0f.com
Received: by 10.229.117.1 with SMTP id o1csp67210qcq;
        Wed, 20 Mar 2013 03:21:45 -0700 (PDT)
X-Received: by 10.66.116.138 with SMTP id jw10mr8192933pab.154.1363774905182;
        Wed, 20 Mar 2013 03:21:45 -0700 (PDT)
Return-Path: <googollee@gmail.com>
Received: from mail-pd0-f170.google.com (mail-pd0-f170.google.com [209.85.192.170])
        by mx.google.com with ESMTPS id ui2si1456826pab.146.2013.03.20.03.21.44
        (version=TLSv1 cipher=ECDHE-RSA-RC4-SHA bits=128/128);
        Wed, 20 Mar 2013 03:21:45 -0700 (PDT)
Received-SPF: pass (google.com: domain of googollee@gmail.com designates 209.85.192.170 as permitted sender) client-ip=209.85.192.170;
Authentication-Results: mx.google.com;
       spf=pass (google.com: domain of googollee@gmail.com designates 209.85.192.170 as permitted sender) smtp.mail=googollee@gmail.com;
       dkim=pass header.i=@gmail.com
Received: by mail-pd0-f170.google.com with SMTP id 4so552277pdd.15
        for <x+259@0d0f.com>; Wed, 20 Mar 2013 03:21:44 -0700 (PDT)
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
        d=gmail.com; s=20120113;
        h=x-received:mime-version:in-reply-to:references:from:date:message-id
         :subject:to:content-type;
        bh=nhHTx5UwQDfrnsHwzXP33DaIicFlsER04l+m55Qfw+Y=;
        b=DjOblLucTJqubTyRDVbL4dTIDfQtZ2ZBKvb5Rlxmyo+jr5WUQP6+sg1DjUwAMJfO+2
         81F14XeyarknCi162o2pfFpgwmafvb/oh0+ErUxYmZxOdItLa+k/i0nBPgZH9RJVe80S
         /4uEBtS/oGbiKbgMvXQnZknwM9ajusmxtrPwvxbovfg+pAR6cQgrz/7KcTajemhoWgHe
         7rXTII6IZlWQArgPrXXN7TsAcwVyxIMNASG4GjATXy8Vga7YqGhwhoT1hoSiNnX7Atmb
         Z5MFPYzFCtpuJKT7W8JB2w6zdM9t16sQu87ZF8ySPSEMwP9vcrul321wAdeCtvSrV2nd
         fsAw==
X-Received: by 10.66.14.232 with SMTP id s8mr8364717pac.13.1363774904509; Wed,
 20 Mar 2013 03:21:44 -0700 (PDT)
MIME-Version: 1.0
Received: by 10.70.51.106 with HTTP; Wed, 20 Mar 2013 03:21:24 -0700 (PDT)
In-Reply-To: <51498d74.a1c1420a.4a65.ffff87cc@mx.google.com>
References: <x+259@exfe.com> <51498d74.a1c1420a.4a65.ffff87cc@mx.google.com>
From: Googol Lee <googollee@gmail.com>
Date: Wed, 20 Mar 2013 18:21:24 +0800
Message-ID: <CAOf82vNTO4nXDAFO+1q7qu-9G1YbSeG67xSeZAxMYwtzLNWt_A@mail.gmail.com>
Subject: =?UTF-8?B?UmU6IEludml0YXRpb246IOWGjeasoWdvb2dsZSBjYWxlbmRhciBAIEZyaSBNYXIgMjIsIA==?=
	=?UTF-8?B?MjAxMyAoeEAwZDBmLmNvbSk=?=
To: x+259@0d0f.com
Content-Type: multipart/mixed; boundary=bcaec520eef1bcacb104d858964e

--bcaec520eef1bcacb104d858964e
Content-Type: multipart/alternative; boundary=bcaec520eef1bcacad04d858964c

--bcaec520eef1bcacad04d858964c
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

=E6=B5=8B=E8=AF=95ics=E6=9B=B4=E6=96=B0


2013/3/20 [DEV]EXFE =C2=B7X=C2=B7 <x@0d0f.com>

> **
>   Updates of =C2=B7X=C2=B7 by Googol Lee, googollee.<http://0d0f.com/#!to=
ken=3D70502fb3ee6e7e2ca359cff4c60d3fc1>   Invitation: =E5=86=8D=E6=AC=A1goo=
gle calendar @ Fri Mar 22, 2013 (x@0d0f.com)
> <http://0d0f.com/#!token=3D70502fb3ee6e7e2ca359cff4c60d3fc1>
>
> Fri, Mar 22(+08:00 CST)<http://0d0f.com/#!token=3D70502fb3ee6e7e2ca359cff=
4c60d3fc1>
>
> =E5=9C=B0=E7=82=B9 <http://0d0f.com/#!token=3D70502fb3ee6e7e2ca359cff4c60=
d3fc1>
>
>  <http://0d0f.com/#!token=3D70502fb3ee6e7e2ca359cff4c60d3fc1>
>         [image: googollee] [image: Googol Lee]   <http://0d0f.com/#!token=
=3D70502fb3ee6e7e2ca359cff4c60d3fc1>   Reply
> this email directly as conversation, or try EXFE<http://itunes.apple.com/=
ua/app/exfe/id514026604>app.
> This update is sent from EXFE <http://0d0f.com> automatically.
> Unsubscribe?<http://0d0f.com/mute/cross?token=3D70502fb3ee6e7e2ca359cff4c=
60d3fc1>
>



--=20
=E6=96=B0=E7=9A=84=E7=90=86=E8=AE=BA=E4=BB=8E=E5=B0=91=E6=95=B0=E4=BA=BA=E7=
=9A=84=E4=B8=BB=E5=BC=A0=E5=88=B0=E4=B8=80=E7=BB=9F=E5=A4=A9=E4=B8=8B=EF=BC=
=8C=E5=B9=B6=E4=B8=8D=E6=98=AF=E5=9B=A0=E4=B8=BA=E8=BF=99=E4=B8=AA=E7=90=86=
=E8=AE=BA=E8=AF=B4=E6=9C=8D=E4=BA=86=E5=88=AB=E4=BA=BA=E6=8A=9B=E5=BC=83=E6=
=97=A7=E8=A7=82=E7=82=B9=EF=BC=8C=E8=80=8C=E6=98=AF=E5=9B=A0=E4=B8=BA=E4=B8=
=80=E4=BB=A3=E4=BA=BA=E7=9A=84=E9=80=9D=E5=8E=BB=E3=80=82

--bcaec520eef1bcacad04d858964c
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

<div dir=3D"ltr">=E6=B5=8B=E8=AF=95ics=E6=9B=B4=E6=96=B0</div><div class=3D=
"gmail_extra"><br><br><div class=3D"gmail_quote">2013/3/20 [DEV]EXFE =C2=B7=
X=C2=B7 <span dir=3D"ltr">&lt;<a href=3D"mailto:x@0d0f.com" target=3D"_blan=
k">x@0d0f.com</a>&gt;</span><br><blockquote class=3D"gmail_quote" style=3D"=
margin:0 0 0 .8ex;border-left:1px #ccc solid;padding-left:1ex">

<u></u>

   =20
       =20
       =20
       =20
   =20
    <div>
        <table border=3D"0" cellpadding=3D"0" cellspacing=3D"0" style=3D"fo=
nt-family:Verdana;font-size:13px;line-height:20px;color:#191919;font-weight=
:normal;width:640px;padding:20px;background-color:#fbfbfb">
            <tbody><tr>
                <td colspan=3D"5" style=3D"color:#333333">
                    <a href=3D"http://0d0f.com/#!token=3D70502fb3ee6e7e2ca3=
59cff4c60d3fc1" style=3D"color:#333333;text-decoration:none" target=3D"_bla=
nk">Updates of <span style=3D"color:#3a6ea5">=C2=B7X=C2=B7</span> by Googol=
 Lee, googollee.</a>
                </td>
            </tr>
            <tr><td colspan=3D"5" height=3D"10"></td></tr>
            <tr>
                <td colspan=3D"5" style=3D"font-size:20px;line-height:26px"=
>
                    <a href=3D"http://0d0f.com/#!token=3D70502fb3ee6e7e2ca3=
59cff4c60d3fc1" style=3D"color:#333333;text-decoration:none;font-weight:lig=
hter" target=3D"_blank">
                        Invitation: =E5=86=8D=E6=AC=A1google calendar @ Fri=
 Mar 22, 2013 (x@0d0f.com)
                    </a>
                </td>
            </tr>
            <tr><td colspan=3D"5" height=3D"10"></td></tr>
            <tr>
                <td valign=3D"top" width=3D"180">
                   =20
                    <p style=3D"font-size:20px;line-height:26px;margin:0">
                        <a href=3D"http://0d0f.com/#!token=3D70502fb3ee6e7e=
2ca359cff4c60d3fc1" style=3D"color:#3a6ea5;text-decoration:none" target=3D"=
_blank">Fri, Mar 22(+08:00 CST)</a>
                    </p>
                   =20
                </td>
                <td width=3D"10"></td>
                <td valign=3D"top" width=3D"190" style=3D"word-break:break-=
all">
                   =20
                    <p style=3D"font-size:20px;line-height:26px;margin:0">
                        <a href=3D"http://0d0f.com/#!token=3D70502fb3ee6e7e=
2ca359cff4c60d3fc1" style=3D"color:#333333;text-decoration:none" target=3D"=
_blank">=E5=9C=B0=E7=82=B9</a>
                    </p>
                    <p style=3D"margin:0">
                        <a href=3D"http://0d0f.com/#!token=3D70502fb3ee6e7e=
2ca359cff4c60d3fc1" style=3D"color:#191919;text-decoration:none" target=3D"=
_blank"></a>
                    </p>
                   =20
                </td>
                <td width=3D"10"></td>
                <td valign=3D"top" width=3D"210">
                   =20
                </td>
            </tr>
            <tr><td colspan=3D"5" height=3D"10"></td></tr>
            <tr>
                <td colspan=3D"5">
                    <table border=3D"0" cellpadding=3D"0" cellspacing=3D"0"=
 style=3D"font-family:Verdana;font-size:13px;line-height:20px;color:#191919=
;font-weight:normal;width:100%;background-color:#fbfbfb">
                    =09
                       =20
                       =20
                       =20
                    </table>
                </td>
            </tr>
            <tr><td colspan=3D"5" height=3D"10"></td></tr>
            <tr>
                <td colspan=3D"5">
                   =20
                   =20
                    <img style=3D"padding-right:5px" width=3D"40" height=3D=
"40" alt=3D"googollee" title=3D"googollee" src=3D"http://api.0d0f.com/v2/av=
atar/render?resolution=3D2x&amp;url=3DaHR0cDovL2FwaS4wZDBmLmNvbS92Mi9hdmF0Y=
XIvZGVmYXVsdD9uYW1lPWdvb2dvbGxlZQ%3D%3D&amp;width=3D40&amp;height=3D40&amp;=
alpha=3D0.33">
                   =20
                    <img style=3D"padding-right:5px" width=3D"40" height=3D=
"40" alt=3D"Googol Lee" title=3D"Googol Lee" src=3D"http://api.0d0f.com/v2/=
avatar/render?resolution=3D2x&amp;url=3DaHR0cDovL3d3dy5ncmF2YXRhci5jb20vYXZ=
hdGFyLzE1YjdmYzFiMTAxZWUyODliODE2Nzg4MTI3ODFhZWE2&amp;width=3D40&amp;height=
=3D40&amp;ishost=3Dtrue">
                   =20
                </td>
            </tr>
            <tr><td colspan=3D"5" height=3D"10"></td></tr>
            <tr>
                <td colspan=3D"5">
                    <a href=3D"http://0d0f.com/#!token=3D70502fb3ee6e7e2ca3=
59cff4c60d3fc1" style=3D"color:#333333;text-decoration:none" target=3D"_bla=
nk"></a>
                </td>
            </tr>
            <tr><td colspan=3D"5" height=3D"20"></td></tr>
            <tr>
                <td colspan=3D"5" style=3D"font-size:11px;line-height:15px;=
color:#7f7f7f">
                    Reply this email directly as conversation, or try <a st=
yle=3D"color:#3a6ea5;text-decoration:none" href=3D"http://itunes.apple.com/=
ua/app/exfe/id514026604" target=3D"_blank">EXFE</a> app.
                    <br>
                    <span style=3D"color:#b2b2b2">This update is sent from =
<a style=3D"color:#3a6ea5;text-decoration:none" href=3D"http://0d0f.com" ta=
rget=3D"_blank">EXFE</a> automatically. <a style=3D"color:#e6e6e6;text-deco=
ration:none" href=3D"http://0d0f.com/mute/cross?token=3D70502fb3ee6e7e2ca35=
9cff4c60d3fc1" target=3D"_blank">Unsubscribe?</a>
                   =20
                    </span>
                </td>
            </tr>
        </tbody></table>
    </div>

</blockquote></div><br><br clear=3D"all"><div><br></div>-- <br>=E6=96=B0=E7=
=9A=84=E7=90=86=E8=AE=BA=E4=BB=8E=E5=B0=91=E6=95=B0=E4=BA=BA=E7=9A=84=E4=B8=
=BB=E5=BC=A0=E5=88=B0=E4=B8=80=E7=BB=9F=E5=A4=A9=E4=B8=8B=EF=BC=8C=E5=B9=B6=
=E4=B8=8D=E6=98=AF=E5=9B=A0=E4=B8=BA=E8=BF=99=E4=B8=AA=E7=90=86=E8=AE=BA=E8=
=AF=B4=E6=9C=8D=E4=BA=86=E5=88=AB=E4=BA=BA=E6=8A=9B=E5=BC=83=E6=97=A7=E8=A7=
=82=E7=82=B9=EF=BC=8C=E8=80=8C=E6=98=AF=E5=9B=A0=E4=B8=BA=E4=B8=80=E4=BB=A3=
=E4=BA=BA=E7=9A=84=E9=80=9D=E5=8E=BB=E3=80=82
</div>

--bcaec520eef1bcacad04d858964c--
--bcaec520eef1bcacb104d858964e
Content-Type: application/octet-stream; name="=?UTF-8?B?5pawZXZlbnQuaWNz?="
Content-Disposition: attachment; filename="=?UTF-8?B?5pawZXZlbnQuaWNz?="
Content-Transfer-Encoding: base64
X-Attachment-Id: f_heicccag1

QkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL0FwcGxlIEluYy4vL01hYyBP
UyBYIDEwLjguMy8vRU4NCkNBTFNDQUxFOkdSRUdPUklBTg0KQkVHSU46VkVWRU5UDQpDUkVBVEVE
OjIwMTMwMzIwVDA4MjMyNVoNClVJRDoyODJGMDgxNi1ENzdGLTQ5OUUtQTVFMi1BNDAxM0QxRDEz
REUNCkRURU5EO1ZBTFVFPURBVEU6MjAxMzAzMjINClRSQU5TUDpUUkFOU1BBUkVOVA0KU1VNTUFS
WTrmlrBldmVudA0KRFRTVEFSVDtWQUxVRT1EQVRFOjIwMTMwMzIxDQpEVFNUQU1QOjIwMTMwMzIw
VDA4MjM0M1oNCkxPQ0FUSU9OOuWcsOeCuQ0KU0VRVUVOQ0U6Mw0KQkVHSU46VkFMQVJNDQpYLVdS
LUFMQVJNVUlEOjk3QUJEREVELTA5NDktNDZBNy04NDkzLUQ3OENEMDA4NjEyOA0KVUlEOjk3QUJE
REVELTA5NDktNDZBNy04NDkzLUQ3OENEMDA4NjEyOA0KVFJJR0dFUjotUFQxNUgNClgtQVBQTEUt
REVGQVVMVC1BTEFSTTpUUlVFDQpBVFRBQ0g7VkFMVUU9VVJJOkJhc3NvDQpBQ1RJT046QVVESU8N
CkVORDpWQUxBUk0NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==
--bcaec520eef1bcacb104d858964e--`

	buf := bytes.NewBufferString(str)
	msg, err := mail.ReadMessage(buf)
	if err != nil {
		t.Fatal(err)
	}
	parser, err := NewParser(msg, &config)
	if err != nil {
		t.Fatal(err)
	}
	to, id := parser.GetTypeID()
	assert.Equal(t, to, "cross_id")
	assert.Equal(t, id, "259")
}

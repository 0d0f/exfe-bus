package notifier

import (
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

var post1 = model.Post{
	ID:        1,
	By:        email1,
	Content:   "email1 post sth",
	Via:       "abc",
	CreatedAt: "2012-10-24 16:31:00",
}

var post2 = model.Post{
	ID:        2,
	By:        twitter3,
	Content:   "twitter3 post sth",
	Via:       "abc",
	CreatedAt: "2012-10-24 16:40:00",
}

func TestConversationUpdateToSelf(t *testing.T) {
	update1 := model.ConversationUpdate{
		To:    remail1,
		Cross: cross,
		Post:  post1,
	}
	updates := []model.ConversationUpdate{update1}

	c := NewConversation(localTemplate, &config)
	private, public, err := c.getContent(updates)
	assert.Equal(t, err.Error(), "can't parse posts: no need send self")
	assert.Equal(t, private, "")
	assert.Equal(t, public, "")
}

func TestConversationUpdateTwitter(t *testing.T) {
	update1 := model.ConversationUpdate{
		To:    rtwitter1,
		Cross: cross,
		Post:  post1,
	}
	update2 := model.ConversationUpdate{
		To:    rtwitter1,
		Cross: cross,
		Post:  post2,
	}
	updates := []model.ConversationUpdate{update1, update2}

	expectPrivate := `email1 name: email1 post sth \((“Test Cross” http://site/url/#!token=recipient_twitter1_token)\)
twitter3 name: twitter3 post sth \((“Test Cross” http://site/url/#!token=recipient_twitter1_token)\)
`
	c := NewConversation(localTemplate, &config)
	private, public, err := c.getContent(updates)
	assert.Equal(t, err, nil)
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "")
}

func TestConversationUpdateEmail(t *testing.T) {
	update1 := model.ConversationUpdate{
		To:    remail1,
		Cross: cross,
		Post:  post1,
	}
	update2 := model.ConversationUpdate{
		To:    remail1,
		Cross: cross,
		Post:  post2,
	}
	updates := []model.ConversationUpdate{update1, update2}

	c := NewConversation(localTemplate, &config)
	private, public, err := c.getContent(updates)

	expectPrivate := "Content-Type: multipart/mixed; boundary=\"56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\"\nReferences: <x+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <email1@domain.com>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nwrcgZW1haWwxIG5hbWUgYXQgMTI6MzFBTSBUaHUsIE9jdCAyNSBzYWlkOgogICAgZW1haWwxIHBvc3Qg\r\nc3RoCsK3IHR3aXR0ZXIzIG5hbWUgYXQgMTI6NDBBTSBUaHUsIE9jdCAyNSBzYWlkOgogICAgdHdpdHRl\r\ncjMgcG9zdCBzdGgKCiMgUmVwbHkgdGhpcyBlbWFpbCBkaXJlY3RseSBhcyBjb252ZXJzYXRpb24uICM=\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgICAgIDxzdHlsZT4KICAgICAgICAgICAgLmV4ZmVfbWFpbF9pZGVudGl0\r\neV9uYW1lIHsKICAgICAgICAgICAgICAgIGNvbG9yOiAjM2E2ZWE1OwogICAgICAgICAgICB9CiAgICAg\r\nICAgICAgIC5leGZlX21haWxfbXNnX2lkZW50aXR5X25hbWUgewogICAgICAgICAgICAgICAgY29sb3I6\r\nICM2NjY2NjY7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfYXQgewogICAg\r\nICAgICAgICAgICAgZm9udC1zaXplOiAxMnB4OwogICAgICAgICAgICAgICAgY29sb3I6ICM5OTk5OTk7\r\nCiAgICAgICAgICAgIH0KICAgICAgICAgICAgLmV4ZmVfbWFpbF9tc2dfdGltZSB7CiAgICAgICAgICAg\r\nICAgICBmb250LXNpemU6IDEycHg7CiAgICAgICAgICAgICAgICBjb2xvcjogIzY2NjY2NjsKICAgICAg\r\nICAgICAgfQogICAgICAgIDwvc3R5bGU+CiAgICA8L2hlYWQ+CiAgICA8Ym9keT4KICAgICAgICA8dGFi\r\nbGUgYm9yZGVyPSIwIiBjZWxscGFkZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJmb250LWZh\r\nbWlseTogVmVyZGFuYTsgZm9udC1zaXplOiAxM3B4OyBsaW5lLWhlaWdodDogMjBweDsgY29sb3I6ICMx\r\nOTE5MTk7IGZvbnQtd2VpZ2h0OiBub3JtYWw7IHdpZHRoOiA2NDBweDsgcGFkZGluZzogMjBweDsgYmFj\r\na2dyb3VuZC1jb2xvcjogI2ZiZmJmYjsiPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8\r\ndGQgY29sc3Bhbj0iNSI+CiAgICAgICAgICAgICAgICAgICAgPHRhYmxlIGJvcmRlcj0iMCIgY2VsbHBh\r\nZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIj4KICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAg\r\nICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8dGQgdmFsaWduPSJ0b3AiIHdpZHRoPSI1MCIgYWxpZ249ImxlZnQiPgogICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxpbWcgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiB0aXRs\r\nZT0iZW1haWwxIG5hbWUiIGFsdD0iZW1haWwxIG5hbWUiIHNyYz0iaHR0cDovL3BhdGgvdG8vZW1haWwx\r\nLmF2YXRhciI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICA8YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIiBz\r\ndHlsZT0idGV4dC1kZWNvcmF0aW9uOiBub25lOyI+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDxzcGFuIHN0eWxlPSJjb2xvcjojMjEyMTIxOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5l\r\nbWFpbDEgcG9zdCBzdGg8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDxi\r\nciAvPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8c3BhbiBzdHlsZT0iZm9udC1z\r\naXplOiAxMnB4OyBjb2xvcjojNzk3OTc5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5lbWFpbDEgbmFt\r\nZTwvc3Bhbj4gPHNwYW4gc3R5bGU9ImZvbnQtc2l6ZTogMTJweDsgY29sb3I6I0E5QTlBOTsgdGV4dC1k\r\nZWNvcmF0aW9uOiBub25lOyI+YXQ8L3NwYW4+IDxzcGFuICBzdHlsZT0iZm9udC1zaXplOiAxMnB4OyBj\r\nb2xvcjojNzk3OTc5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij4xMjozMUFNIFRodSwgT2N0IDI1PC9z\r\ncGFuPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDwvYT4KICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgICAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgIDx0cj48dGQgY29sc3Bhbj0iMiIgaGVpZ2h0PSIyMCI+PC90ZD48L3RyPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgICAgICAgICAgICAgPHRkIHZhbGlnbj0idG9wIiB3aWR0aD0iNTAiIGFsaWduPSJsZWZ0Ij4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8aW1nIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgdGl0\r\nbGU9InR3aXR0ZXIzIG5hbWUiIGFsdD0idHdpdHRlcjMgbmFtZSIgc3JjPSJodHRwOi8vcGF0aC90by90\r\nd2l0dGVyMy5hdmF0YXIiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAg\r\nICAgICAgICAgICAgICAgICAgIDx0ZCB2YWxpZ249InRvcCI+CiAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgPGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90\r\nb2tlbiIgc3R5bGU9InRleHQtZGVjb3JhdGlvbjogbm9uZTsiPgogICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgICAgICAgICA8c3BhbiBzdHlsZT0iY29sb3I6IzIxMjEyMTsgdGV4dC1kZWNvcmF0aW9uOiBu\r\nb25lOyI+dHdpdHRlcjMgcG9zdCBzdGg8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\nICAgICAgIDxiciAvPgogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA8c3BhbiBzdHls\r\nZT0iZm9udC1zaXplOiAxMnB4OyBjb2xvcjojNzk3OTc5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij50\r\nd2l0dGVyMyBuYW1lPC9zcGFuPiA8c3BhbiBzdHlsZT0iZm9udC1zaXplOiAxMnB4OyBjb2xvcjojQTlB\r\nOUE5OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5hdDwvc3Bhbj4gPHNwYW4gIHN0eWxlPSJmb250LXNp\r\nemU6IDEycHg7IGNvbG9yOiM3OTc5Nzk7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPjEyOjQwQU0gVGh1\r\nLCBPY3QgMjU8L3NwYW4+CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgPC9hPgogICAgICAg\r\nICAgICAgICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgICAgICAgICAgICAgPC90cj4KICAg\r\nICAgICAgICAgICAgICAgICAgICAgPHRyPjx0ZCBjb2xzcGFuPSIyIiBoZWlnaHQ9IjIwIj48L3RkPjwv\r\ndHI+CiAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICAgICAgIDwvdGFibGU+CiAg\r\nICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGNv\r\nbHNwYW49IjUiIGhlaWdodD0iMjAiPjwvdGQ+PC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAg\r\nICAgICAgPHRkIGNvbHNwYW49IjUiIHN0eWxlPSJmb250LXNpemU6IDExcHg7IGxpbmUtaGVpZ2h0OiAx\r\nNXB4OyBjb2xvcjogIzdGN0Y3RjsiPgogICAgICAgICAgICAgICAgICAgIFJlcGx5IHRoaXMgZW1haWwg\r\nZGlyZWN0bHkgYXMgY29udmVyc2F0aW9uLjwhLS0sIG9yIHRyeSA8YSBzdHlsZT0iY29sb3I6ICMzYTZl\r\nYTU7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiIGhyZWY9Imh0dHA6Ly9hcHAvdXJsIj5FWEZFPC9hPiBh\r\ncHAuLS0+CiAgICAgICAgICAgICAgICAgICAgPGJyIC8+CiAgICAgICAgICAgICAgICAgICAgPHNwYW4g\r\nc3R5bGU9ImNvbG9yOiAjQjJCMkIyIj5UaGlzIHVwZGF0ZSBpcyBzZW50IGZyb20gPGEgc3R5bGU9ImNv\r\nbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwi\r\nPkVYRkU8L2E+IGF1dG9tYXRpY2FsbHkuIDxhIHN0eWxlPSJjb2xvcjogI0U2RTZFNjsgdGV4dC1kZWNv\r\ncmF0aW9uOiBub25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsL3MvcmVwb3J0U3BhbT90b2tlbj1yZWNp\r\ncGllbnRfZW1haWwxX3Rva2VuIj5VbnN1YnNjcmliZT88L2E+CiAgICAgICAgICAgICAgICAgICAgPCEt\r\nLQogICAgICAgICAgICAgICAgICAgIFlvdSBjYW4gY2hhbmdlCiAgICAgICAgICAgICAgICAgICAgPGEg\r\nc3R5bGU9ImNvbG9yOiAjQjJCMkIyOyB0ZXh0LWRlY29yYXRpb246IHVuZGVsaW5lOyIgaHJlZj0iIj5u\r\nb3RpZmljYXRpb24gcHJlZmVyZW5jZTwvYT4uCiAgICAgICAgICAgICAgICAgICAgLS0+CiAgICAgICAg\r\nICAgICAgICAgICAgPC9zcGFuPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4K\r\nICAgICAgICA8L3RhYmxlPgogICAgPC9ib2R5Pgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75\nContent-Disposition: attachment; filename=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Transfer-Encoding: base64\n\nQkVHSU46VkNBTEVOREFSDQpWRVJTSU9OOjIuMA0KUFJPRElEOi0vL2V4ZmUvL2V4ZmUuY29tIC8vDQpY\r\nLVdSLUNBTE5BTUU6VGVzdCBDcm9zcw0KWC1XUi1DQUxERVNDOmV4ZmUgY2FsDQpYLVdSLVRJTUVaT05F\r\nOiswODAwDQpCRUdJTjpWRVZFTlQNClVJRDohMTIzQGV4ZmUNCkRUU1RBTVA6DQpERVNDUklQVElPTjp0\r\nZXN0IGNyb3NzIGRlc2NyaXB0aW9uDQpEVFNUQVJUOlZBTFVFPURBVEU6DQpMT0NBVElPTjoNClNVTU1B\r\nUlk6VGVzdCBDcm9zcw0KVVJMOmh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lwaWVudF9lbWFpbDFf\r\ndG9rZW4NCkVORDpWRVZFTlQNCkVORDpWQ0FMRU5EQVINCg==\n--56040bc4f71301a3dc363b960b1796dafbb8b190894fd231dda878b5be75--\n"
	assert.Equal(t, err, nil)
	t.Logf("-----start------\n%s\n--------end-------", private)
	assert.Equal(t, private, expectPrivate)
	assert.Equal(t, public, "")
}

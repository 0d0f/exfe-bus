package notifier

import (
	"encoding/json"
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

func TestUserWelcomeEmail(t *testing.T) {
	arg := model.UserWelcome{}
	arg.To = remail1
	arg.NeedVerify = true

	err := arg.Parse(&config)
	assert.Equal(t, err, nil)
	content, err := GetContent(localTemplate, "user_welcome", arg.To, arg)
	assert.Equal(t, err, nil)
	t.Logf("content:---------start---------\n%s\n---------end----------", content)
	expectPrivate := "Content-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <email1@domain.com>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x@test.com>\nSubject: =?utf-8?B?V2VsY29tZSB0byBFWEZF?=\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nSGksIGVtYWlsMSBuYW1lLgoKV2VsY29tZSB0byBFWEZFISBBbiB1dGlsaXR5IGZvciBoYW5naW5nIG91\r\ndCB3aXRoIGZyaWVuZHMuCgpQbGVhc2UgY2xpY2sgaGVyZSB0byB2ZXJpZnkgeW91ciBpZGVudGl0eTog\r\naHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbgoKwrdYwrcoY3Jvc3Mp\r\nIGlzIGEgZ2F0aGVyaW5nIG9mIHBlb3BsZSwgb24gcHVycG9zZSBvciBub3QuIFdlIHNhdmUgeW91IGZy\r\nb20gY2FsbGluZyB1cCBldmVyeSBvbmUgUlNWUCwgbG9zaW5nIGluIGVuZGxlc3MgZW1haWxzIG1lc3Nh\r\nZ2VzIG9mZiB0aGUgcG9pbnQuCgpFWEZFIHlvdXIgZnJpZW5kcy4gR2F0aGVyIGEgwrdYwrcKCuKAnFJv\r\nbWUgd2Fzbid0IGJ1aWx0IGluIGEgZGF5LuKAnSBFWEZFIFvLiMmba3NmaV0gaXMgc3RpbGwgaW4gcGls\r\nb3Qgc3RhZ2UuIFdl4oCZcmUgYnVpbGRpbmcgdXAgYmxvY2tzLCBjb25zZXF1ZW50bHkgc29tZSBidWdz\r\nIG9yIHVuZmluaXNoZWQgcGFnZXMgbWF5IGhhcHBlbi4gQW55IGZlZWRiYWNrLCBwbGVhc2UgZW1haWwg\r\ndG8gZmVlZGJhY2tAZXhmZS5jb20uIE11Y2ggYXBwcmVjaWF0ZWQuCgpUaGlzIGlzIHNlbnQgdG8gZW1h\r\naWwxQGRvbWFpbi5jb20gcGVyIGlkZW50aXR5IHJlZ2lzdHJhdGlvbiByZXF1ZXN0IG9uIEVYRkUoIGh0\r\ndHA6Ly9zaXRlL3VybCApLg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgPC9oZWFkPgogICAgPGJvZHk+CiAgICAgICAgPHRhYmxlIGJvcmRlcj0i\r\nMCIgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0iZm9udC1mYW1pbHk6IFZlcmRh\r\nbmE7IGZvbnQtc2l6ZTogMTRweDsgbGluZS1oZWlnaHQ6IDIwcHg7IGNvbG9yOiAjMTkxOTE5OyBmb250\r\nLXdlaWdodDogbm9ybWFsOyB3aWR0aDogNjQwcHg7IHBhZGRpbmc6IDEwcHggMjBweCAzMHB4IDIwcHg7\r\nIGJhY2tncm91bmQtY29sb3I6ICNmYmZiZmI7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgPHRkIHZhbGlnbj0idG9wIj4KICAgICAgICAgICAgICAgICAgICA8cCBzdHlsZT0iZm9udC1zaXpl\r\nOiAxOHB4OyBsaW5lLWhlaWdodDogMjRweDsiPgogICAgICAgICAgICAgICAgICAgIEhpLCBlbWFpbDEg\r\nbmFtZQogICAgICAgICAgICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICA8cD4KICAgICAg\r\nICAgICAgICAgICAgICAgICAgPHNwYW4gc3R5bGU9ImZvbnQtc2l6ZTogMThweDsgbGluZS1oZWlnaHQ6\r\nIDI0cHg7Ij5XZWxjb21lIHRvIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9u\r\nOiBub25lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsIj5FWEZFPC9hPiE8L3NwYW4+CiAgICAgICAgICAg\r\nICAgICAgICAgICAgIDxicj4KICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4+QW4gdXRpbGl0eSBm\r\nb3IgaGFuZ2luZyBvdXQgd2l0aCBmcmllbmRzLjwvc3Bhbj4KICAgICAgICAgICAgICAgICAgICA8L3A+\r\nCiAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgPHA+CiAgICAgICAgICAgICAg\r\nICAgICAgICAgIDxzcGFuIHN0eWxlPSJmb250LXNpemU6MTRweDsgY29sb3I6IzMzMzsgIj5QbGVhc2Ug\r\nY2xpY2sgaGVyZSB0byB2ZXJpZnkgeW91ciBpZGVudGl0eTo8L3NwYW4+IDxhIHN0eWxlPSJjb2xvcjoj\r\nMTkxOTE5OyB0ZXh0LWRlY29yYXRpb246IHVuZGVybGluZTsiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8j\r\nIXRva2VuPXJlY2lwaWVudF9lbWFpbDFfdG9rZW4iPmh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPXJlY2lw\r\n4oCmPC9hPgogICAgICAgICAgICAgICAgICAgIDwvcD4KICAgICAgICAgICAgICAgICAgICAKICAgICAg\r\nICAgICAgICAgICAgICA8cD4KICAgICAgICAgICAgICAgICAgICAgICAgPHNwYW4gc3R5bGU9ImNvbG9y\r\nOiAjM2E2ZWE1OyI+wrdYwrc8L3NwYW4+IChjcm9zcykgaXMgYSBnYXRoZXJpbmcgb2YgcGVvcGxlLCBv\r\nbiBwdXJwb3NlIG9yIG5vdC4gV2Ugc2F2ZSB5b3UgZnJvbSBjYWxsaW5nIHVwIGV2ZXJ5IG9uZSBSU1ZQ\r\nLCBsb3NpbmcgaW4gZW5kbGVzcyBlbWFpbHMgbWVzc2FnZXMgb2ZmIHRoZSBwb2ludC4KICAgICAgICAg\r\nICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAgICAgICAgPHA+CiAgICAgICAgICAgICAgICAgICAg\r\nICAgIDxzcGFuIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsiPkVYRkU8L3NwYW4+IHlvdXIgZnJpZW5kcy4g\r\nR2F0aGVyIGEgPHNwYW4gc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyI+wrdYwrc8L3NwYW4+CiAgICAgICAg\r\nICAgICAgICAgICAgPC9wPgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAg\r\nICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIGhlaWdodD0iMjAiPjwvdGQ+CiAgICAgICAg\r\nICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj4KICAgICAgICAgICAgICAgIDx0ZD4KICAgICAgICAgICAg\r\nICAgICAgICA8aW1nIHN0eWxlPSJmbG9hdDogbGVmdDsgbWFyZ2luLWxlZnQ6IDIwcHg7IG1hcmdpbi1y\r\naWdodDogNDBweDsgbWFyZ2luLWJvdHRvbTogMTBweDsiIHNyYz0iaHR0cDovL3NpdGUvaW1nL2VtYWls\r\nL3JvbWUucG5nIj4KICAgICAgICAgICAgICAgICAgICA8ZGl2IHN0eWxlPSJjb2xvcjojMzMzMzMzOyI+\r\nCiAgICAgICAgICAgICAgICAgICAgICAgIDxwIHN0eWxlPSJmb250LXNpemU6IDE4cHg7IGxpbmUtaGVp\r\nZ2h0OiAyNHB4OyBtYXJnaW46IDAgMCAyMHB4IDAiPgogICAgICAgICAgICAgICAgICAgICAgICAgICAg\r\n4oCcUm9tZSB3YXNuJ3QgYnVpbHQgaW4gYSBkYXku4oCdCiAgICAgICAgICAgICAgICAgICAgICAgIDwv\r\ncD4KICAgICAgICAgICAgICAgICAgICAgICAgPHAgc3R5bGU9Im1hcmdpbjogMDsiPgogICAgICAgICAg\r\nICAgICAgICAgICAgICA8YSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9u\r\nZTsiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybCI+RVhGRTwvYT4gW8uIyZtrc2ZpXSBpcyBzdGlsbCBpbiA8\r\nc3BhbiBzdHlsZSA9ImZvbnQtd2VpZ2h0OiBib2xkOyI+cGlsb3Q8L3NwYW4+IHN0YWdlLiBXZeKAmXJl\r\nIGJ1aWxkaW5nIHVwIGJsb2NrcyBvZiBpdCwgY29uc2VxdWVudGx5IHNvbWUgYnVncyBvciB1bmZpbmlz\r\naGVkIHBhZ2VzIG1heSBoYXBwZW4uIEFueSBmZWVkYmFjaywgcGxlYXNlIGVtYWlsIDxhIHN0eWxlPSJm\r\nb250LXN0eWxlOiBpdGFsaWM7IGNvbG9yOiAjMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IHVuZGVybGlu\r\nZTsiIGhyZWY9Im1haWx0bzpmZWVkYmFja0BleGZlLmNvbSI+ZmVlZGJhY2tAZXhmZS5jb208L2E+LiBN\r\ndWNoIGFwcHJlY2lhdGVkLgogICAgICAgICAgICAgICAgICAgICAgICA8L3A+CiAgICAgICAgICAgICAg\r\nICAgICAgPC9kaXY+CiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAg\r\nICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgaGVpZ2h0PSIyMCI+PC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkIHN0eWxlPSJjb2xvcjogIzY2\r\nNjY2NjsgZm9udC1zaXplOiAxMXB4OyI+CiAgICAgICAgICAgICAgICAgICAgVGhpcyBpcyBzZW50IHRv\r\nIDxhIHN0eWxlPSJmb250LXN0eWxlOiBpdGFsaWM7IGNvbG9yOiAjNjY2NjY2OyB0ZXh0LWRlY29yYXRp\r\nb246IG5vbmU7IiBocmVmPSJtYWlsdG86ZW1haWwxQGRvbWFpbi5jb20iPmVtYWlsMUBkb21haW4uY29t\r\nPC9hPiBwZXIgaWRlbnRpdHkgcmVnaXN0cmF0aW9uIHJlcXVlc3Qgb24gPGEgc3R5bGU9ImNvbG9yOiAj\r\nM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8vc2l0ZS91cmwiPkVYRkU8\r\nL2E+LgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4KICAgICAgICA8L3RhYmxl\r\nPgogICAgPC9ib2R5Pgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n"
	assert.Equal(t, content, expectPrivate)
}

func TestUserWelcomeTwitter(t *testing.T) {
	arg := model.UserWelcome{}
	arg.To = rtwitter1
	arg.NeedVerify = true

	err := arg.Parse(&config)
	assert.Equal(t, err, nil)
	content, err := GetContent(localTemplate, "user_welcome", arg.To, arg)
	assert.Equal(t, err, nil)
	t.Logf("content:---------start---------\n%s\n---------end----------", content)
	expectPrivate := "Welcome to EXFE! An utility for hanging out with friends.Please click here to verify your identity: \\(http://site/url/#!token=recipient_twitter1_token\\)"
	assert.Equal(t, content, expectPrivate)
}

func TestUserConfirmEmail(t *testing.T) {
	arg := model.UserConfirm{}
	arg.To = remail1
	arg.By = email2

	d, _ := json.Marshal(arg)
	fmt.Println(string(d))
	t.Errorf("show")

	err := arg.Parse(&config)
	assert.Equal(t, err, nil)
	content, err := GetContent(localTemplate, "user_confirm", arg.To, arg)
	assert.Equal(t, err, nil)
	t.Logf("content:---------start---------\n%s\n---------end----------", content)
	expectPrivate := "Content-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <email1@domain.com>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x@test.com>\nSubject: =?utf-8?B?RVhGRSBpZGVudGl0eSB2ZXJpZmljYXRpb24=?=\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nSGksIGVtYWlsMSBuYW1lLgoKWW91ciBlbWFpbCBlbWFpbDFAZG9tYWluLmNvbSBoYXMgYmVlbiByZXF1\r\nZXN0ZWQgZm9yIHZlcmlmaWNhdGlvbiBieSBlbWFpbDIgbmFtZSBvbiBFWEZFLgoKUGxlYXNlIGNsaWNr\r\nIGhlcmUgdG8gdmVyaWZ5OiBodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rv\r\na2VuCgpFWEZFIGlzIGEgdXRpbGl0eSBmb3IgaGFuZ2luZyBvdXQgd2l0aCBmcmllbmRzLiBXZSBzYXZl\r\nIHlvdSBmcm9tIGNhbGxpbmcgdXAgZXZlcnkgb25lIFJTVlAsIGxvc2luZyBpbiBlbmRsZXNzIGVtYWls\r\ncyBhbmQgbWVzc2FnZXMgb2ZmIHRoZSBwb2ludC4KClRoaXMgZW1haWwgaXMgc2VudCB0byBlbWFpbDFA\r\nZG9tYWluLmNvbSBwZXIgaWRlbnRpdHkgdmVyaWZpY2F0aW9uIHJlcXVlc3Qgb24gRVhGRSggaHR0cDov\r\nL3NpdGUvdXJsICku\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgPC9oZWFkPgogICAgPGJvZHk+CiAgICAgICAgPHRhYmxlIGJvcmRlcj0i\r\nMCIgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0iZm9udC1mYW1pbHk6IFZlcmRh\r\nbmE7IGZvbnQtc2l6ZTogMTRweDsgbGluZS1oZWlnaHQ6IDIwcHg7IGNvbG9yOiAjMzMzMzMzOyBmb250\r\nLXdlaWdodDogbm9ybWFsOyB3aWR0aDogNjQwcHg7IHBhZGRpbmc6IDIwcHg7IGJhY2tncm91bmQtY29s\r\nb3I6ICNmYmZiZmI7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkPgogICAgICAg\r\nICAgICAgICAgICAgIEhpLCA8c3BhbiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7Ij5lbWFpbDEgbmFtZS48\r\nL3NwYW4+CiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgICAgICA8\r\ndHI+PHRkIGhlaWdodD0iMjAiPjwvdGQ+PC90cj4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAg\r\nICAgPHRkPgogICAgICAgICAgICAgICAgICAgIFlvdXIgZW1haWwgPGEgaHJlZj0ibWFpbHRvOmVtYWls\r\nMUBkb21haW4uY29tIiBzdHlsZT0idGV4dC1kZWNvcmF0aW9uOiBub25lOyBjb2xvcjogIzMzMzMzMzsg\r\nZm9udC1zdHlsZTogaXRhbGljOyI+ZW1haWwxQGRvbWFpbi5jb208L2E+IGhhcyBiZWVuIHJlcXVlc3Rl\r\nZCBmb3IgdmVyaWZpY2F0aW9uIGJ5IDxzcGFuIHN0eWxlPSJjb2xvcjogIzE5MTkxOTsiPmVtYWlsMiBu\r\nYW1lPC9zcGFuPiBvbiA8YSBzdHlsZT0iY29sb3I6ICMzYTZlYTU7IHRleHQtZGVjb3JhdGlvbjogbm9u\r\nZTsiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybCI+RVhGRTwvYT4uCiAgICAgICAgICAgICAgICA8L3RkPgog\r\nICAgICAgICAgICA8L3RyPgogICAgICAgICAgICA8dHI+PHRkIGhlaWdodD0iMjAiPjwvdGQ+PC90cj4K\r\nICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkPgogICAgICAgICAgICAgICAgICAgIDxh\r\nIHN0eWxlPSJjb2xvcjojMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIj5QbGVhc2UgY2xpY2sgaGVyZSB0\r\nbyB2ZXJpZnk6IDxzcGFuIHN0eWxlPSJ0ZXh0LWRlY29yYXRpb246IHVuZGVybGluZTsiPmh0dHA6Ly9z\r\naXRlL3VybC8jIXRva2VuPXJlY2lw4oCmPC9zcGFuPjwvYT4KICAgICAgICAgICAgICAgIDwvdGQ+CiAg\r\nICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgaGVpZ2h0PSIyMCI+PC90ZD48L3RyPgog\r\nICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgPGEg\r\nc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJodHRwOi8v\r\nc2l0ZS91cmwiPkVYRkU8L2E+IGlzIGEgdXRpbGl0eSBmb3IgaGFuZ2luZyBvdXQgd2l0aCBmcmllbmRz\r\nLiBXZSBzYXZlIHlvdSBmcm9tIGNhbGxpbmcgdXAgZXZlcnkgb25lIFJTVlAsIGxvc2luZyBpbiBlbmRs\r\nZXNzIGVtYWlscyBhbmQgbWVzc2FnZXMgb2ZmIHRoZSBwb2ludC4KICAgICAgICAgICAgICAgIDwvdGQ+\r\nCiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgaGVpZ2h0PSI0MCI+PC90ZD48L3Ry\r\nPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQgc3R5bGU9ImZvbnQtc2l6ZTogMTFw\r\neDsgbGluZS1oZWlnaHQ6IDEzcHg7IGNvbG9yOiAjNjY2NjY2OyI+CiAgICAgICAgICAgICAgICAgICAg\r\nVGhpcyBlbWFpbCBpcyBzZW50IHRvIDxhIHN0eWxlPSJjb2xvcjogIzY2NjY2NjsgdGV4dC1kZWNvcmF0\r\naW9uOiBub25lOyBmb250LXN0eWxlOiBpdGFsaWM7IiBocmVmPSJtYWlsdG86ZW1haWwxQGRvbWFpbi5j\r\nb20iPmVtYWlsMUBkb21haW4uY29tPC9hPiBwZXIgaWRlbnRpdHkgdmVyaWZpY2F0aW9uIHJlcXVlc3Qg\r\nb24gPGEgc3R5bGU9ImNvbG9yOiAjM2E2ZWE1OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IiBocmVmPSJo\r\ndHRwOi8vc2l0ZS91cmwiPkVYRkU8L2E+LgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAg\r\nPC90cj4KICAgICAgICA8L3RhYmxlPgogICAgPC9ib2R5Pgo8L2h0bWw+Cg==\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n"
	assert.Equal(t, content, expectPrivate)
}

func TestUserConfirmTwitter(t *testing.T) {
	arg := model.UserConfirm{}
	arg.To = rtwitter1
	arg.By = email2

	err := arg.Parse(&config)
	assert.Equal(t, err, nil)
	content, err := GetContent(localTemplate, "user_confirm", arg.To, arg)
	assert.Equal(t, err, nil)
	t.Logf("content:---------start---------\n%s\n---------end----------", content)
	expectPrivate := "\\(You(twitter1@domain.com)\\) has been requested for verification by \\(email2 name\\) on EXFE. Please click here to verify: \\(http://site/url/#!token=recipient_twitter1_token\\)"
	assert.Equal(t, content, expectPrivate)
}

func TestUserResetEmail(t *testing.T) {
	arg := model.ThirdpartTo{}
	arg.To = remail1

	err := arg.Parse(&config)
	assert.Equal(t, err, nil)
	content, err := GetContent(localTemplate, "user_resetpass", arg.To, arg)
	assert.Equal(t, err, nil)
	t.Logf("content:---------start---------\n%s\n---------end----------", content)
	expectPrivate := "Content-Type: multipart/alternative; boundary=\"bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\"\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <email1@domain.com>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <x@test.com>\nSubject: =?utf-8?B?RVhGRSByZXNldCBmb3Jnb3R0ZW4gcGFzc3dvcmQ=?=\n\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\nSGksIGVtYWlsMSBuYW1lLgoKWW91IGp1c3QgcmVxdWVzdGVkIHRvIHJlc2V0IHlvdXIgZm9yZ290dGVu\r\nIEVYRkUgcGFzc3dvcmQuCgpQbGVhc2UgY2xpY2sgaGVyZSB0byBzZXQgbmV3IHBhc3N3b3JkOiBodHRw\r\nOi8vc2l0ZS91cmwvIyF0b2tlbj1yZWNpcGllbnRfZW1haWwxX3Rva2VuIChUaGlzIHNpbmdsZS11c2Ug\r\nbGluayB3aWxsIGJlIGV4cGlyZWQgaW4gMSBkYXkuKQoKVGhpcyBlbWFpbCBpcyBzZW50IHRvIGVtYWls\r\nMUBkb21haW4uY29tIHBlciBmb3Jnb3QgcGFzc3dvcmQgcmVxdWVzdCBvbiBFWEZFKCBodHRwOi8vc2l0\r\nZS91cmwgKS4K\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIGh0bWw+CjxodG1sPgogICAgPGhlYWQ+CiAgICAgICAgPHRpdGxlPjwvdGl0bGU+CiAg\r\nICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVudC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNo\r\nYXJzZXQ9VVRGLTgiPgogICAgPC9oZWFkPgogICAgPGJvZHk+CiAgICAgICAgPHRhYmxlIGJvcmRlcj0i\r\nMCIgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0iZm9udC1mYW1pbHk6IFZlcmRh\r\nbmE7IGZvbnQtc2l6ZTogMTRweDsgbGluZS1oZWlnaHQ6IDIwcHg7IGNvbG9yOiAjMzMzMzMzOyBmb250\r\nLXdlaWdodDogbm9ybWFsOyB3aWR0aDogNjQwcHg7IHBhZGRpbmc6IDIwcHg7IGJhY2tncm91bmQtY29s\r\nb3I6ICNmYmZiZmI7Ij4KICAgICAgICAgICAgPHRyPgogICAgICAgICAgICAgICAgPHRkPgogICAgICAg\r\nICAgICAgICAgICAgIEhpLCA8c3BhbiBzdHlsZT0iY29sb3I6ICMxOTE5MTk7Ij5lbWFpbDEgbmFtZTwv\r\nc3Bhbj4KICAgICAgICAgICAgICAgIDwvdGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0\r\ncj48dGQgaGVpZ2h0PSIyMCI+PC90ZD48L3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAg\r\nICA8dGQ+CiAgICAgICAgICAgICAgICAgICAgWW91IGp1c3QgcmVxdWVzdGVkIHRvIHJlc2V0IHlvdXIg\r\nZm9yZ290dGVuIDxhIHN0eWxlPSJjb2xvcjogIzNhNmVhNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyIg\r\naHJlZj0iaHR0cDovL3NpdGUvdXJsIj5FWEZFPC9hPiBwYXNzd29yZC4KICAgICAgICAgICAgICAgIDwv\r\ndGQ+CiAgICAgICAgICAgIDwvdHI+CiAgICAgICAgICAgIDx0cj48dGQgaGVpZ2h0PSIyMCI+PC90ZD48\r\nL3RyPgogICAgICAgICAgICA8dHI+CiAgICAgICAgICAgICAgICA8dGQ+CiAgICAgICAgICAgICAgICAg\r\nICAgUGxlYXNlIGNsaWNrIGhlcmUgdG8gc2V0IG5ldyBwYXNzd29yZDogPGEgc3R5bGU9ImNvbG9yOiMz\r\nMzMzMzM7IHRleHQtZGVjb3JhdGlvbjogdW5kZXJsaW5lOyIgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMh\r\ndG9rZW49cmVjaXBpZW50X2VtYWlsMV90b2tlbiI+aHR0cDovL3NpdGUvdXJsLyMhdG9rZW49cmVjaXDi\r\ngKY8L2E+PGJyIC8+CiAgICAgICAgICAgICAgICAgICAgVGhpcyBzaW5nbGUtdXNlIGxpbmsgd2lsbCBi\r\nZSBleHBpcmVkIGluIDEgZGF5LgogICAgICAgICAgICAgICAgPC90ZD4KICAgICAgICAgICAgPC90cj4K\r\nICAgICAgICAgICAgPHRyPjx0ZCBoZWlnaHQ9IjQwIj48L3RkPjwvdHI+CiAgICAgICAgICAgIDx0cj4K\r\nICAgICAgICAgICAgICAgIDx0ZCBzdHlsZT0iZm9udC1zaXplOiAxMXB4OyBsaW5lLWhlaWdodDogMTNw\r\neDsgY29sb3I6ICM2NjY2NjY7Ij4KICAgICAgICAgICAgICAgICAgICBUaGlzIGVtYWlsIGlzIHNlbnQg\r\ndG8gPGEgc3R5bGU9ImNvbG9yOiAjNjY2NjY2OyB0ZXh0LWRlY29yYXRpb246IG5vbmU7IGZvbnQtc3R5\r\nbGU6IGl0YWxpYzsiIGhyZWY9Im1haWx0bzplbWFpbDFAZG9tYWluLmNvbSI+ZW1haWwxQGRvbWFpbi5j\r\nb208L2E+IHBlciBmb3Jnb3QgcGFzc3dvcmQgcmVxdWVzdCBvbiA8YSBzdHlsZT0iY29sb3I6ICMzYTZl\r\nYTU7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiIGhyZWY9Imh0dHA6Ly9zaXRlL3VybCI+RVhGRTwvYT4u\r\nCiAgICAgICAgICAgICAgICA8L3RkPgogICAgICAgICAgICA8L3RyPgogICAgICAgIDwvdGFibGU+CiAg\r\nICA8L2JvZHk+CjwvaHRtbD4K\n--bf4da18faf42d84da6be65288a47b5d8ca22e5f1c8d3af179533aa734d60--\n"
	assert.Equal(t, content, expectPrivate)
}

func TestUserResetTwitter(t *testing.T) {
	arg := model.ThirdpartTo{}
	arg.To = rtwitter1

	err := arg.Parse(&config)
	assert.Equal(t, err, nil)
	content, err := GetContent(localTemplate, "user_resetpass", arg.To, arg)
	assert.Equal(t, err, nil)
	t.Logf("content:---------start---------\n%s\n---------end----------", content)
	expectPrivate := "Please click here to set new password: \\(http://site/url/#!token=recipient_twitter1_token\\) (This single-use link will be expired in 1 day.)"
	assert.Equal(t, content, expectPrivate)
}

package notifier

import (
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
)

func TestCrossDigestV3Email(t *testing.T) {
	widget := map[string]interface{}{
		"image":     "RedRiverValley.jpg",
		"widget_id": 0,
		"id":        0,
		"type":      "Background",
	}
	cross1 := cross
	cross1.Exfee = exfee1
	cross1.Widgets = append(cross1.Widgets, widget)

	cross1.Exfee.Invitations[0].RsvpStatus = model.RsvpAccepted
	cross1.Exfee.Invitations[1].RsvpStatus = model.RsvpAccepted
	cross1.Exfee.Invitations[2].RsvpStatus = model.RsvpNoresponse
	cross1.Exfee.Invitations[3].RsvpStatus = model.RsvpDeclined

	arg := map[string]interface{}{
		"To":     remail1,
		"Cross":  cross1,
		"Config": &config,
	}

	text, err := GenerateContent(localTemplate, "cross_digest_v3", remail1.Provider, remail1.Language, arg)
	assert.Equal(t, err, nil)
	t.Logf("text:-----start------\n%s\n-------end-------", text)
	expect := "Content-Type: multipart/mixed; boundary=\"mixsplitter\"\nReferences: <+123@exfe.com>\nTo: =?utf-8?B?ZW1haWwxIG5hbWU=?= <email1@domain.com>\nFrom: =?utf-8?B?YnVzaW5lc3MgdGVzdGVy?= <+123@test.com>\nSubject: =?utf-8?B?VGVzdCBDcm9zcw==?=\n\n--mixsplitter\nContent-Type: multipart/alternative; boundary=\"alternativesplitter\"\n\n--alternativesplitter\nContent-Type: text/plain; charset=utf-8\nContent-Transfer-Encoding: base64\n\ntext\n\n--alternativesplitter\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\n\nPCFET0NUWVBFIEhUTUwgUFVCTElDICItLy9XM0MvL0RURCBYSFRNTCAxLjAgVHJhbnNpdGlvbmFsIC8v\r\nRU4iICJodHRwOi8vd3d3LnczLm9yZy9UUi94aHRtbDEvRFREL3hodG1sMS10cmFuc2l0aW9uYWwuZHRk\r\nIj4KCjxodG1sPgo8aGVhZD4KCTx0aXRsZT48L3RpdGxlPgoJPG1ldGEgaHR0cC1lcXVpdj0iQ29udGVu\r\ndC1UeXBlIiBjb250ZW50PSJ0ZXh0L2h0bWw7IGNoYXJzZXQ9dXRmLTgiPgo8L2hlYWQ+Cjxib2R5IHN0\r\neWxlPSJmb250LWZhbWlseTpIZWx2ZXRpY2EgTmV1ZSxIZWx2ZXRpY2EsQXJpYWwsc2Fucy1zZXJpZjsg\r\nbWF4LXdpZHRoOjY0MHB4OyBiYWNrZ3JvdW5kLWNvbG9yOndoaXRlOyBtYXJnaW46MDsgcGFkZGluZzow\r\nOyBjb2xvcjojMzMzMzMzOyBmb250LXNpemU6MTRweDsgbGluZS1oZWlnaHQ6MjBweDsgZm9udC13ZWln\r\naHQ6MzAwOyI+Cgk8dGFibGUgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0id2lk\r\ndGg6MTAwJTsgYm9yZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+CgkJPHRi\r\nb2R5PgoJCQk8dHI+PHRkIGNvbHNwYW49IjMiPgoJCQkJPHRhYmxlIGNlbGxwYWRkaW5nPSIwIiBjZWxs\r\nc3BhY2luZz0iMCIgYmdjb2xvcj0iIzU0NTQ1NCIgc3R5bGU9ImRpc3BsYXk6aW5saW5lLWJsb2NrOyB3\r\naWR0aDoxMDAlOyBib3JkZXItc3BhY2luZzowOyB2ZXJ0aWNhbC1hbGlnbjp0b3A7IGJhY2tncm91bmQt\r\naW1hZ2U6dXJsKGNpZDp0aXRsZV9iZ19qcGcpOyI+CgkJCQkJPHRib2R5PgoJCQkJCQk8dHI+CgkJCQkJ\r\nCQk8dGQgc3R5bGU9IndpZHRoOjU0cHg7IGhlaWdodDoxMDAlOyBwYWRkaW5nLXJpZ2h0OjZweDsgdmVy\r\ndGljYWwtYWxpZ246dG9wOyB0ZXh0LWFsaWduOnJpZ2h0OyI+PC90ZD4KCQkJCQkJCTx0ZCBzdHlsZT0i\r\nd2lkdGg6MnB4OyBoZWlnaHQ6MTAwJTsgdmVydGljYWwtYWxpZ246dG9wOyI+PC90ZD4KCQkJCQkJCTx0\r\nZCBzdHlsZT0icGFkZGluZzo1cHggMTBweCA1cHggMTBweDsgbWF4LXdpZHRoOjQ5OHB4OyI+CgkJCQkJ\r\nCQkJPHRhYmxlIGNlbGxwYWRkaW5nPSIwIiBjZWxsc3BhY2luZz0iMCIgc3R5bGU9IndpZHRoOjEwMCU7\r\nIGJvcmRlci1zcGFjaW5nOjA7IGJvcmRlci1jb2xsYXBzZTpjb2xsYXBzZTsiPgoJCQkJCQkJCQk8dGJv\r\nZHk+CgkJCQkJCQkJCQk8dHI+PHRkIHN0eWxlPSJjb2xvcjp3aGl0ZTsgcGFkZGluZy1ib3R0b206MTJw\r\neDsgbWF4LWhlaWdodDoyMHB4OyBvdmVyZmxvdzpoaWRkZW47IGRpc3BsYXk6aW5saW5lLWJsb2NrOyB0\r\nZXh0LW92ZXJmbG93OmVsbGlwc2lzOyBsaW5lLWhlaWdodDoyNXB4OyB0ZXh0LWFsaWduOmxlZnQ7Ij7C\r\nt1jCtyBkaWdlc3QgdG9kYXkuPC90ZD48L3RyPgoJCQkJCQkJCQkJPHRyPjx0ZCBzdHlsZT0iaGVpZ2h0\r\nOjEycHg7IGNvbG9yOiNGRjdFOTg7IGZvbnQtc2l6ZToxMHB4OyBsaW5lLWhlaWdodDoxMnB4OyI+VGl0\r\nbGUgdXBkYXRlZC48L3RkPjwvdHI+CgkJCQkJCQkJCQk8dHI+PHRkIHN0eWxlPSJtYXgtaGVpZ2h0OjQ4\r\ncHg7IG92ZXJmbG93OmhpZGRlbjsgZGlzcGxheTppbmxpbmUtYmxvY2s7IHRleHQtb3ZlcmZsb3c6ZWxs\r\naXBzaXM7IGNvbG9yOndoaXRlOyBmb250LXNpemU6MjJweDsgbGluZS1oZWlnaHQ6MjRweDsgZm9udC13\r\nZWlnaHQ6NTAwOyB0ZXh0LWFsaWduOmxlZnQ7Ij48YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tl\r\nbj0iIHN0eWxlPSJjb2xvcjp3aGl0ZTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+VGVzdCBDcm9zczwv\r\nYT48L3RkPjwvdHI+CgkJCQkJCQkJCTwvdGJvZHk+CgkJCQkJCQkJPC90YWJsZT4KCQkJCQkJCTwvdGQ+\r\nCgkJCQkJCTwvdHI+CgkJCQkJPC90Ym9keT4KCQkJCTwvdGFibGU+CgkJCTwvdGQ+PC90cj4KCQkJPHRy\r\nPgoJCQkJPHRkIHN0eWxlPSJ3aWR0aDo1NHB4OyBoZWlnaHQ6MTAwJTsgcGFkZGluZy1yaWdodDo2cHg7\r\nIHZlcnRpY2FsLWFsaWduOnRvcDsgdGV4dC1hbGlnbjpyaWdodDsiPjwvdGQ+CgkJCQk8dGQgc3R5bGU9\r\nIndpZHRoOjJweDsgaGVpZ2h0OjEwMCU7IGJhY2tncm91bmQtY29sb3I6I0U2RTZFNjsgdmVydGljYWwt\r\nYWxpZ246dG9wOyI+PC90ZD4KCQkJCTx0ZCBzdHlsZT0icGFkZGluZzo1cHggMTBweCAyMHB4IDEwcHg7\r\nIj4KCQkJCQk8dGFibGUgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0id2lkdGg6\r\nMTAwJTsgYm9yZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+CgkJCQkJCTx0\r\nYm9keT4KCQkJCQkJCTx0cj48dGQgc3R5bGU9ImhlaWdodDoxMnB4OyBjb2xvcjojRkY3RTk4OyBmb250\r\nLXNpemU6MTBweDsgbGluZS1oZWlnaHQ6MTJweDsiPlVwZGF0ZWQ8L3RkPjwvdHI+CgkJCQkJCQk8dHI+\r\nPHRkPjxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPSIgc3R5bGU9ImNvbG9yOiMzMzMzMzM7\r\nIHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPnRlc3QgY3Jvc3MgZGVzY3JpcHRpb248L2E+PC90ZD48L3Ry\r\nPgoJCQkJCQk8L3Rib2R5PgoJCQkJCTwvdGFibGU+CgkJCQk8L3RkPgoJCQk8L3RyPgoJCQk8dHI+CgkJ\r\nCQk8dGQgc3R5bGU9IndpZHRoOjU0cHg7IGhlaWdodDoxMDAlOyBwYWRkaW5nOjE1cHggNnB4IDAgMDsg\r\ndmVydGljYWwtYWxpZ246dG9wOyB0ZXh0LWFsaWduOnJpZ2h0OyI+PGltZyBzdHlsZT0ib3V0bGluZTpu\r\nb25lOyB0ZXh0LWRlY29yYXRpb246bm9uZTsgdmVydGljYWwtYWxpZ246dG9wOyIgc3JjPSJjaWQ6dGlt\r\nZV8zMGJsdWVfMnhfcG5nIiB3aWR0aD0iMzBweCIgaGVpZ2h0PSIzMHB4IiAvPjwvdGQ+CgkJCQk8dGQg\r\nc3R5bGU9IndpZHRoOjJweDsgaGVpZ2h0OjEwMCU7IGJhY2tncm91bmQtY29sb3I6I0U2RTZFNjsgdmVy\r\ndGljYWwtYWxpZ246dG9wOyI+PC90ZD4KCQkJCTx0ZCBzdHlsZT0icGFkZGluZzo1cHggMTBweCAyMHB4\r\nIDEwcHg7Ij4KCQkJCQk8dGFibGUgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0i\r\nd2lkdGg6MTAwJTsgYm9yZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+CgkJ\r\nCQkJCTx0Ym9keT4KCQkJCQkJCTx0cj48dGQgc3R5bGU9ImhlaWdodDoxMnB4OyBjb2xvcjojRkY3RTk4\r\nOyBmb250LXNpemU6MTBweDsgbGluZS1oZWlnaHQ6MTJweDsiPlVwZGF0ZWQ8L3RkPjwvdHI+CgkJCQkJ\r\nCQkKCQkJCQkJCTx0cj48dGQgc3R5bGU9ImhlaWdodDoyNHB4OyBjb2xvcjojM0E2RUE1OyBmb250LXNp\r\nemU6MThweDsgbGluZS1oZWlnaHQ6MjRweDsgZm9udC13ZWlnaHQ6NTAwOyI+PGEgaHJlZj0iaHR0cDov\r\nL3NpdGUvdXJsLyMhdG9rZW49IiBzdHlsZT0iY29sb3I6IzNBNkVBNTsgdGV4dC1kZWNvcmF0aW9uOiBu\r\nb25lOyI+VGltZTwvYT48L3RkPjwvdHI+CgkJCQkJCQk8dHI+PHRkPjxhIGhyZWY9Imh0dHA6Ly9zaXRl\r\nL3VybC8jIXRva2VuPSIgc3R5bGU9ImNvbG9yOiMzMzMzMzM7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsi\r\nPlRvIGJlIGRlY2lkZWQ8L2E+PC90ZD48L3RyPgoJCQkJCQkJCgkJCQkJCTwvdGJvZHk+CgkJCQkJPC90\r\nYWJsZT4KCQkJCTwvdGQ+CgkJCTwvdHI+CgkJCTx0cj4KCQkJCTx0ZCBzdHlsZT0id2lkdGg6NTRweDsg\r\naGVpZ2h0OjEwMCU7IHBhZGRpbmc6MTVweCA2cHggMCAwOyB2ZXJ0aWNhbC1hbGlnbjp0b3A7IHRleHQt\r\nYWxpZ246cmlnaHQ7Ij48aW1nIHN0eWxlPSJvdXRsaW5lOm5vbmU7IHRleHQtZGVjb3JhdGlvbjpub25l\r\nOyB2ZXJ0aWNhbC1hbGlnbjp0b3A7IiBzcmM9ImNpZDpwbGFjZV8zMGJsdWVfMnhfcG5nIiB3aWR0aD0i\r\nMzBweCIgaGVpZ2h0PSIzMHB4IiAvPjwvdGQ+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjJweDsgaGVpZ2h0\r\nOjEwMCU7IGJhY2tncm91bmQtY29sb3I6I0U2RTZFNjsgdmVydGljYWwtYWxpZ246dG9wOyI+PC90ZD4K\r\nCQkJCTx0ZCBzdHlsZT0icGFkZGluZzo1cHggMTBweCAyMHB4IDEwcHg7Ij4KCQkJCQk8dGFibGUgY2Vs\r\nbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBzdHlsZT0id2lkdGg6MTAwJTsgYm9yZGVyLXNwYWNp\r\nbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+CgkJCQkJCTx0Ym9keT4KCQkJCQkJCTx0cj48\r\ndGQgc3R5bGU9ImhlaWdodDoxMnB4OyBjb2xvcjojRkY3RTk4OyBmb250LXNpemU6MTBweDsgbGluZS1o\r\nZWlnaHQ6MTJweDsiPlVwZGF0ZWQ8L3RkPjwvdHI+CgkJCQkJCQkKCQkJCQkJCTx0cj48dGQgc3R5bGU9\r\nImhlaWdodDoyNHB4OyBjb2xvcjojM0E2RUE1OyBmb250LXNpemU6MThweDsgbGluZS1oZWlnaHQ6MjRw\r\neDsgZm9udC13ZWlnaHQ6NTAwOyI+PGEgaHJlZj0iaHR0cDovL3NpdGUvdXJsLyMhdG9rZW49IiBzdHls\r\nZT0iY29sb3I6IzNBNkVBNTsgdGV4dC1kZWNvcmF0aW9uOiBub25lOyI+UGxhY2U8L2E+PC90ZD48L3Ry\r\nPgoJCQkJCQkJPHRyPjx0ZD48YSBocmVmPSJodHRwOi8vc2l0ZS91cmwvIyF0b2tlbj0iIHN0eWxlPSJj\r\nb2xvcjojMzMzMzMzOyB0ZXh0LWRlY29yYXRpb246IG5vbmU7Ij5UbyBiZSBkZWNpZGVkPC9hPjwvdGQ+\r\nPC90cj4KCQkJCQkJCQoJCQkJCQk8L3Rib2R5PgoJCQkJCTwvdGFibGU+CgkJCQk8L3RkPgoJCQk8L3Ry\r\nPgoJCQk8dHI+CgkJCQkKCQkJCQoJCQkJPHRkIHN0eWxlPSJ3aWR0aDo1NHB4OyBoZWlnaHQ6MTAwJTsg\r\ncGFkZGluZzoyNXB4IDZweCAwIDA7IHZlcnRpY2FsLWFsaWduOnRvcDsgdGV4dC1hbGlnbjpyaWdodDsi\r\nPgoJCQkJCQoJCQkJCTxpbWcgc3R5bGU9Im91dGxpbmU6bm9uZTsgdGV4dC1kZWNvcmF0aW9uOm5vbmU7\r\nIHZlcnRpY2FsLWFsaWduOnRvcDsiIHNyYz0iY2lkOnJzdnBfYWNjZXB0ZWRfMjZibHVlXzJ4X3BuZyIg\r\nd2lkdGg9IjMwcHgiIGhlaWdodD0iMzBweCIgLz4KCQkJCQkKCQkJCQkKCQkJCQkKCQkJCTwvdGQ+CgkJ\r\nCQk8dGQgc3R5bGU9IndpZHRoOjJweDsgaGVpZ2h0OjEwMCU7IGJhY2tncm91bmQtY29sb3I6I0U2RTZF\r\nNjsgdmVydGljYWwtYWxpZ246dG9wOyI+PC90ZD4KCQkJCTx0ZCBzdHlsZT0icGFkZGluZzo1cHggMTBw\r\neCAyMHB4IDEwcHg7Ij4KCQkJCQk8dGFibGUgY2VsbHBhZGRpbmc9IjAiIGNlbGxzcGFjaW5nPSIwIiBz\r\ndHlsZT0id2lkdGg6MTAwJTsgYm9yZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNl\r\nOyI+CgkJCQkJCTx0Ym9keT4KCQkJCQkJCQoJCQkJCQkJPHRyPjx0ZD5Zb3VyIHBhcnRpY2lwYXRpb24g\r\nc3RhdHVzIGlzOjwvdGQ+PC90cj4KCQkJCQkJCTx0cj48dGQgaGVpZ2h0PSIzMHB4IiBzdHlsZT0ibGlu\r\nZS1oZWlnaHQ6MzBweDsgY29sb3I6IzdGN0Y3RjsiPjxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRv\r\na2VuPSIgc3R5bGU9ImNvbG9yOiMzMzMzMzM7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPjxzcGFuIHN0\r\neWxlPSJjb2xvcjojM0E2RUE1OyBmb250LXdlaWdodDo1MDA7IGJhY2tncm91bmQtY29sb3I6I0U2RTZF\r\nNjsgYm9yZGVyLXJhZGl1czo0cHg7IHBhZGRpbmc6NXB4IDIwcHg7IHZlcnRpY2FsLWFsaWduOnRvcDsg\r\nbWFyZ2luLXJpZ2h0OjhweDsiPkFjY2VwdGVkPC9zcGFuPiBzZXQgYnkgZW1haWwxIG5hbWU8L2E+PC90\r\nZD48L3RyPgoJCQkJCQkJCgkJCQkJCQkKCQkJCQkJCQoJCQkJCQk8L3Rib2R5PgoJCQkJCTwvdGFibGU+\r\nCgkJCQk8L3RkPgoJCQk8L3RyPgoJCQk8dHI+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjU0cHg7IGhlaWdo\r\ndDoxMDAlOyBwYWRkaW5nOjVweCA2cHggMCAwOyB2ZXJ0aWNhbC1hbGlnbjp0b3A7IHRleHQtYWxpZ246\r\ncmlnaHQ7Ij48aW1nIHN0eWxlPSJvdXRsaW5lOm5vbmU7IHRleHQtZGVjb3JhdGlvbjpub25lOyB2ZXJ0\r\naWNhbC1hbGlnbjp0b3A7IiBzcmM9ImNpZDpleGZlZV8zMGJsdWVfMnhfcG5nIiB3aWR0aD0iMzBweCIg\r\naGVpZ2h0PSIzMHB4IiAvPjwvdGQ+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjJweDsgaGVpZ2h0OjEwMCU7\r\nIGJhY2tncm91bmQtY29sb3I6I0U2RTZFNjsgdmVydGljYWwtYWxpZ246dG9wOyI+PC90ZD4KCQkJCTx0\r\nZCBzdHlsZT0icGFkZGluZzo1cHggMTBweCA1cHggMTBweDsiPgoJCQkJCTx0YWJsZSBjZWxscGFkZGlu\r\nZz0iMCIgY2VsbHNwYWNpbmc9IjAiIHN0eWxlPSJ3aWR0aDoxMDAlOyBib3JkZXItc3BhY2luZzowOyBi\r\nb3JkZXItY29sbGFwc2U6Y29sbGFwc2U7Ij4KCQkJCQkJPHRib2R5PgoJCQkJCQkJPHRyPjx0ZCBoZWln\r\naHQ9IjMwcHgiPjxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC8jIXRva2VuPSIgc3R5bGU9ImNvbG9yOiMz\r\nMzMzMzM7IHRleHQtZGVjb3JhdGlvbjogbm9uZTsiPjxzcGFuIHN0eWxlPSJoZWlnaHQ6MjRweDsgY29s\r\nb3I6IzNBNkVBNTsgZm9udC1zaXplOjE4cHg7IGxpbmUtaGVpZ2h0OjI0cHg7IGZvbnQtd2VpZ2h0OjUw\r\nMDsiPkV4ZmVlcyZuYnNwOyZuYnNwOzI8L3NwYW4+LzQgYWNjZXB0ZWQ8L2E+PC90ZD48L3RyPgoJCQkJ\r\nCQkJPCEtLQoJCQkJCQkJPHRyPjx0ZCBzdHlsZT0iaGVpZ2h0OjEycHg7IGNvbG9yOiNGRjdFOTg7IGZv\r\nbnQtc2l6ZToxMHB4OyBsaW5lLWhlaWdodDoxMnB4OyI+VXBkYXRlZDwvdGQ+PC90cj4KCQkJCQkJCTx0\r\ncj48dGQ+QWNjZXB0ZWQ6IDxzcGFuIHN0eWxlPSJmb250LXdlaWdodDo1MDA7IiBzdHlsZT0iY29sb3I6\r\nIzNBNkVBNTsiPkRhdmU8L3NwYW4+LCA8c3BhbiBzdHlsZT0iZm9udC13ZWlnaHQ6NTAwOyIgc3R5bGU9\r\nImNvbG9yOiMzQTZFQTU7Ij5Ud28zMzwvc3Bhbj4gYW5kIDMgb3RoZXJzPC90ZD48L3RyPgoJCQkJCQkJ\r\nPHRyPjx0ZD5OZXdseSBpbnZpdGVkOiA8c3BhbiBzdHlsZT0iZm9udC13ZWlnaHQ6NTAwOyI+Sm9l4oCo\r\nPC9zcGFuPjwvdGQ+PC90cj4KCQkJCQkJCTx0cj48dGQ+VW5hdmFpbGFibGU6IDxzcGFuIHN0eWxlPSJm\r\nb250LXdlaWdodDo1MDA7Ij5kbS48L3NwYW4+PC90ZD48L3RyPgoJCQkJCQkJPHRyPjx0ZD5SZW1vdmVk\r\nOiA8c3BhbiBzdHlsZT0iZm9udC13ZWlnaHQ6NTAwOyI+PGRlbD5Kb2VsPC9kZWw+PC9zcGFuPjwvdGQ+\r\nPC90cj4KCQkJCQkJCS0tPgoJCQkJCQk8L3Rib2R5PgoJCQkJCTwvdGFibGU+CgkJCQk8L3RkPgoJCQk8\r\nL3RyPgoJCQkKCQkJCgkJCTx0ciBzdHlsZT0idmVydGljYWwtYWxpZ246dG9wOyI+CgkJCQk8dGQgc3R5\r\nbGU9IndpZHRoOjU0cHg7IGhlaWdodDoxMDAlOyBwYWRkaW5nLXJpZ2h0OjZweDsgdmVydGljYWwtYWxp\r\nZ246dG9wOyB0ZXh0LWFsaWduOnJpZ2h0OyI+PGltZyBzdHlsZT0ib3V0bGluZTpub25lOyB0ZXh0LWRl\r\nY29yYXRpb246bm9uZTsgdmVydGljYWwtYWxpZ246dG9wOyIgc3JjPSJjaWQ6aTExIiB3aWR0aD0iMjRw\r\neCIgaGVpZ2h0PSIxOXB4IiBzdHlsZT0idmVydGljYWwtYWxpZ246dG9wOyIgLz48L3RkPgoJCQkJPHRk\r\nIHN0eWxlPSJ3aWR0aDoycHg7IGhlaWdodDoxMDAlOyBiYWNrZ3JvdW5kLWNvbG9yOiMzNzg0RDU7IHZl\r\ncnRpY2FsLWFsaWduOnRvcDsiPjwvdGQ+CgkJCQk8dGQgc3R5bGU9InBhZGRpbmc6MCAxMHB4OyI+CgkJ\r\nCQkJPHRhYmxlIGNlbGxwYWRkaW5nPSIwIiBjZWxsc3BhY2luZz0iMCIgc3R5bGU9InRhYmxlLWxheW91\r\ndDpmaXhlZDsgd2lkdGg6MTAwJTsgYm9yZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxh\r\ncHNlOyI+PHRib2R5PgoJCQkJCQk8dHI+PHRkIHN0eWxlPSJvdmVyZmxvdzpoaWRkZW47IHdoaXRlLXNw\r\nYWNlOm5vd3JhcDsgdGV4dC1vdmVyZmxvdzplbGxpcHNpczsiPgoJCQkJCQkJZW1haWwxIG5hbWU8c3Bh\r\nbiBzdHlsZT0iY29sb3I6IzdGN0Y3RjsgbWFyZ2luLWxlZnQ6MTVweDsgZm9udC1zdHlsZTppdGFsaWM7\r\nIj5lbWFpbDFAZG9tYWluLmNvbUBlbWFpbDwvc3Bhbj4KCQkJCQkJPC90ZD48L3RyPgoJCQkJCTwvdGJv\r\nZHk+PC90YWJsZT4KCQkJCTwvdGQ+CgkJCTwvdHI+CgkJCQoJCQkKCQkJCgkJCTx0ciBzdHlsZT0idmVy\r\ndGljYWwtYWxpZ246dG9wOyI+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjU0cHg7IGhlaWdodDoxMDAlOyBw\r\nYWRkaW5nLXJpZ2h0OjZweDsgdmVydGljYWwtYWxpZ246dG9wOyB0ZXh0LWFsaWduOnJpZ2h0OyI+PGlt\r\nZyBzdHlsZT0ib3V0bGluZTpub25lOyB0ZXh0LWRlY29yYXRpb246bm9uZTsgdmVydGljYWwtYWxpZ246\r\ndG9wOyIgc3JjPSJjaWQ6aTEyIiB3aWR0aD0iMjRweCIgaGVpZ2h0PSIxOXB4IiBzdHlsZT0idmVydGlj\r\nYWwtYWxpZ246dG9wOyIgLz48L3RkPgoJCQkJPHRkIHN0eWxlPSJ3aWR0aDoycHg7IGhlaWdodDoxMDAl\r\nOyBiYWNrZ3JvdW5kLWNvbG9yOiMzNzg0RDU7IHZlcnRpY2FsLWFsaWduOnRvcDsiPjwvdGQ+CgkJCQk8\r\ndGQgc3R5bGU9InBhZGRpbmc6MCAxMHB4OyI+CgkJCQkJPHRhYmxlIGNlbGxwYWRkaW5nPSIwIiBjZWxs\r\nc3BhY2luZz0iMCIgc3R5bGU9InRhYmxlLWxheW91dDpmaXhlZDsgd2lkdGg6MTAwJTsgYm9yZGVyLXNw\r\nYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+PHRib2R5PgoJCQkJCQk8dHI+PHRkIHN0\r\neWxlPSJvdmVyZmxvdzpoaWRkZW47IHdoaXRlLXNwYWNlOm5vd3JhcDsgdGV4dC1vdmVyZmxvdzplbGxp\r\ncHNpczsiPgoJCQkJCQkJZW1haWwyIG5hbWU8c3BhbiBzdHlsZT0iY29sb3I6IzdGN0Y3RjsgbWFyZ2lu\r\nLWxlZnQ6MTVweDsgZm9udC1zdHlsZTppdGFsaWM7Ij5lbWFpbDJAZG9tYWluLmNvbUBlbWFpbDwvc3Bh\r\nbj4KCQkJCQkJPC90ZD48L3RyPgoJCQkJCTwvdGJvZHk+PC90YWJsZT4KCQkJCTwvdGQ+CgkJCTwvdHI+\r\nCgkJCQoJCQkKCQkJCgkJCQoJCQkKCQkJCgkJCQoJCQkKCQkJCgkJCQoJCQkKCQkJCgkJCTx0ciBzdHls\r\nZT0idmVydGljYWwtYWxpZ246dG9wOyI+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjU0cHg7IGhlaWdodDox\r\nMDAlOyBwYWRkaW5nLXJpZ2h0OjZweDsgdmVydGljYWwtYWxpZ246dG9wOyB0ZXh0LWFsaWduOnJpZ2h0\r\nOyI+PGltZyBzdHlsZT0ib3V0bGluZTpub25lOyB0ZXh0LWRlY29yYXRpb246bm9uZTsgdmVydGljYWwt\r\nYWxpZ246dG9wOyIgc3JjPSJjaWQ6aTIyIiB3aWR0aD0iMjRweCIgaGVpZ2h0PSIxOXB4IiBzdHlsZT0i\r\ndmVydGljYWwtYWxpZ246dG9wOyIgLz48L3RkPgoJCQkJPHRkIHN0eWxlPSJ3aWR0aDoycHg7IGhlaWdo\r\ndDoxMDAlOyBiYWNrZ3JvdW5kLWNvbG9yOiNFNkU2RTY7IHZlcnRpY2FsLWFsaWduOnRvcDsiPjwvdGQ+\r\nCgkJCQk8dGQgc3R5bGU9InBhZGRpbmc6MCAxMHB4OyI+CgkJCQkJPHRhYmxlIGNlbGxwYWRkaW5nPSIw\r\nIiBjZWxsc3BhY2luZz0iMCIgc3R5bGU9InRhYmxlLWxheW91dDpmaXhlZDsgd2lkdGg6MTAwJTsgYm9y\r\nZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+PHRib2R5PgoJCQkJCQk8dHI+\r\nPHRkIHN0eWxlPSJvdmVyZmxvdzpoaWRkZW47IHdoaXRlLXNwYWNlOm5vd3JhcDsgdGV4dC1vdmVyZmxv\r\ndzplbGxpcHNpczsiPgoJCQkJCQkJdHdpdHRlcjMgbmFtZTxzcGFuIHN0eWxlPSJjb2xvcjojN0Y3RjdG\r\nOyBtYXJnaW4tbGVmdDoxNXB4OyBmb250LXN0eWxlOml0YWxpYzsiPnR3aXR0ZXIzQGRvbWFpbi5jb21A\r\ndHdpdHRlcjwvc3Bhbj4KCQkJCQkJPC90ZD48L3RyPgoJCQkJCTwvdGJvZHk+PC90YWJsZT4KCQkJCTwv\r\ndGQ+CgkJCTwvdHI+CgkJCQoJCQkKCQkJCgkJCQoJCQkKCQkJCgkJCQoJCQkKCQkJCgkJCQoJCQkKCQkJ\r\nCgkJCTx0ciBzdHlsZT0idmVydGljYWwtYWxpZ246dG9wOyI+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjU0\r\ncHg7IGhlaWdodDoxMDAlOyBwYWRkaW5nLXJpZ2h0OjZweDsgdmVydGljYWwtYWxpZ246dG9wOyB0ZXh0\r\nLWFsaWduOnJpZ2h0OyI+PGltZyBzdHlsZT0ib3V0bGluZTpub25lOyB0ZXh0LWRlY29yYXRpb246bm9u\r\nZTsgdmVydGljYWwtYWxpZ246dG9wOyIgc3JjPSJjaWQ6aTMyIiB3aWR0aD0iMjRweCIgaGVpZ2h0PSIx\r\nOXB4IiBzdHlsZT0idmVydGljYWwtYWxpZ246dG9wOyIgLz48L3RkPgoJCQkJPHRkIHN0eWxlPSJ3aWR0\r\naDoycHg7IGhlaWdodDoxMDAlOyBiYWNrZ3JvdW5kLWNvbG9yOiM3RjdGN0Y7IHZlcnRpY2FsLWFsaWdu\r\nOnRvcDsiPjwvdGQ+CgkJCQk8dGQgc3R5bGU9InBhZGRpbmc6MCAxMHB4OyI+CgkJCQkJPHRhYmxlIGNl\r\nbGxwYWRkaW5nPSIwIiBjZWxsc3BhY2luZz0iMCIgc3R5bGU9InRhYmxlLWxheW91dDpmaXhlZDsgd2lk\r\ndGg6MTAwJTsgYm9yZGVyLXNwYWNpbmc6MDsgYm9yZGVyLWNvbGxhcHNlOmNvbGxhcHNlOyI+PHRib2R5\r\nPgoJCQkJCQk8dHI+PHRkIHN0eWxlPSJvdmVyZmxvdzpoaWRkZW47IHdoaXRlLXNwYWNlOm5vd3JhcDsg\r\ndGV4dC1vdmVyZmxvdzplbGxpcHNpczsiPgoJCQkJCQkJZmFjZWJvb2s0IG5hbWU8c3BhbiBzdHlsZT0i\r\nY29sb3I6IzdGN0Y3RjsgbWFyZ2luLWxlZnQ6MTVweDsgZm9udC1zdHlsZTppdGFsaWM7Ij5mYWNlYm9v\r\nazRAZG9tYWluLmNvbUBmYWNlYm9vazwvc3Bhbj4KCQkJCQkJPC90ZD48L3RyPgoJCQkJCTwvdGJvZHk+\r\nPC90YWJsZT4KCQkJCTwvdGQ+CgkJCTwvdHI+CgkJCQoJCQkKCQkJPHRyIHN0eWxlPSJ2ZXJ0aWNhbC1h\r\nbGlnbjp0b3AiPgoJCQkJPHRkPjwvdGQ+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjJweDsgaGVpZ2h0OjEw\r\nMCU7IGJhY2tncm91bmQtY29sb3I6I0U2RTZFNjsgdmVydGljYWwtYWxpZ246dG9wOyI+PC90ZD4KCQkJ\r\nCTx0ZCBoZWlnaHQ9IjMwcHgiPiZuYnNwOzwvdGQ+CgkJCTwvdHI+CgkJCTx0ciBjb2xvcj0iIzdGN0Y3\r\nRiIgYmdjb2xvcj0iI0VFRUVFRSI+CgkJCQk8dGQgc3R5bGU9IndpZHRoOjU0cHg7IGhlaWdodDoxMDAl\r\nOyBwYWRkaW5nLXJpZ2h0OjZweDsgdmVydGljYWwtYWxpZ246dG9wOyB0ZXh0LWFsaWduOnJpZ2h0OyI+\r\nPC90ZD4KCQkJCTx0ZD48L3RkPgoJCQkJPHRkIHN0eWxlPSJjb2xvcjojN0Y3RjdGOyBmb250LXNpemU6\r\nMTFweDsgbGluZS1oZWlnaHQ6MTNweDsgcGFkZGluZzo4cHggMCA4cHggMTBweDsiPlJlcGx5IHRoaXMg\r\nZW1haWwgYXMgZ3JvdXAgY29udmVyc2F0aW9uLCDigJhjY+KAmSBwZW9wbGUgdG8gaW52aXRlLiBHZXQg\r\nPGEgaHJlZj0iaHR0cDovL2FwcC91cmwiIHN0eWxlPSJjb2xvcjojM0E2RUE1OyB0ZXh0LWRlY29yYXRp\r\nb246bm9uZTsiPkVYRkU8L2E+IGFwcCA8c3BhbiBzdHlsZT0iZm9udC1zdHlsZTogaXRhbGljIj5mcmVl\r\nPC9zcGFuPiB0byBlbmdhZ2UgZWFzaWVyLjxiciAvPlRoaXMgZW1haWwgaXMgZ2VuZXJhdGVkIGJ5IEVY\r\nRkUuIDxhIGhyZWY9Imh0dHA6Ly9zaXRlL3VybC9tdXRlL2Nyb3NzP3Rva2VuPSIgc3R5bGU9ImNvbG9y\r\nOiM3RjdGN0Y7Ij5VbnN1YnNjcmliZTwvYT4gZnVydGhlciB1cGRhdGVzPCEtLSBvciA8YSBocmVmPSIv\r\ncHJlZmVyZW5jZSIgc3R5bGU9ImNvbG9yOiM3RjdGN0Y7Ij5jaGFuZ2Ugbm90aWZpY2F0aW9uIHByZWZl\r\ncmVuY2U8L2E+LS0+LjwvdGQ+CgkJCTwvdHI+CgkJPC90Ym9keT4KCTwvdGFibGU+CjwvYm9keT4KPC9o\r\ndG1sPg==\n\n--alternativesplitter--\n\n--mixsplitter\nContent-Disposition: attachment; filename=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Type: text/calendar; charset=utf-8; name=\"=?UTF-8?B?VGVzdCBDcm9zcy5pY3M=?=\"\nContent-Transfer-Encoding: base64\n\nics\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"title_bg.jpg\"\nContent-Transfer-Encoding: base64\nContent-ID: <title_bg_jpg>\n\n\n\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"i11.jpg\"\nContent-Transfer-Encoding: base64\nContent-ID: <i11>\n\nPD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IURPQ1RZUEUgaHRtbCBQ\r\nVUJMSUMgIi0vL1czQy8vRFREIFhIVE1MIDEuMCBUcmFuc2l0aW9uYWwvL0VOIgogICAgICAgICAiaHR0\r\ncDovL3d3dy53My5vcmcvVFIveGh0bWwxL0RURC94aHRtbDEtdHJhbnNpdGlvbmFsLmR0ZCI+CjxodG1s\r\nIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hodG1sIiB4bWw6bGFuZz0iZW4iIGxhbmc9ImVu\r\nIj4KIDxoZWFkPgogIDx0aXRsZT40MDQgLSBOb3QgRm91bmQ8L3RpdGxlPgogPC9oZWFkPgogPGJvZHk+\r\nCiAgPGgxPjQwNCAtIE5vdCBGb3VuZDwvaDE+CiA8L2JvZHk+CjwvaHRtbD4K\n\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"i12.jpg\"\nContent-Transfer-Encoding: base64\nContent-ID: <i12>\n\nPD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IURPQ1RZUEUgaHRtbCBQ\r\nVUJMSUMgIi0vL1czQy8vRFREIFhIVE1MIDEuMCBUcmFuc2l0aW9uYWwvL0VOIgogICAgICAgICAiaHR0\r\ncDovL3d3dy53My5vcmcvVFIveGh0bWwxL0RURC94aHRtbDEtdHJhbnNpdGlvbmFsLmR0ZCI+CjxodG1s\r\nIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hodG1sIiB4bWw6bGFuZz0iZW4iIGxhbmc9ImVu\r\nIj4KIDxoZWFkPgogIDx0aXRsZT40MDQgLSBOb3QgRm91bmQ8L3RpdGxlPgogPC9oZWFkPgogPGJvZHk+\r\nCiAgPGgxPjQwNCAtIE5vdCBGb3VuZDwvaDE+CiA8L2JvZHk+CjwvaHRtbD4K\n\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"i22.jpg\"\nContent-Transfer-Encoding: base64\nContent-ID: <i22>\n\nPD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IURPQ1RZUEUgaHRtbCBQ\r\nVUJMSUMgIi0vL1czQy8vRFREIFhIVE1MIDEuMCBUcmFuc2l0aW9uYWwvL0VOIgogICAgICAgICAiaHR0\r\ncDovL3d3dy53My5vcmcvVFIveGh0bWwxL0RURC94aHRtbDEtdHJhbnNpdGlvbmFsLmR0ZCI+CjxodG1s\r\nIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hodG1sIiB4bWw6bGFuZz0iZW4iIGxhbmc9ImVu\r\nIj4KIDxoZWFkPgogIDx0aXRsZT40MDQgLSBOb3QgRm91bmQ8L3RpdGxlPgogPC9oZWFkPgogPGJvZHk+\r\nCiAgPGgxPjQwNCAtIE5vdCBGb3VuZDwvaDE+CiA8L2JvZHk+CjwvaHRtbD4K\n\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"i32.jpg\"\nContent-Transfer-Encoding: base64\nContent-ID: <i32>\n\nPD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pgo8IURPQ1RZUEUgaHRtbCBQ\r\nVUJMSUMgIi0vL1czQy8vRFREIFhIVE1MIDEuMCBUcmFuc2l0aW9uYWwvL0VOIgogICAgICAgICAiaHR0\r\ncDovL3d3dy53My5vcmcvVFIveGh0bWwxL0RURC94aHRtbDEtdHJhbnNpdGlvbmFsLmR0ZCI+CjxodG1s\r\nIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hodG1sIiB4bWw6bGFuZz0iZW4iIGxhbmc9ImVu\r\nIj4KIDxoZWFkPgogIDx0aXRsZT40MDQgLSBOb3QgRm91bmQ8L3RpdGxlPgogPC9oZWFkPgogPGJvZHk+\r\nCiAgPGgxPjQwNCAtIE5vdCBGb3VuZDwvaDE+CiA8L2JvZHk+CjwvaHRtbD4K\n\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"time_30blue2x.png\"\nContent-Transfer-Encoding: base64\nContent-ID: <time_30blue_2x_png>\n\n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\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"place_30blue2x.png\"\nContent-Transfer-Encoding: base64\nContent-ID: <place_30blue_2x_png>\n\n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\r\nLNgs+CUX/A8OrCTYs58dKAAAAABJRU5ErkJggg==\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"exfee_30blue2x.png\"\nContent-Transfer-Encoding: base64\nContent-ID: <exfee_30blue_2x_png>\n\n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\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"rsvp_accepted_26blue2x.png\"\nContent-Transfer-Encoding: base64\nContent-ID: <rsvp_accepted_26blue_2x_png>\n\niVBORw0KGgoAAAANSUhEUgAAADQAAAA0CAMAAADypuvZAAABCFBMVEX///86b6U3baQ6b6U4b6I6bqU5\r\nbKQ6bqVAgL85aqM6bqU6bqU6bqUudKI6aKI6bqU8baM6bqU6b6U6bqU6bqU6bqUrgKpJbbY6bqU7bqM6\r\nbqVVVao5caE7b6M6bqU6bqU6bqU8aaU6b6Y8b6Y5bKY6bqU6b6Q5baM7bqU6bqU5b6U4cKI7b6U6b6Q7\r\nb6MAAP86bag7a6Y6bqU6bKY6bqU6b6UzZpk5b6Q6bqU7baM4bKQ6baY4a6M2bKI6b6U6cKU6bqU4bKMz\r\nZqo5bqU6bqU6bqU6bqU6b6Y6bqU6bqU6bqU3b6Y6bqY5b6Q6bqVAYJ86bqU6bqU7baU9bao6bqQ4baZA\r\ncJ86bqU/zGEqAAAAV3RSTlMA2xzwN/kt/AQkovPyCxZPL/Wf3+X9BgfDZPcDG0W81/oRmjxQWDVn3plV\r\nKcRcTgEjK+tCf9QFWv51O7EyIekw6kAP1c7YpIzjp+AXxXjvCM+5WxV2TRA7R66CAAABBElEQVR4Xu3S\r\nx1ICURCF4avjREBAQRwy5pxzzjmn8/5vYjlVFtzqxbFZueDff73oOub/1mtyQG/GXUetVtYArWo1AK3a\r\nqgNqNQao1S6gVhNFqNVBHp05U4ZXLcCqr5+bzIXeVPb1xlx3Yaa7MLUmMWEszPwTMc8v/puxe58hxswC\r\n/pxlSgvMLAKAu9Rhlq+YWU0hUetttMFMZhNJyJ/+mu0dYkwkRrxXZiZGOyf3Yw6PmDm2rp54xpydMxNc\r\nwqoQBq98O9mUrW5upZHl7iDiG71/4EZWisCN7LHIjSyb5kbmDUvDGxyVhjcyBGl4H740vM+0xtjvkIa/\r\nQxr+jsRo+/ob6PUNLt7A+K+yBCEAAAAASUVORK5CYII=\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"rsvp_pending_26g5_2x.png\"\nContent-Transfer-Encoding: base64\nContent-ID: <rsvp_pending_26g5_2x_png>\n\niVBORw0KGgoAAAANSUhEUgAAADQAAAA0CAMAAADypuvZAAAAMFBMVEX///+AgICBgYGAgICAgICAgICH\r\nh4eAgICAgICBgYGAgICAgICAgICCgoKAgICAgIDM00BJAAAAD3RSTlMARHfuZpkRuyJViMzdM6pSPOBW\r\nAAAAz0lEQVR42u2U0Q7DIAhFFVBR2/n/f7s+NJLFRqHZntbz1sRjKFx0D9+HCdJBAK9WqLYOBtYofmuf\r\nwNqBNpDiwintgjq3crukzJy4tWu8/oeEMJGqnNqdY8D+Pamun6FzyLiuzw/V0HpYMN5rkJxBYn8icxNJ\r\nT5KCtUgf0OLYq5PwYlQqXFsnK52MQ8iX0JDWNYzikFMSpAf6CaH6fRD4TnqytEAP9Pn8PZEA2Oh4lABp\r\n2Q/HbL1ky/XIw/9rqW96vBHY4iwUWSYDPqQXuYcpb/feGW1CjldpAAAAAElFTkSuQmCC\n\n--mixsplitter\nContent-Type: image/jpeg; name=\"rsvp_unavailable_26g5_2x.png\"\nContent-Transfer-Encoding: base64\nContent-ID: <rsvp_unavailable_26g5_2x_png>\n\niVBORw0KGgoAAAANSUhEUgAAADQAAAA0CAYAAADFeBvrAAACyUlEQVR42u2aP0xTURTGa0xMQCMLboTB\r\nQaLCTJjUgUo6uHVAo4uwKaRJQ6lpXxqMOgojqJsaEwdGV53sWhL/7N0wJCwMIKF+33Bzk9b42nPOfTHk\r\nDl94Cb3fOb/wcs65p+Q6nc6pUgSKQBEoAkWgHlUqlZFGo7FWLBbPWidET3ozhghICNOEOtCrXC53xgqG\r\nXvDcojdjSKDkMF5PrYD4l3G+UigFjFeSJMtaGHo4Pw2UDsbrBLqrgJmnB720UHoYr0NoTvCazfGs89FC\r\n2cB4HdTr9el+YfhZnvHn9VB2MF6/oGtpMLVa7So/684poNKByuXyeQcjVBsa+8drNga1Ff5N5tgfkO8H\r\nL3lYoe/Q6F9gRqFvSu915pgO1Bv8mTJwExp2fnyGvio9n6uqHAxWlQl8gs5RfFZ6rVr0IUI9cn1CqPfQ\r\nO8X5E/Sqx6aTAgwfwviYATIWYy6EmOUIdQ/mRxnCHDFm0GkbQe64zh5Yh4w1aH7S2WvWdfhAOmCMLG+s\r\nHFtuIPB+AJh9emd/Bfez2J4hzB49tVdwLdQUEtk1gNmllzYfq5vmFeVs1sawOhFiSaJZmrQUQC16/BdA\r\npVJpCAl9NnjlvtBLD6RfOW0bFoVt7WpMu3J6HaBsv6F35kAc4wM21heZAqGLl0KPPoyRBRBh7rurRGAx\r\nxoOgQAhQEA6mbynhgFoIAoQuPiMcSLdYuSg8b0oGVca2BCLMdeHMtsGK1VUZNwQ+jD1pAgSjcahtWan4\r\nO+FqbFwFVK1WL8Hkh3iZYb98+cmcxEAwSASVaWmAIrMkqJiJDMi/8yt9LkeOUdIXBW1gsV9/5sKcLNZY\r\n+ZSi8NstMyTiWXqkFIW8Zdkm1GWoJVhmaJcvO4xt31j9Ev9D1zLjtuFXkvmuXvcRujCIh2jKBkSFrwF+\r\n3lKD9ELdpDf0hLGynLYv2oH0esd/vIhAESgCmegPTyst6rzGGWQAAAAASUVORK5CYII=\n\n--mixsplitter--\n"
	assert.Equal(t, text, expect)
}

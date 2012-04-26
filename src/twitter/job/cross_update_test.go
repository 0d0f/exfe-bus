package twitter_job

import (
	"testing"
)

func TestCrossUpdateMessage(t *testing.T) {
	cross := CrossUpdateArg{}
	sameTitle := cross.sameTitleMessage("Dinner 12:03am This week 2012-03-04", "To meet someone in some party fdafdasffd", "place title here", "place description maybe very very very very long")
	t.Errorf("%d\n%s", len(sameTitle), sameTitle)
	diffTitle := cross.diffTitleMessage("Dinner 12:03am This week 2012-03-04", "To check someone in some party fdafdasffd", "place title here", "place description maybe very very very very long", "To meet someone in some party afdasfadsfasdf")
	t.Errorf("%d\n%s", len(diffTitle), diffTitle)
}

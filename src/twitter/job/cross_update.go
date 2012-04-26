package twitter_job

import (
	"exfe"
	"strings"
	"gobus"
	"fmt"
)

type CrossUpdateArg struct {
	Cross exfe.Cross
	Old_cross exfe.Cross
}

func (a CrossUpdateArg) UpdateMessage(to_invitation *exfe.Invitation) (string, error) {
	time, err := a.Cross.Time.StringInZone(to_invitation.Identity.Timezone)
	if err != nil {
		return "", err
	}
	place1 := a.Cross.Place.Title
	place2 := a.Cross.Place.Description
	title := a.Cross.Title
	old_title := a.Old_cross.Title

	if title == old_title {
		return a.sameTitleMessage(time, title, place1, place2), nil
	}
	return a.diffTitleMessage(time, title, place1, place2, old_title), nil
}

const messageMaxLen = 140 - 29 /* len("Update http://t.co/fbqqsjky:\n") */ - 5 /* reserved */
const titleMaxLen = 20
const newTitleMaxLen = 13

func (a CrossUpdateArg) sameTitleMessage(time, title, place1, place2 string) string {
	meta := fmt.Sprintf("%s \n%s \n%s", time, place1, place2)

	if len(meta) + len(title) + 2 > messageMaxLen {
		title = strings.Trim(title[0:titleMaxLen], " \n") + "…"
	}
	if len(meta) + len(title) + 2 > messageMaxLen {
		metaLen := messageMaxLen - len(title) - 5
		meta = strings.Trim(meta[0:metaLen], " \n") + "…"
	}
	return fmt.Sprintf("%s \n%s", meta, title)
}

func (a CrossUpdateArg) diffTitleMessage(time, new_title, place1, place2, old_title string) string {
	meta := fmt.Sprintf("%s \n%s \n%s", time, place1, place2)
	title := fmt.Sprintf("\"%s\"\nchanged from \"%s\"", new_title, old_title)

	if len(meta) + len(title) + 2 > messageMaxLen {
		new_title = strings.Trim(new_title[0:newTitleMaxLen], " \n") + "…"
		title = fmt.Sprintf("\"%s\"\nchanged from \"%s\"", new_title, old_title)
	}
	if len(meta) + len(title) + 2 > messageMaxLen {
		old_title = strings.Trim(old_title[0:titleMaxLen], " \n") + "…"
		title = fmt.Sprintf("\"%s\"\nchanged from \"%s\"", new_title, old_title)
	}
	if len(meta) + len(title) + 2 > messageMaxLen {
		metaLen := messageMaxLen - len(title) - 5
		meta = strings.Trim(meta[0:metaLen], " \n") + "…"
	}
	return fmt.Sprintf("%s \n%s", meta, title)
}

type Cross_update struct {
	Config *Config
	Client *gobus.Client
}


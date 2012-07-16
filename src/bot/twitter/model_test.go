package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestTweetModel(t *testing.T) {
	dm := `{"direct_message":{"id_str":"223662005347823616","recipient_id":491159882,"created_at":"Fri Jul 13 06:15:52 +0000 2012","sender_screen_name":"googollee","recipient":{"id":491159882,"statuses_count":1,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme1\/bg.png","default_profile_image":true,"favourites_count":0,"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme1\/bg.png","following":true,"profile_link_color":"0084B4","utc_offset":null,"name":"Li Zhaohai","notifications":false,"profile_use_background_image":true,"contributors_enabled":false,"geo_enabled":false,"protected":false,"profile_text_color":"333333","id_str":"491159882","default_profile":true,"profile_image_url":"http:\/\/a0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","show_all_inline_media":false,"followers_count":4,"profile_sidebar_border_color":"C0DEED","description":null,"url":null,"screen_name":"lzh429","profile_background_tile":false,"created_at":"Mon Feb 13 09:52:49 +0000 2012","listed_count":0,"friends_count":3,"lang":"en","profile_sidebar_fill_color":"DDEEF6","verified":false,"time_zone":null,"is_translator":false,"location":null,"profile_image_url_https":"https:\/\/si0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","profile_background_color":"C0DEED","follow_request_sent":false},"sender":{"id":56591660,"statuses_count":43821,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/profile_background_images\/371844736\/2_mg_black2.gif","default_profile_image":false,"favourites_count":39,"profile_background_image_url":"http:\/\/a0.twimg.com\/profile_background_images\/371844736\/2_mg_black2.gif","following":false,"profile_link_color":"c45858","utc_offset":28800,"name":"Googol | \uff0f\u4eba\u25d5 \u203f\u203f \u25d5\u4eba\uff3c","notifications":false,"profile_use_background_image":true,"contributors_enabled":false,"geo_enabled":true,"protected":false,"profile_text_color":"333333","id_str":"56591660","default_profile":false,"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2391387815\/touemfxmkdqy9uounkj1_normal.jpeg","show_all_inline_media":false,"followers_count":2032,"profile_sidebar_border_color":"dfc1eb","description":"\uff0f\u4eba\u25d5 \u203f\u203f \u25d5\u4eba\uff3c  #nowplaying G\u5f26\u4e0a\u7684\u548f\u53f9\u8c03","url":"http:\/\/gplus.to\/googollee","screen_name":"googollee","profile_background_tile":true,"created_at":"Tue Jul 14 03:25:45 +0000 2009","listed_count":155,"friends_count":905,"lang":"en","profile_sidebar_fill_color":"bdbdbd","verified":false,"time_zone":"Beijing","is_translator":false,"location":"Beijing, China","profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2391387815\/touemfxmkdqy9uounkj1_normal.jpeg","profile_background_color":"ebc1de","follow_request_sent":false},"recipient_screen_name":"lzh429","sender_id":56591660,"id":223662005347823616,"text":"fdasfdasf"}}`
	tweet := `{"in_reply_to_status_id_str":null,"contributors":null,"place":null,"in_reply_to_screen_name":"lzh429","text":"@lzh429 fafdsf","favorited":false,"in_reply_to_user_id_str":"491159882","coordinates":null,"geo":null,"retweet_count":0,"created_at":"Fri Jul 13 06:16:23 +0000 2012","source":"web","in_reply_to_user_id":491159882,"in_reply_to_status_id":null,"retweeted":false,"id_str":"223662138227568640","truncated":false,"user":{"favourites_count":0,"friends_count":3,"profile_background_color":"C0DEED","following":null,"profile_background_tile":false,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme1\/bg.png","followers_count":4,"profile_image_url":"http:\/\/a0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","contributors_enabled":false,"geo_enabled":false,"created_at":"Mon Feb 13 09:52:49 +0000 2012","profile_sidebar_fill_color":"DDEEF6","description":null,"listed_count":0,"follow_request_sent":null,"time_zone":null,"url":null,"verified":false,"profile_sidebar_border_color":"C0DEED","default_profile":true,"show_all_inline_media":false,"is_translator":false,"notifications":null,"profile_use_background_image":true,"protected":false,"profile_image_url_https":"https:\/\/si0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","location":null,"id_str":"491159882","profile_text_color":"333333","name":"Li Zhaohai","statuses_count":2,"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme1\/bg.png","id":491159882,"default_profile_image":true,"lang":"en","utc_offset":null,"profile_link_color":"0084B4","screen_name":"lzh429"},"id":223662138227568640,"entities":{"user_mentions":[{"indices":[0,7],"name":"Li Zhaohai","id_str":"491159882","id":491159882,"screen_name":"lzh429"}],"urls":[],"hashtags":[]}}`

	{
		buf := bytes.NewBufferString(dm)
		d := json.NewDecoder(buf)
		var tw Tweet
		e := d.Decode(&tw)
		fmt.Println(e)
		i := (&tw).ToInput()
		expect := "56591660"
		if i.ID != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ID)
		}
		expect = "googollee"
		if i.ScreenName != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ScreenName)
		}
		expect = ""
		if i.Iom != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Iom)
		}
		expect = "2012-07-13 06:15:52 +0000"
		if i.CreatedAt != expect {
			t.Errorf("expect: %s, got: %s", expect, i.CreatedAt)
		}
		expect = "fdasfdasf"
		if i.Text != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Text)
		}
	}

	{
		buf := bytes.NewBufferString(tweet)
		d := json.NewDecoder(buf)
		var tw Tweet
		e := d.Decode(&tw)
		fmt.Println(e)
		i := (&tw).ToInput()
		expect := "491159882"
		if i.ID != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ID)
		}
		expect = "lzh429"
		if i.ScreenName != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ScreenName)
		}
		expect = ""
		if i.Iom != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Iom)
		}
		expect = "2012-07-13 06:16:23 +0000"
		if i.CreatedAt != expect {
			t.Errorf("expect: %s, got: %s", expect, i.CreatedAt)
		}
		expect = "@lzh429 fafdsf"
		if i.Text != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Text)
		}
	}
}

func TestStreamingReader(t *testing.T) {
	dm := `{"direct_message":{"id_str":"223662005347823616","recipient_id":491159882,"created_at":"Fri Jul 13 06:15:52 +0000 2012","sender_screen_name":"googollee","recipient":{"id":491159882,"statuses_count":1,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme1\/bg.png","default_profile_image":true,"favourites_count":0,"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme1\/bg.png","following":true,"profile_link_color":"0084B4","utc_offset":null,"name":"Li Zhaohai","notifications":false,"profile_use_background_image":true,"contributors_enabled":false,"geo_enabled":false,"protected":false,"profile_text_color":"333333","id_str":"491159882","default_profile":true,"profile_image_url":"http:\/\/a0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","show_all_inline_media":false,"followers_count":4,"profile_sidebar_border_color":"C0DEED","description":null,"url":null,"screen_name":"lzh429","profile_background_tile":false,"created_at":"Mon Feb 13 09:52:49 +0000 2012","listed_count":0,"friends_count":3,"lang":"en","profile_sidebar_fill_color":"DDEEF6","verified":false,"time_zone":null,"is_translator":false,"location":null,"profile_image_url_https":"https:\/\/si0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","profile_background_color":"C0DEED","follow_request_sent":false},"sender":{"id":56591660,"statuses_count":43821,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/profile_background_images\/371844736\/2_mg_black2.gif","default_profile_image":false,"favourites_count":39,"profile_background_image_url":"http:\/\/a0.twimg.com\/profile_background_images\/371844736\/2_mg_black2.gif","following":false,"profile_link_color":"c45858","utc_offset":28800,"name":"Googol | \uff0f\u4eba\u25d5 \u203f\u203f \u25d5\u4eba\uff3c","notifications":false,"profile_use_background_image":true,"contributors_enabled":false,"geo_enabled":true,"protected":false,"profile_text_color":"333333","id_str":"56591660","default_profile":false,"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2391387815\/touemfxmkdqy9uounkj1_normal.jpeg","show_all_inline_media":false,"followers_count":2032,"profile_sidebar_border_color":"dfc1eb","description":"\uff0f\u4eba\u25d5 \u203f\u203f \u25d5\u4eba\uff3c  #nowplaying G\u5f26\u4e0a\u7684\u548f\u53f9\u8c03","url":"http:\/\/gplus.to\/googollee","screen_name":"googollee","profile_background_tile":true,"created_at":"Tue Jul 14 03:25:45 +0000 2009","listed_count":155,"friends_count":905,"lang":"en","profile_sidebar_fill_color":"bdbdbd","verified":false,"time_zone":"Beijing","is_translator":false,"location":"Beijing, China","profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2391387815\/touemfxmkdqy9uounkj1_normal.jpeg","profile_background_color":"ebc1de","follow_request_sent":false},"recipient_screen_name":"lzh429","sender_id":56591660,"id":223662005347823616,"text":"fdasfdasf"}}`
	tweet := `{"in_reply_to_status_id_str":null,"contributors":null,"place":null,"in_reply_to_screen_name":"lzh429","text":"@lzh429 fafdsf","favorited":false,"in_reply_to_user_id_str":"491159882","coordinates":null,"geo":null,"retweet_count":0,"created_at":"Fri Jul 13 06:16:23 +0000 2012","source":"web","in_reply_to_user_id":491159882,"in_reply_to_status_id":null,"retweeted":false,"id_str":"223662138227568640","truncated":false,"user":{"favourites_count":0,"friends_count":3,"profile_background_color":"C0DEED","following":null,"profile_background_tile":false,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme1\/bg.png","followers_count":4,"profile_image_url":"http:\/\/a0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","contributors_enabled":false,"geo_enabled":false,"created_at":"Mon Feb 13 09:52:49 +0000 2012","profile_sidebar_fill_color":"DDEEF6","description":null,"listed_count":0,"follow_request_sent":null,"time_zone":null,"url":null,"verified":false,"profile_sidebar_border_color":"C0DEED","default_profile":true,"show_all_inline_media":false,"is_translator":false,"notifications":null,"profile_use_background_image":true,"protected":false,"profile_image_url_https":"https:\/\/si0.twimg.com\/sticky\/default_profile_images\/default_profile_3_normal.png","location":null,"id_str":"491159882","profile_text_color":"333333","name":"Li Zhaohai","statuses_count":2,"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme1\/bg.png","id":491159882,"default_profile_image":true,"lang":"en","utc_offset":null,"profile_link_color":"0084B4","screen_name":"lzh429"},"id":223662138227568640,"entities":{"user_mentions":[{"indices":[0,7],"name":"Li Zhaohai","id_str":"491159882","id":491159882,"screen_name":"lzh429"}],"urls":[],"hashtags":[]}}`
	ender := "\r\n"
	streaming := ""
	for _, s := range []string{dm, tweet} {
		streaming = streaming + s + ender
	}

	buf := bytes.NewBufferString(streaming)
	reader := NewStreamingReader(buf)

	{
		d, err := reader.ReadTweet()
		t.Logf("%v", err)
		t.Logf("%+v", d)
		i := d.ToInput()
		expect := "56591660"
		if i.ID != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ID)
		}
		expect = "googollee"
		if i.ScreenName != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ScreenName)
		}
		expect = ""
		if i.Iom != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Iom)
		}
		expect = "2012-07-13 06:15:52 +0000"
		if i.CreatedAt != expect {
			t.Errorf("expect: %s, got: %s", expect, i.CreatedAt)
		}
		expect = "fdasfdasf"
		if i.Text != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Text)
		}
	}

	{
		tw, err := reader.ReadTweet()
		t.Logf("%v", err)
		t.Logf("%+v", tw)
		i := tw.ToInput()
		expect := "491159882"
		if i.ID != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ID)
		}
		expect = "lzh429"
		if i.ScreenName != expect {
			t.Errorf("expect: %s, got: %s", expect, i.ScreenName)
		}
		expect = ""
		if i.Iom != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Iom)
		}
		expect = "2012-07-13 06:16:23 +0000"
		if i.CreatedAt != expect {
			t.Errorf("expect: %s, got: %s", expect, i.CreatedAt)
		}
		expect = "@lzh429 fafdsf"
		if i.Text != expect {
			t.Errorf("expect: %s, got: %s", expect, i.Text)
		}
	}
}

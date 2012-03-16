#!/usr/bin/env ruby
# coding: utf-8

require 'redis'

queue = "resque:queue:twitter"

job1 = "{\"class\":\"twitter_job\",\"args\":{\"title\":\"abc\",\"description\":\"fdafadsf\",\"begin_at\":{\"time\":\"8:00 AM\",\"date\":\"Mar 21, 2012\",\"datetime\":\"8:00 AM, Mar 21, 2012\",\"time_type\":\"\"},\"time_type\":\"Anytime\",\"place_line1\":\"here\",\"place_line2\":\"\",\"cross_id\":100092,\"cross_id_base62\":\"q2o\",\"invitation_id\":\"300\",\"token\":\"f80db16449e49abd8b7f493af4b2671d\",\"identity_id\":\"174\",\"host_identity_id\":174,\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"host_identity\":{\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\"},\"rsvp_status\":1,\"by_identity\":{\"id\":\"174\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"provider\":\"twitter\"},\"to_identity\":{\"id\":\"174\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"provider\":\"twitter\"},\"to_identity_time_zone\":\"+08:00\",\"invitations\":[{\"invitation_id\":\"300\",\"state\":1,\"by_identity_id\":\"174\",\"token\":\"f80db16449e49abd8b7f493af4b2671d\",\"updated_at\":\"2012-03-16 19:42:00\",\"identity_id\":\"174\",\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"identities\":[{\"identity_id\":\"174\",\"status\":\"3\",\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\"}],\"user_id\":131},{\"invitation_id\":\"301\",\"state\":0,\"by_identity_id\":\"174\",\"token\":\"eed03ac1a2b35745003bc7205028957d\",\"updated_at\":\"2012-03-16 19:42:00\",\"identity_id\":\"175\",\"provider\":\"twitter\",\"external_identity\":\"@lzh429@twitter\",\"name\":\"Li Zhaohai\",\"bio\":\"\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"external_username\":\"lzh429\",\"user_id\":0}]},\"id\":\"6245b470aef1c287e36ba83d4e7c351f\"}"

job2 = "{\"class\":\"twitter_job\",\"args\":{\"title\":\"abc\",\"description\":\"fdafadsf\",\"begin_at\":{\"time\":\"8:00 AM\",\"date\":\"Mar 21, 2012\",\"datetime\":\"8:00 AM, Mar 21, 2012\",\"time_type\":\"\"},\"time_type\":\"Anytime\",\"place_line1\":\"here\",\"place_line2\":\"\",\"cross_id\":100092,\"cross_id_base62\":\"q2o\",\"invitation_id\":\"301\",\"token\":\"eed03ac1a2b35745003bc7205028957d\",\"identity_id\":\"175\",\"host_identity_id\":174,\"provider\":\"twitter\",\"external_identity\":\"@lzh429@twitter\",\"name\":\"Li Zhaohai\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"host_identity\":{\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\"},\"rsvp_status\":0,\"by_identity\":{\"id\":\"174\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"provider\":\"twitter\"},\"to_identity\":{\"id\":\"175\",\"external_identity\":\"\",\"name\":\"Li Zhaohai\",\"bio\":\"\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"external_username\":\"lzh429\",\"provider\":\"twitter\"},\"to_identity_time_zone\":\"+08:00\",\"invitations\":[{\"invitation_id\":\"300\",\"state\":1,\"by_identity_id\":\"174\",\"token\":\"f80db16449e49abd8b7f493af4b2671d\",\"updated_at\":\"2012-03-16 19:42:00\",\"identity_id\":\"174\",\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"identities\":[{\"identity_id\":\"174\",\"status\":\"3\",\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\"}],\"user_id\":131},{\"invitation_id\":\"301\",\"state\":0,\"by_identity_id\":\"174\",\"token\":\"eed03ac1a2b35745003bc7205028957d\",\"updated_at\":\"2012-03-16 19:42:00\",\"identity_id\":\"175\",\"provider\":\"twitter\",\"external_identity\":\"@lzh429@twitter\",\"name\":\"Li Zhaohai\",\"bio\":\"\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"external_username\":\"lzh429\",\"user_id\":0}]},\"id\":\"ab41f20613af73deb0e6a9df02886724\"}"

job3 = "{\"class\":\"twitter_job\",\"args\":{\"title\":\"Aabc\",\"description\":\"fdafadsf\",\"begin_at\":{\"time\":\"8:00 AM\",\"date\":\"Mar 21, 2012\",\"datetime\":\"8:00 AM, Mar 21, 2012\",\"time_type\":\"\"},\"place_line1\":\"here\",\"place_line2\":\"\",\"cross_id\":100092,\"cross_id_base62\":\"q2o\",\"invitation_id\":\"303\",\"token\":\"297bac94b63b1b5381ae46f57d94ab7f\",\"identity_id\":\"175\",\"host_identity_id\":174,\"provider\":\"twitter\",\"external_identity\":\"\",\"name\":\"Li Zhaohai\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"host_identity\":{\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\"},\"rsvp_status\":0,\"to_identity_time_zone\":\"+08:00\",\"by_identity\":{\"id\":\"174\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"provider\":\"twitter\"},\"to_identity\":{\"id\":\"175\",\"external_identity\":\"\",\"name\":\"Li Zhaohai\",\"bio\":\"\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"external_username\":\"lzh429\",\"provider\":\"twitter\"},\"invitations\":[{\"invitation_id\":\"300\",\"state\":1,\"token\":\"f80db16449e49abd8b7f493af4b2671d\",\"updated_at\":\"2012-03-16 19:42:00\",\"by_identity_id\":\"174\",\"identity_id\":\"174\",\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\",\"identities\":[{\"identity_id\":\"174\",\"status\":\"3\",\"provider\":\"twitter\",\"external_identity\":\"twitter_56591660\",\"name\":\"Googol | ／人◕ ‿‿ ◕人＼\",\"bio\":\"／人◕ ‿‿ ◕人＼ #nowplaying G弦上的咏叹调\",\"avatar_file_name\":\"http://a0.twimg.com/profile_images/1896720358/Untitled_reasonably_small.jpg\",\"external_username\":\"googollee\"}],\"user_id\":131},{\"invitation_id\":\"303\",\"state\":0,\"token\":\"297bac94b63b1b5381ae46f57d94ab7f\",\"updated_at\":\"2012-03-16 20:41:44\",\"by_identity_id\":\"174\",\"identity_id\":\"175\",\"provider\":\"twitter\",\"external_identity\":\"\",\"name\":\"Li Zhaohai\",\"bio\":\"\",\"avatar_file_name\":\"http://a0.twimg.com/sticky/default_profile_images/default_profile_3_reasonably_small.png\",\"external_username\":\"lzh429\",\"user_id\":0}]},\"id\":\"24b1911927edac27a0d89319e9d96477\"}"

redis = Redis.new
redis.rpush queue, job1
redis.rpush queue, job2
redis.rpush queue, job3

curl "http://127.0.0.1:23333/v3/notifier/cross/digest" -d '[{
  "to":{
    "external_id":"googollee@gmail.com",
    "external_username":"googollee@gmail.com",
    "provider":"email",
    "name":"Googol Lee",
    "token":"141e2197cd155316dddd54044c5f0e3b",
    "identity_id":572,
    "user_id":384},
  "cross_id":100354,
  "updated_at":"2013-03-11 00:00:00"}]'

curl "http://127.0.0.1:23333/v3/notifier/cross/invitation" -d '{
  "to":{
    "external_id":"googollee@gmail.com",
    "external_username":"googollee@gmail.com",
    "provider":"email",
    "name":"Googol Lee",
    "token":"141e2197cd155316dddd54044c5f0e3b",
    "identity_id":572,
    "user_id":384},
  "by":{
    "external_id":"googollee@gmail.com",
    "external_username":"googollee@gmail.com",
    "provider":"email",
    "name":"Googol Lee",
    "token":"141e2197cd155316dddd54044c5f0e3b",
    "identity_id":572,
    "user_id":384},
  "cross":{"title":"Meet Googol","description":"\u518d\u6b21\u6d4b\u8bd5\u518d\u6b21\u6d4b\u8bd5","time":{"begin_at":{"date_word":"","date":"2013-04-17","time_word":"","time":"","timezone":"+08:00 CST","id":0,"type":"EFTime"},"origin":"2013-04-17","outputformat":0,"id":0,"type":"CrossTime"},"place":{"title":"aaaaa","description":"","lng":"","lat":"","provider":"","external_id":"","created_at":"2012-12-19 17:50:52 +0000","updated_at":"2013-04-28 15:05:23 +0000","id":289,"type":"Place"},"attribute":{"state":"published","closed":false},"exfee":{"invitations":[{"identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"invited_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"by_identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"updated_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"rsvp_status":"ACCEPTED","via":"","created_at":"2012-12-19 17:50:52 +0000","updated_at":"2012-12-19 17:50:52 +0000","token":"141e2197cd155316dddd54044c5f0e3b","host":true,"mates":0,"remark":[],"id":1198,"type":"invitation"},{"identity":{"name":"googollee","nickname":"","bio":"","provider":"twitter","connected_user_id":-573,"external_id":"","external_username":"googollee","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=googollee","created_at":"2012-12-19 17:51:02 +0000","updated_at":"2012-12-19 17:53:17 +0000","order":999,"unreachable":true,"id":573,"type":"identity"},"invited_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"by_identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"updated_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"rsvp_status":"NORESPONSE","via":"","created_at":"2012-12-19 17:51:02 +0000","updated_at":"2012-12-19 17:51:02 +0000","token":"911c2e88d4d7b63cfc76a04a0e51fafb","host":false,"mates":0,"remark":[],"id":1199,"type":"invitation"},{"identity":{"name":"googollee","nickname":"","bio":"","provider":"email","connected_user_id":475,"external_id":"googollee@gmail.com","external_username":"googollee@gmail.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/15b7fc1b101ee289b81678812781aea6","created_at":"2012-12-21 13:05:46 +0000","updated_at":"2013-03-06 17:05:37 +0000","order":999,"unreachable":false,"id":574,"type":"identity"},"invited_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"by_identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"updated_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"rsvp_status":"NORESPONSE","via":"","created_at":"2012-12-21 13:05:46 +0000","updated_at":"2012-12-21 13:05:46 +0000","token":"80269cf3f80057d2835ed6f560ae524c","host":false,"mates":0,"remark":[],"id":1200,"type":"invitation"},{"identity":{"name":"891","nickname":"","bio":"","provider":"phone","connected_user_id":-614,"external_id":"+8613488802891","external_username":"+8613488802891","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=891","created_at":"2013-03-24 18:17:16 +0000","updated_at":"2013-03-24 18:17:16 +0000","order":999,"unreachable":false,"id":614,"type":"identity"},"invited_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"by_identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"updated_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"rsvp_status":"NORESPONSE","via":"","created_at":"2013-03-24 18:17:16 +0000","updated_at":"2013-03-24 18:17:16 +0000","token":"3fd83413d0ed25e53f6f3ce535a4ba4e","host":false,"mates":0,"remark":[],"id":1303,"type":"invitation"},{"identity":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"invited_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"by_identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"updated_by":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"rsvp_status":"NORESPONSE","via":"","created_at":"2013-03-24 16:08:19 +0000","updated_at":"2013-03-24 16:08:19 +0000","token":"a23a1caf2da597e378ef72d7006110fc","host":false,"mates":0,"remark":[],"id":1302,"type":"invitation"},{"identity":{"name":"leaskh","nickname":"","bio":"","provider":"phone","connected_user_id":327,"external_id":"+8618675530413","external_username":"+8618675530413","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=%2B8618675530413","created_at":"2013-02-22 15:16:48 +0000","updated_at":"2013-03-12 03:30:02 +0000","order":999,"unreachable":false,"id":586,"type":"identity"},"invited_by":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"by_identity":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"updated_by":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"rsvp_status":"NORESPONSE","via":"","created_at":"2013-04-01 12:03:27 +0000","updated_at":"2013-04-02 01:40:31 +0000","token":"f7f8fe7bea4e14ded6e15fc20a7ab8b3","host":false,"mates":0,"remark":[],"id":1318,"type":"invitation"},{"identity":{"name":"L","nickname":"","bio":"","provider":"phone","connected_user_id":-640,"external_id":"+8611111111","external_username":"+8611111111","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=L","created_at":"2013-04-02 15:46:37 +0000","updated_at":"2013-04-02 15:46:37 +0000","order":999,"unreachable":false,"id":640,"type":"identity"},"invited_by":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"by_identity":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:* Connection #0 to host api.panda.0d0f.com left intact > \/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"updated_by":{"name":"Leask Huang","nickname":"","bio":"","provider":"email","connected_user_id":379,"external_id":"i@leaskh.com","external_username":"i@leaskh.com","avatar_filename":"http:\/\/www.gravatar.com\/avatar\/5b53fb71b6f36f46fe9cb14eb5acd56f","created_at":"2012-12-11 15:15:35 +0000","updated_at":"2013-04-28 02:32:55 +0000","order":0,"unreachable":false,"id":569,"type":"identity"},"rsvp_status":"NORESPONSE","via":"","created_at":"2013-04-02 15:46:37 +0000","updated_at":"2013-04-02 15:46:37 +0000","token":"9f1fd2a1bba790c726c8739aae4a259b","host":false,"mates":0,"remark":[],"id":1331,"type":"invitation"}],"items":7,"total":7,"accepted":1,"name":"Meet Googol","id":110220,"type":"exfee","hosts":[384],"updated_at":"2013-04-28 15:05:23 +0000"},"widget":[{"image":"RedRiverValley.jpg","widget_id":0,"id":0,"type":"Background"}],"conversation_count":0,"relative":[],"type":"Cross","created_at":"2012-12-19 17:50:52 +0000","by_identity":{"name":"Googol","nickname":"","bio":"","provider":"email","connected_user_id":384,"external_id":"googollee@163.com","external_username":"googollee@163.com","avatar_filename":"http:\/\/api.panda.0d0f.com\/v2\/avatar\/default?name=Googol","created_at":"2012-12-17 18:09:51 +0000","updated_at":"2013-01-14 17:03:13 +0000","order":999,"unreachable":false,"id":572,"type":"identity"},"id":100354,"updated_at":"2013-04-28 15:05:23 +0000","updated":{"conversation":{"updated_at":"2013-03-24 10:15:08 +0000"},"exfee":{"updated_at":"2013-04-02 07:46:37 +0000"},"background":{"updated_at":"2013-04-09 13:14:26 +0000"},"title":{"updated_at":"2013-04-28 07:03:56 +0000"},"description":{"updated_at":"2013-04-28 07:03:53 +0000"},"time":{"updated_at":"2013-04-28 15:03:45 +0000"},"place":{"updated_at":"2013-04-28 07:05:23 +0000"}}}}'
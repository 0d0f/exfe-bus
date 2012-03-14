#!/usr/bin/env ruby
# coding: utf-8

require 'json'
require 'redis'

queue = "gobus:queue:twitter:directmessage"
tweet_name = {
  :Id => "#{queue}:1",
  :Arg => {
    :ClientToken => "VC3OxLBNSGPLOZ2zkgisA",
    :ClientSecret => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
    :AccessToken => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
    :AccessSecret => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
    :Message => "就是测个试Name",
    :ToUserName => "lzh429",
  },
  :NeedReply => true,
}

tweet_id = {
  :Id => "#{queue}:2",
  :Arg => {
    :ClientToken => "VC3OxLBNSGPLOZ2zkgisA",
    :ClientSecret => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
    :AccessToken => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
    :AccessSecret => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
    :Message => "就是测个试ID",
    :ToUserId => "56591660",
  },
  :NeedReply => true,
}

redis = Redis.new
redis.rpush queue, tweet_name.to_json
redis.rpush queue, tweet_id.to_json

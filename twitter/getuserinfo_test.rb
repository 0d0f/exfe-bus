#!/usr/bin/env ruby
# coding: utf-8

require 'json'
require 'redis'

queue = "gobus:queue:twitter:userinfo"
job1 = {
  :Id => "#{queue}:1",
  :Data => {
    :ClientToken => "VC3OxLBNSGPLOZ2zkgisA",
    :ClientSecret => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
    :AccessToken => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
    :AccessSecret => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
    :ScreenName => "lzh429",
  }
}

job2 = {
  :Id => "#{queue}:2",
  :Data => {
    :ClientToken => "VC3OxLBNSGPLOZ2zkgisA",
    :ClientSecret => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
    :AccessToken => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
    :AccessSecret => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
    :UserId => "491159882",
  }
}

redis = Redis.new
redis.rpush queue, job1.to_json
redis.rpush queue, job2.to_json
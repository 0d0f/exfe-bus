#!/usr/bin/env ruby
# coding: utf-8

dataWithName = {
  "ClientToken" => "VC3OxLBNSGPLOZ2zkgisA",
  "ClientSecret" => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
  "AccessToken" => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
  "AccessSecret" => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
  "Message" => "就是测个试Name",
  "ToUserName" => "lzh429",
}

dataWithId = {
  "ClientToken" => "VC3OxLBNSGPLOZ2zkgisA",
  "ClientSecret" => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
  "AccessToken" => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
  "AccessSecret" => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
  "Message" => "就是测个试Id",
  "ToUserId" => "491159882",
}

require 'redis'
require 'json'

redis = Redis.new
redis.rpush 'resque:twitter:directmessage', dataWithName.to_json
redis.rpush 'resque:twitter:directmessage', dataWithId.to_json

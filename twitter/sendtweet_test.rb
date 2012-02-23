#!/usr/bin/env ruby
# coding: utf-8

data = {
  "ClientToken" => "VC3OxLBNSGPLOZ2zkgisA",
  "ClientSecret" => "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
  "AccessToken" => "491159882-urND5ZaHmUPWNgvpr5coIifkApmKsmjGtX69Bn51",
  "AccessSecret" => "5kwJdfqd6xL93BvPisYaRVzk5VlOEMhQwAk2aPMxy6s",
  "Tweet" => "就是测个试",
}

require 'redis'
require 'json'

redis = Redis.new
redis.rpush 'resque:twitter:tweet', data.to_json

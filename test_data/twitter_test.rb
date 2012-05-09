#!/usr/bin/env ruby

require 'redis'

redis = Redis.new
i = 0

open("./test_data/twitter_test.data").each do |l|
  redis.rpush "gobus:queue:twitter_job", l
  i += 1
end

redis.set "gobus:queue:twitter_job:idcount", i

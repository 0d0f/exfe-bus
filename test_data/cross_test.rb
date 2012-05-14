#!/usr/bin/env ruby

require 'redis'

redis = Redis.new
i = 0

open("./test_data/cross_test.data").each do |l|
  next if l.length == 0
  redis.rpush "gobus:queue:cross", l
  i += 1
end

redis.set "gobus:queue:cross:idcount", i

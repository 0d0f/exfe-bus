#!/usr/bin/env ruby
# coding: utf-8

require "redis"
require 'json'

queue = "gobus:queue:mail:sender"
job = {
  :Id => "#{queue}:1",
  :Data => {
    :To => [
      {
        :Mail => "lzh@exfe.com",
        :Name => "Li Zhaohai",
      },
      {
        :Mail => "googollee@hotmail.com",
        :Name => "Googol Lee",
      }
    ],
    :From => {
      :Mail => "googollee@gmail.com",
      :Name => "Googol Lee",
    },
    :Subject => "A Test",
    :Message => "Just a mail test",
  },
  :NeedReturn => true
}

redis = Redis.new
redis.rpush queue, job.to_json

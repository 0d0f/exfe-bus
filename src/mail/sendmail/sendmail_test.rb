#!/usr/bin/env ruby
# coding: utf-8

require "redis"
require 'json'

queue = "gobus:queue:mail:sender"
job = {
  :Id => "#{queue}:1",
  :Arg => {
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
    :Text => "Just a mail test",
    :Html => "<h1>Just mail test</h1>",
    :FileParts => [
      {
        :Name => "test.txt",
        :Content => "dGVzdCBkYXRh",
      },
    ]
  },
  :NeedReply => true
}

redis = Redis.new
redis.rpush queue, job.to_json

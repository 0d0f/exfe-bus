#!/usr/bin/env ruby
# coding: utf-8

require 'json'
require 'redis'

job = {
  "class"=>"welcomeandactivecode_job",
  "args"=>
  {
    "identityid"=>166,
    "external_identity"=>"googollee@hotmail.com",
    "name"=>"Googol",
    "avatar_file_name"=> "http://www.gravatar.com/avatar/fc1d342bf78fffe115d867168598c9a5?d=http%3A%2F%2Fimages.exfe.com%2Fweb%2F80_80_default.png",
    "activecode"=>"e84dcb54bce7e038fcb6700f20939e421331105824",
    "token"=> "4d6a41344d544d334e2b6d43347544684b75567043513549424e47487852526930356e346c5769596c6554764434345551326537797278356f485877766b773769745070714c5866666a6f7133374477724b6a4c674d374338533062676c33356e486545415032417846516957656a554134372b6f3134585a477445437a57456e4b5a63786854644e704c71616d774f675171383d"
  },
  "id"=>"ffb6fa51c03ee68068b5d2ee30c35d5b"
}

queue = "resque:queue:email"

redis = Redis.new
redis.rpush queue, job.to_json

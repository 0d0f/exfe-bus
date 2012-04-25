<?php
require_once 'Redisent.php';

class Gobus
{
    public static $redis = null;

    public static function setBackend($server)
    {
        list($host, $port) = explode(':', $server);
        self::$redis = new Redisent($host, $port);
    }

    public static function send($queue_name, $method, $arg, $max_retry = 5)
    {
        $queue = "gobus:queue:" . $queue_name;
        $idcount = $queue . ":idcount";
        $id = self::$redis->incr($idcount);

        $meta = array(
            'id' => $queue . ":" . $id,
            'method' => $method,
            'arg' => $arg,
            'maxRetry' => intval($max_retry),
            "needReply" => false
        );

        $data = json_encode($method).json_encode($meta);
        self::$redis->rpush($queue, $data);
    }
}
?>

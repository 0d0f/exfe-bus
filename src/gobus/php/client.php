<?php
require_once 'Redisent.php';

class Gobus
{
    public static $host = null;
    public static $port = null;

    public static function setBackend($server)
    {
        list(self::$host, self::$port) = explode(':', $server);
    }

    public static function send($queue_name, $method, $arg, $max_retry = 5)
    {
        $redis = new Redis();
        $redis->connect(self::$host, self::$port);

        $queue = "gobus:queue:" . $queue_name;
        $idcount = $queue . ":idcount";
        $id = $redis->incr($idcount);

        $meta = array(
            'id' => $queue . ":" . $id,
            'method' => $method,
            'arg' => $arg,
            'maxRetry' => intval($max_retry),
            "needReply" => false
        );

        $data = json_encode($method).json_encode($meta);
        $redis->rPush($queue, $data);
    }
}
?>

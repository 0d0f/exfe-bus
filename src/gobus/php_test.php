<?php
require_once 'php/client.php';


Gobus::setBackend("127.0.0.1:6379");
Gobus::send("php", "Batch", 3);
?>

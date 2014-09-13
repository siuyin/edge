<?php

$context = new ZMQContext();

//  Socket to talk to server
echo "Connecting to server.\n";
$requester = new ZMQSocket($context, ZMQ::SOCKET_PUSH);
$requester->connect("ipc:///var/www/socks/sms.ipc");

$requester->send("{\"dest\":\"+6591720889\",\"msg\":\"brown fox\"}");
?>

<?php
/*
*  Hello World client
*  Connects REQ socket to tcp://localhost:5555
*  Sends "Hello" to server, expects "World" back
* @author Ian Barber <ian(dot)barber(at)gmail(dot)com>
*/

$context = new ZMQContext();

//  Socket to talk to server
echo "Connecting to server.\n";
$requester = new ZMQSocket($context, ZMQ::SOCKET_REQ);
$requester->connect("ipc:///var/www/socks/buy.ipc");

    $requester->send("{\"user_id\":1,\"item_id\":2}");

    $reply = $requester->recv();
    printf ("Received reply [%s]\n", $reply);
?>

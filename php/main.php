<?php

$csvString = file_get_contents("./syukujitsu.csv");
$datas = explode(",", $csvString);

var_dump($datas[2]);
// Topic: Variables and Arrays
<?php

$name = "toofan";
$version = 1.0;
$active = true;

$scores = [
    "alice" => 92,
    "bob"   => 85,
    "eve"   => 97,
];

foreach ($scores as $user => $score) {
    echo "$user: $score\n";
}

echo count($scores) . " users\n";

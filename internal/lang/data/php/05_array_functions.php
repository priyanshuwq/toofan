// Topic: Array Functions and Closures

$users = [
    ["name" => "Alice", "score" => 92],
    ["name" => "Bob",   "score" => 45],
    ["name" => "Eve",   "score" => 78],
];

$passing = array_filter($users, fn($u) => $u["score"] >= 60);

$names = array_map(fn($u) => strtoupper($u["name"]), $passing);

usort($passing, fn($a, $b) => $b["score"] <=> $a["score"]);

foreach ($passing as $u) {
    echo "{$u['name']}: {$u['score']}\n";
}

echo implode(", ", $names) . "\n";

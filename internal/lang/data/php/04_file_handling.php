// Topic: File Handling and JSON

function loadConfig(string $path): array {
    if (!file_exists($path)) {
        return [];
    }
    $raw = file_get_contents($path);
    return json_decode($raw, associative: true) ?? [];
}

function saveConfig(string $path, array $data): void {
    file_put_contents($path, json_encode($data, JSON_PRETTY_PRINT));
}

$config = loadConfig("config.json");
$config["debug"] = false;
$config["version"] = "1.0.0";
saveConfig("config.json", $config);

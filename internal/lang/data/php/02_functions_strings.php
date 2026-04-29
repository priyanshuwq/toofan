// Topic: Functions and Strings

function greet(string $name, int $age): string {
    return "Hi, I am $name and I am $age years old.";
}

function slugify(string $text): string {
    $lower = strtolower($text);
    return str_replace(" ", "-", $lower);
}

echo greet("Alice", 25) . "\n";
echo slugify("Hello World") . "\n";

$parts = explode(",", "go,rust,php");
echo implode(" | ", $parts) . "\n";

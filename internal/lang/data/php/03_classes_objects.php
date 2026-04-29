// Topic: Classes and Objects

class User {
    public function __construct(
        public readonly string $name,
        public readonly int $age,
        public ?string $bio = null,
    ) {}

    public function greet(): string {
        $desc = $this->bio ?? "no bio";
        return "Hi, I am {$this->name} ({$this->age}) - $desc";
    }
}

$user = new User(name: "Alice", age: 25);
echo $user->greet() . "\n";

$admin = new User(name: "Bob", age: 30, bio: "System admin");
echo $admin->greet() . "\n";

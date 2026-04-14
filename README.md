<div align="center">

# toofan

**A minimal, lightning-fast typing test for your terminal**  
*Practice with english words or real code snippets. No browser, no account, everything stays local.*

<br>

<img src="assets/main.gif" alt="toofan demo" width="750">

<br>

[**Install**](#install) &nbsp;·&nbsp; [**FAQ**](#faq) &nbsp;·&nbsp; [**Contribute**](#contributing)

</div>

---

## Features

- **Two modes** : English words or real code snippets
- **Hand-written code lessons** : Practical, self-contained examples that teach a concept while you type
- **Infinite mode** : No timer, type the full snippet at your own pace
- **Lesson picker** : Choose what to practice instead of getting random code
- **Live WPM** : Speed updates in real time as you type
- **Personal bests** : Tracked per duration and mode
- **Activity map** : Track how consistent you are across days
- **Backup and restore** : Save your data and import it on any machine
- **Multiple themes** : Match your terminal aesthetic
- Press `?` inside the app for all available keybindings

## Code Snippets

Type real, practical code instead of random keywords. Each snippet is a small, self-contained example that builds muscle memory for `{}`, `=>`, `()`, `<-`, and all the symbols you actually use. Pick a lesson with `ctrl+o` or let it choose randomly. Set the timer to `∞` and the test ends when you finish the snippet.

<div align="center">
<img src="assets/code-snippets.png" width="560">
</div>

## Supported Languages (contribute to add more)

Go, JavaScript, Lua, Shell, Dart, etc.

<div align="center">
<img src="assets/languages.png" width="560">
</div>

## Lessons (for code mode)

Hand-written lessons organized by topic. Each one teaches a real concept while you type. Not generated, not random.

<div align="center">
<img src="assets/lession.png" width="560">
</div>

## Themes

I've always wanted my terminal to look good and match my setup. So toofan comes with multiple themes you can cycle through to match your aesthetic. Hit `ctrl+t` to switch.

<div align="center">
<img src="assets/theme.png" width="560">
</div>

## Profile

Your typing history, personal bests, rank, and activity map all in one place. Press `ctrl+p` to open it anytime.

<div align="center">
<img src="assets/profile-page.png" width="600">
</div>

## Install

**Quick Install (Mac / Linux):**
```bash
curl -fsSL https://raw.githubusercontent.com/vyrx-dev/toofan/master/install.sh | sh
```
*(If you set up your subdomain redirect, you can change this to `curl -fsSL https://toofan.vyrx.dev | sh`)*

<br>

**Build from Source (Requires Go installed):**
```sh
git clone https://github.com/vyrx-dev/toofan.git
cd toofan
go build -o toofan .
sudo mv toofan /usr/local/bin/toofan
```

### Package Managers

Coming soon: AUR (paru/yay), Homebrew, Nix.

## FAQ

<details>
<summary>How are stats calculated?</summary>

```text
raw      = total_chars / 5 / elapsed_minutes
wpm      = (total_chars - uncorrected_errors) / 5 / elapsed_minutes
accuracy = (total_chars - all_mistakes) / total_chars × 100
```

- **wpm** - your net speed. Every 5 characters count as one "word". Uncorrected mistakes are subtracted.
- **accuracy** - counts every wrong keystroke, even if you corrected it with backspace.
- **raw** - your gross speed before any penalty.
- **errors** - press `e` on the results page to see exactly which words you mistyped.
</details>

<details>
<summary>How are ranks decided?</summary>

Your rank is based on your average WPM over the last 10 word-mode tests:

```text
grandma      # below 30 wpm
noob         # 30 to 50 wpm
mid          # 50 to 80 wpm
tryhard      # 80 to 120 wpm
toofan       # 120+ wpm
```
</details>

<details>
<summary>Where are my files stored?</summary>

Everything lives in `~/.toofan/` as plain text files:

- `config.txt` : Your selected duration, mode, language, and theme
- `results.txt` : Every test result (date, wpm, accuracy, duration, mode)
- `pb.txt` : Your personal bests per mode and duration
</details>

<details>
<summary>Can I backup my data?</summary>

Yes. Press `ctrl+s` to save a backup and `ctrl+r` to restore from one. Backups are saved to `~/Toofan/` and can be moved between machines.
</details>

<details>
<summary>Does it work offline?</summary>

Yes. Everything runs locally and is embedded in the binary. No internet needed.
</details>

## Roadmap

- [ ] More language support (python, rust, c, typescript, etc.)
- [ ] Difficulty levels for english words
- [ ] AUR, Homebrew, Nix, and curl install
- [ ] Fix top pane alignment to match bottom panes in profile
- [ ] Proper documentation for AI and contributors to understand the project

## Contributing

- New snippets : Drop a file in `internal/lang/data/<language>/lessons/` and rebuild
- New languages : Just a folder with lesson files
- New themes : One Go file with a color palette
- Bug fixes and UX improvements

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) : TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) : Terminal styling

<br>

<div align="center">
<a href="https://www.star-history.com/#vyrx-dev/toofan&type=Date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=vyrx-dev/toofan&type=Date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=vyrx-dev/toofan&type=Date&legend=top-left" />
   <img alt="toofan Star History Chart" src="https://api.star-history.com/svg?repos=vyrx-dev/toofan&type=Date&legend=top-left" />
 </picture>
</a>
</div>

---

> **Enjoying toofan?** Consider dropping a ⭐ or sharing it online. A mention is always appreciated :)

# Totem Trainer — Design

## Overview

A terminal typing trainer purpose-built for the Totem 38-key split keyboard running Colemak DH. Built with Go + Bubble Tea (Charm ecosystem). Parses the actual ZMK keymap as source of truth.

## Goals

- Progressive learning: home row → full alpha → numbers → symbols → mixed real-world
- Adaptive per-key drilling (keybr-style) within each stage
- Persistent stats with historical trends and keyboard heatmap
- Flexible enough to work as the keymap evolves

## Architecture

```
trainer/
├── cmd/totem-trainer/main.go
├── internal/
│   ├── keymap/       # ZMK keymap parser
│   ├── lesson/       # Lesson engine + adaptive algorithm
│   ├── stats/        # Per-key stats, history, persistence
│   └── ui/           # Bubble Tea models + views
├── go.mod
├── go.sum
└── data/
    └── wordlists/    # Built-in word lists
```

Dependencies: `bubbletea`, `lipgloss`, `bubbles`
Data: `~/.config/totem-trainer/` (history.json, progress.json)

## ZMK Keymap Parser

Minimal parser — extracts only what we need from the `.keymap` file:

- Layer names and bindings from the `keymap` node
- Behavior references: `&kp`, `&mt`, `&lt`, `&trans`, `&none`
- Combo definitions

Skips hardware config, includes, and non-keymap nodes. Simple state-machine lexer, not a full devicetree parser.

```go
type Keymap struct {
    Layers []Layer
    Combos []Combo
}
type Layer struct {
    Name     string
    Bindings []Key // 38 entries by position
}
type Key struct {
    Position int
    Type     string // "kp", "mt", "lt", "trans", "none"
    Tap      string // character on tap
    Hold     string // modifier or layer on hold
}
```

## Lesson Engine

### Stages

| Stage | Keys | Source |
|-------|------|--------|
| Home row | a r s t d h n e i o | BASE row 1 |
| Top row | q w f p g j l u y ; | BASE row 0 |
| Bottom row | z x c v b k m , . / + pinky | BASE row 2 |
| Full alpha | All above | BASE layer |
| Numbers | 0-9 via NAV | NAV layer |
| Symbols | ! @ # $ % ^ & * ~ ` _ \| | SYM layer |
| Mixed | Real terminal commands | All layers |

### Adaptive Algorithm

- Track per-key: accuracy %, avg latency, last N attempts
- Weighted pool: low-accuracy/high-latency keys get more representation
- New keys unlock into pool when existing keys hit >95% accuracy
- Start each stage with 3-4 keys, grow organically

### Exercise Generation

- Filter English wordlist to currently-unlocked keys
- Terminal snippets for symbol/number stages
- Bigram drills for weak finger transitions

## Typing Flow

1. Launch → Main menu (Practice / Stats / Quit)
2. Practice → Lesson picker with stage unlock status
3. Select stage → Exercise generated, biased toward weak keys
4. Type → Real-time green/red feedback, live WPM/accuracy
5. Complete → Results with per-key breakdown
6. Loop → Back to picker or quit

## Stats & Persistence

Session record appended to `~/.config/totem-trainer/history.json`:

```json
{
  "timestamp": "2026-03-08T14:30:00Z",
  "stage": "home_row",
  "duration_secs": 120,
  "wpm": 34.5,
  "accuracy": 0.92,
  "total_keys": 248,
  "per_key": { "a": {"hits": 30, "misses": 2, "avg_ms": 180} }
}
```

TUI stats views:
- WPM/accuracy sparkline trend
- Totem keyboard heatmap (physical layout, colored by accuracy)
- Stage progress overview

## Stack

- Go + Bubble Tea + lipgloss + bubbles
- No external databases — JSON files only
- Reads `config/totem.keymap` as source of truth

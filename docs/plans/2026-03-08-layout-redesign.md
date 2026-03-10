0# Totem Keyboard Layout Redesign

## Goal

Redesign the Totem layout to optimize for a **terminal-first workflow**: tmux, shell,
and Claude Code. Must not sacrifice the ability to code in the future.

## Approach: "Terminal-First Rethink"

- Colemak DHm alphas are untouchable
- Right pinky key dedicated to tmux prefix
- TAB freed from dual-role (plain TAB for shell autocomplete)
- Combos for high-frequency terminal symbols on base layer
- Directional keys use n/e/i/o (home row, zero finger movement) across all layers
- Separate NUM and NAV layers (not combined)
- SYM layer reorganized around actual usage
- MOUSE layer toggled via both inner thumbs
- Future VIM layer ready when needed (remap vim to match n/e/i/o directions)

## Dependencies

### zmk-helpers (urob)

Add to `config/west.yml`:

```yaml
manifest:
  remotes:
    - name: urob
      url-base: https://github.com/urob
  projects:
    - name: zmk-helpers
      remote: urob
      revision: v0.3
```

Provides:
- **`ZMK_LAYER`** — simplified layer definitions
- **`ZMK_COMBO`** — combos using readable key labels instead of numbers
- **`ZMK_MACRO`** — multi-key sequence macros
- **`ZMK_HOLD_TAP`** — custom hold-tap behaviors
- **`ZMK_MOD_MORPH`** — modifier-dependent key behavior
- **`ZMK_UNICODE_SINGLE`** / **`ZMK_UNICODE_PAIR`** — direct Unicode character output
- **Key labels for Totem** (`totem.h`) — human-readable position names

---

## Totem Physical Layout Reference

38-key split columnar keyboard:
- 5 keys per hand on top two rows
- 6 keys per hand on bottom row (1 extra pinky key per side)
- 3 thumb keys per hand
- Pinky keys sit at home row height, offset outward

### Key Position Numbers

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │  0  │  1  │  2  │  3  │  4  │   │  5  │  6  │  7  │  8  │  9  │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ 20  │ 10  │ 11  │ 12  │ 13  │ 14  │   │ 15  │ 16  │ 17  │ 18  │ 19  │ 31  │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ 21  │ 22  │ 23  │ 24  │ 25  │   │ 26  │ 27  │ 28  │ 29  │ 30  │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ 32  │ 33  │ 34  │   │ 35  │ 36  │ 37  │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

### Key Labels (from zmk-helpers `totem.h`)

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ LT4 │ LT3 │ LT2 │ LT1 │ LT0 │   │ RT0 │ RT1 │ RT2 │ RT3 │ RT4 │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ LB5 │ LM4 │ LM3 │ LM2 │ LM1 │ LM0 │   │ RM0 │ RM1 │ RM2 │ RM3 │ RM4 │ RB5 │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ LB4 │ LB3 │ LB2 │ LB1 │ LB0 │   │ RB0 │ RB1 │ RB2 │ RB3 │ RB4 │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ LH2 │ LH1 │ LH0 │   │ RH0 │ RH1 │ RH2 │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Naming convention:** L/R = left/right, T/M/B/H = top/middle/bottom/thumb, 0 = innermost → 5 = outermost.

**Group labels:** `KEYS_L`, `KEYS_R`, `THUMBS_L`, `THUMBS_R`, `THUMBS`

---

## Proposed Layout

### Layer Access Summary

**Single-thumb hold layers:**

| Hold | Layer | Content hand |
|---|---|---|
| Left middle thumb `LH1` (TAB/NAV) | NAV | Right hand |
| Left inner thumb `LH2` (DEL/NUM) | NUM | Right hand |
| Right middle thumb `RH1` (BSPC/SYM) | SYM | Both hands |

**Two-thumb hold layers** (one thumb holds two adjacent keys, opposite hand is free):

| Hold | Layer | Content hand |
|---|---|---|
| Left `LH2` + `LH1` (DEL + TAB) | FN | Right hand |
| Left `LH1` + `LH0` (TAB + SPACE) | MEDIA | Right hand |
| Right `RH0` + `RH1` (ENTER + BSPC) | SYSTEM | Left hand |

**Toggle layers** (tap combo to enter/exit, both hands free):

| Tap combo | Layer |
|---|---|
| `LH2` + `RH0` (both inner thumbs) | MOUSE |
**Safety layer:**

| Hold combo | Layer |
|---|---|
| `LB5` + `RB5` (both pinky keys) | DANGER (reset/bootloader) |

---

### BASE (Layer 0) — Colemak DHm + Home-Row Mods

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │  q  │  w  │  f  │  p  │  b  │   │  j  │  l  │  u  │  y  │  ;  │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ESC │a/GUI│r/ALT│s/CTL│t/SFT│  g  │   │  m  │n/SFT│e/CTL│i/ALT│o/GUI│TMUX │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │  z  │  x  │  c  │  d  │  v  │   │  k  │  h  │  ,  │  .  │  /  │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ DEL/│ TAB/│SPACE│   │ENTER│BSPC/│ DEL │
                  │ NUM │ NAV │     │   │     │ SYM │     │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Home-row mods:** tap-preferred, 170ms tapping term, 100ms quick-tap.

**Key changes from original:**
- DEL/NUM `LH2` — hold for NUM layer, tap for DEL
- TAB/NAV `LH1` — hold for NAV layer, tap for TAB (quick tap = autocomplete, hold = arrows)
- BSPC/SYM `RH1` — hold for SYM layer, tap for backspace (most-used right thumb)
- Right outer thumb `RH2` → DEL (available on both sides)
- Right pinky `RB5` → TMUX prefix (dedicated key, sends Ctrl+B or custom prefix)
- Left pinky `LB5` → ESC (unchanged, good position)

**Base layer combos** (simultaneous press):

| Combo | Key Labels | Output | Reason |
|-------|------------|--------|--------|
| Tilde | `LT1` + `LT0` (f + p) | `~` | Home dir paths, used constantly |
| Pipe | `RT1` + `RT2` (l + u) | `\|` | Command chaining |
| Backtick | `LT3` + `LT2` (w + f) | `` ` `` | Shell command substitution |
| Underscore | `LB3` + `LB2` (x + c) | `_` | File names, variables |
| Mouse toggle | `LH2` + `RH0` | Toggle MOUSE layer | Both inner thumbs |

---

### NAV (Layer 1) — Hold left middle thumb `LH1` (TAB/NAV)

Arrows and navigation on right hand. Left hand has mods. Left thumb holds TAB, right hand is completely free.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │HOME │PG_DN│PG_UP│ END │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ GUI │ ALT │ CTL │SHIFT│ ___ │   │ ___ │  ←  │  ↓  │  ↑  │  →  │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │█████│ ___ │   │ ___ │ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Arrows on `RM1`/`RM2`/`RM3`/`RM4` (n/e/i/o)** — all fingers stay on home row, zero movement
- **HOME/END/PG_UP/PG_DN on top row** — same columns as arrows
- **Left home row = mods** — allows NAV + modifier combos (e.g., Shift+Arrow to select text)
- **Left thumb hold, right hand free** — no same-hand contortion
- **TAB quick-tap safe** — 170ms tapping term + 100ms quick-tap means sharp taps always register as TAB

---

### NUM (Layer 2) — Hold left inner thumb `LH2` (DEL/NUM)

Numpad on right hand. Left hand has mods.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │  `  │  7  │  8  │  9  │  ~  │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ GUI │ ALT │ CTL │SHIFT│ ___ │   │  -  │  4  │  5  │  6  │  0  │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │  /  │  1  │  2  │  3  │  .  │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │█████│ ___ │ ___ │   │  =  │  +  │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Numpad on right hand** — 1-9 in phone-style grid, 0 on home row pinky
- **`~` and `` ` `` on top row** — also available as base combos, but here for convenience
- **`-`, `/`, `=`, `+` around the numpad** — arithmetic operators where you'd expect them
- **Left home row = mods** — allows NUM + modifier combos
- **Tmux window switching flow:** right pinky (TMUX) → left thumb hold (NUM) → right hand number

---

### SYM (Layer 3) — Hold right inner thumb `RH0` (BSPC/SYM)

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │  @  │  #  │  $  │  %  │  ^  │   │  &  │  *  │  +  │  <  │  >  │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│  £  │  ~  │  |  │  `  │  '  │  "  │   │  !  │  {  │  }  │  -  │  =  │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │  \  │  (  │  )  │  [  │  ]  │   │  _  │  ;  │  :  │  ?  │ ___ │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ___ │   │█████│ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **`~ | ` ' "` on left home row** — most-used terminal symbols at fingertips
- **`! { } - =` on right home row** — bang, braces, operators (replaces wasted mods)
- **`£` on left pinky extra** — easy access without taking a prime position
- **Brackets grouped on left bottom** — `()` middle, `[]` inner, `\` pinky
- **`_ ; : ?` on right bottom** — useful characters accessible but not prime
- **`@ # $ % ^` on top left, `& * + < >` on top right** — less frequent symbols

---

### MOUSE (Layer 4) — Toggle via both inner thumbs `LH2` + `RH0` (DEL + ENTER)

Mouse movement on right home row (n/e/i/o), scroll on right bottom row, clicks on left home row. Both hands free after toggle.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ ___ │RCLK │MCLK │LCLK │ ___ │   │ ___ │ MV← │ MV↓ │ MV↑ │ MV→ │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │SC_L │SC_DN│SC_UP│SC_R │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Movement on `RM1`/`RM2`/`RM3`/`RM4` (n/e/i/o)** — directional consistency with NAV arrows
- **Scroll on `RB1`/`RB2`/`RB3`/`RB4` (h/,/.//)** — same columns as movement, one row down
- **Clicks on left home row** — LCLK on index (t), MCLK on middle (s), RCLK on ring (r)
- **Toggle, not hold** — both hands free for mouse control
- **Exit** — tap both inner thumbs again to return to BASE

---

### FN (Layer 5) — Hold left thumbs `LH2` + `LH1` (DEL + TAB)

F-keys on right hand, mods on left. Same grid as NUM for muscle memory.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ F7  │ F8  │ F9  │ F12 │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ GUI │ ALT │ CTL │SHIFT│ ___ │   │ ___ │ F4  │ F5  │ F6  │ F11 │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ F1  │ F2  │ F3  │ F10 │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │█████│█████│ ___ │   │ ___ │ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **F-keys on right hand** — same grid positions as NUM (muscle memory: F1=1, F4=4, F7=7)
- **Left home row = mods** — allows Ctrl+F5 etc. with one hand on mods, other on F-keys
- **Left thumb holds both `LH2` + `LH1`** — one thumb, two adjacent keys, right hand fully free

---

### MEDIA (Layer 5) — Hold left thumbs `LH1` + `LH0` (TAB + SPACE)

Media and brightness on right hand.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │BRI- │BRI+ │ ___ │ ___ │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ ___ │ ___ │ ___ │ ___ │ ___ │   │MUTE │VOL- │VOL+ │PREV │NEXT │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │█████│█████│   │PLAY │ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Media on right home row** — MUTE/VOL-/VOL+/PREV/NEXT are the most used, at fingertips
- **Brightness on right top row** — less frequent, easy reach
- **PLAY on right inner thumb `RH0`** — quick toggle play/pause
- **Left thumb holds both `LH1` + `LH0`** — one thumb, two adjacent keys, right hand fully free

---

### SYSTEM (Layer 6) — Hold right thumbs `RH0` + `RH1` (ENTER + BSPC)

Bluetooth and output toggle on left hand.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │O_TOG│BTNXT│BTPRV│ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │BTCLR│BT 0 │BT 1 │BT 2 │   │ ___ │ ___ │ ___ │ ___ │ ___ │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ___ │   │█████│█████│ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **BT/system on left hand** — right thumb holds `RH0` + `RH1`, left hand is fully free
- **O_TOG** — toggle between BLE and USB output
- **BT_CLR on bottom row** — harder to hit accidentally than home row
- **BT 0/1/2** — profile selection for connecting to different devices

---

### DANGER (Layer 7) — Hold both pinky keys `LB5` + `RB5` (ESC + TMUX)

Hardest combo to trigger accidentally. Only reset/bootloader.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│█████│ ___ │ ___ │RESET│ ___ │ ___ │   │ ___ │ ___ │RESET│ ___ │ ___ │█████│
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │BOOT │ ___ │ ___ │   │ ___ │ ___ │BOOT │ ___ │ ___ │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Both pinky keys held** — deliberate, two-hand gesture, nearly impossible to trigger accidentally
- **RESET/BOOT mirrored on both halves** — middle finger positions (`LM2`/`RM2` and `LB2`/`RB2`), flash either half
- **Everything else blank** — no accidental keypresses while in this layer

---

## Directional Consistency

All directional controls use `RM1`/`RM2`/`RM3`/`RM4` (n/e/i/o, right hand home row, zero finger movement):

| Layer | RM1 (n) | RM2 (e) | RM3 (i) | RM4 (o) |
|-------|---------|---------|---------|---------|
| NAV | ← | ↓ | ↑ | → |
| MOUSE | MV← | MV↓ | MV↑ | MV→ |
| Future VIM | (remap in neovim config to match) |

Scroll on MOUSE uses `RB1`/`RB2`/`RB3`/`RB4` (bottom row, same columns):

| RB1 (h) | RB2 (,) | RB3 (.) | RB4 (/) |
|---------|---------|---------|---------|
| SC_L | SC_DN | SC_UP | SC_R |

---

## Two-Thumb Layer Access

Two adjacent thumb keys on one hand activate a layer on the opposite hand. One thumb can press both adjacent keys.

| Combo | Hold hand | Content hand | Layer |
|---|---|---|---|
| `LH2` + `LH1` (DEL + TAB) | Left | Right | FN |
| `LH1` + `LH0` (TAB + SPACE) | Left | Right | MEDIA |
| `RH0` + `RH1` (ENTER + BSPC) | Right | Left | SYSTEM |
| `LB5` + `RB5` (ESC + TMUX) | Both pinkies | Both | DANGER |

---

## Future Considerations

### VIM Layer
When neovim practice ramps up:
- Remap vim's hjkl to match n/e/i/o directional convention
- Potentially add a dedicated VIM layer with motion keys (w/b/e word movement)
- Toggle on/off rather than hold (for sustained vim sessions)

### Multiple Profiles
ZMK doesn't support runtime profile switching, but we can:
- Build multiple firmware files with different layouts
- Use conditional layers or toggles for mode switching
- Create a "coding mode" that swaps SYM layer priorities

---

## Summary of Changes from Original

| Change | Why |
|--------|-----|
| Colemak DHm (corrected) | Proper matrix variant for columnar board |
| Right pinky `RB5` → TMUX prefix | Eliminates Ctrl+B two-key combo |
| TAB freed from dual-role | Fixes shell autocomplete misfires |
| DEL/NUM `LH2` on left inner thumb | Numpad on right hand, left hand free |
| TAB/NAV `LH1` on left middle thumb | Hold for arrows on right hand, tap for TAB |
| BSPC/SYM `RH1` on right middle thumb | Backspace on strongest thumb position |
| Separate NUM and NAV layers | No same-hand thumb+finger contortion |
| Arrows on n/e/i/o (`RM1-4`) | Home row, zero finger movement |
| MOUSE toggle via `LH2` + `RH0` | Both hands free for mouse control |
| ADJ split into FN/MEDIA/SYSTEM/DANGER | Focused layers, ergonomic two-thumb access |
| SYM right bottom row → operators | Media moved to dedicated MEDIA layer |
| Base layer combos for ~ \| ` _ | High-frequency symbols without layer switch |
| SYM reorganized | Terminal symbols promoted, coding symbols retained |
| Q+W → ESC combo removed | Redundant with ESC pinky key |
| zmk-helpers dependency added | Clean macros, combos, key labels, unicode support |

---

## Design Status

**Status: DESIGN IN PROGRESS**

Resolved:
- [x] Base layer (Colemak DHm + home-row mods + combos)
- [x] NAV layer (arrows on n/e/i/o, mods on left)
- [x] NUM layer (numpad on right, mods on left)
- [x] SYM layer (terminal symbols promoted, brackets grouped)
- [x] MOUSE layer (toggle, movement on n/e/i/o)
- [x] Directional consistency across layers
- [x] zmk-helpers integration

- [x] ADJ split into FN, MEDIA, SYSTEM, DANGER layers
- [x] Two-thumb layer access design
Open:
- [ ] VIM layer design (deferred until neovim usage increases)

# Totem Keyboard Layout Redesign

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

---

## Totem Physical Layout Reference

38-key split columnar keyboard:
- 5 keys per hand on top two rows
- 6 keys per hand on bottom row (1 extra pinky key per side)
- 3 thumb keys per hand
- Pinky keys sit at home row height, offset outward

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

---

## Proposed Layout

### Layer Access Summary

| Hold / Action | Layer |
|---|---|
| Left inner thumb (DEL/NUM) | NUM |
| Right outer thumb (NAV) | NAV |
| Right middle thumb (BSPC/SYM) | SYM |
| NAV + SYM | ADJ |
| Tap both inner thumbs together | MOUSE (toggle on/off) |

---

### BASE (Layer 0) — Colemak DHm + Home-Row Mods

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │  q  │  w  │  f  │  p  │  b  │   │  j  │  l  │  u  │  y  │  ;  │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ESC │a/GUI│r/ALT│s/CTL│t/SFT│  g  │   │m/SFT│n/CTL│e/ALT│i/GUI│  o  │TMUX │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │  z  │  x  │  c  │  d  │  v  │   │  k  │  h  │  ,  │  .  │  /  │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ DEL/│ TAB │SPACE│   │ENTER│BSPC/│ NAV │
                  │ NUM │     │     │   │     │ SYM │     │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Home-row mods:** tap-preferred, 170ms tapping term, 100ms quick-tap.

**Key changes from original:**
- TAB is plain (no dual-role) — fixes shell autocomplete misfires
- DEL/NUM — hold for NUM layer, tap for DEL
- BSPC/SYM — hold for SYM layer, tap for backspace (most-used right thumb)
- NAV — right outer thumb, hold for NAV layer
- Right pinky → TMUX prefix (dedicated key, sends Ctrl+B or custom prefix)
- Left pinky → ESC (unchanged, good position)

**Base layer combos** (simultaneous press):

| Combo | Keys | Output | Reason |
|-------|------|--------|--------|
| Tilde | f + p | `~` | Home dir paths, used constantly |
| Pipe | l + u | `\|` | Command chaining |
| Backtick | w + f | `` ` `` | Shell command substitution |
| Underscore | x + c | `_` | File names, variables |
| Mouse toggle | DEL/NUM + ENTER | Toggle MOUSE layer | Both inner thumbs |

---

### NAV (Layer 1) — Hold right outer thumb

Arrows and navigation on right hand. Left hand free.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │HOME │PG_DN│PG_UP│ END │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ GUI │ ALT │ CTL │SHIFT│ ___ │   │ ___ │  ←  │  ↓  │  ↑  │  →  │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ___ │   │ ADJ │ ___ │█████│
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Arrows on n/e/i/o** — all fingers stay on home row, zero movement
- **HOME/END/PG_UP/PG_DN on top row** — same columns as arrows (n/e/i/o → l/u/y/;)
- **Left home row = mods** — allows NAV + modifier combos (e.g., Shift+Arrow to select text)
- **Left hand otherwise empty** — clean, no accidental presses
- **ADJ on right inner thumb** — NAV + ADJ to reach adjust layer

---

### NUM (Layer 2) — Hold left inner thumb (DEL/NUM)

Numpad on right hand. Left hand free.

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

### SYM (Layer 3) — Hold BSPC (right middle thumb)

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │  !  │  @  │  #  │  $  │  %  │   │  ^  │  &  │  *  │  '  │  "  │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │  ~  │  `  │  _  │  |  │  {  │   │  }  │SHIFT│ CTRL│ ALT │ GUI │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │  \  │  (  │  )  │  [  │  ]  │   │MUTE │VOL- │VOL+ │PREV │NEXT │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ADJ │   │PLAY │█████│ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **`~`, `` ` ``, `_`, `|` on left home row** — most-used terminal symbols, best finger positions
- **Brackets grouped logically on left** — `{}` on home inner, `()` on bottom middle, `[]` on bottom inner. All pairs adjacent.
- **Right home row = one-shot mods** — allows SYM + modifier combos without contortion
- **Media controls on right bottom** — less frequent, out of the way
- **`\` on left bottom pinky** — available but not prime real estate
- **`'` and `"` on right top** — easy access for shell quoting
- **ADJ on left inner thumb** — SYM + ADJ to reach adjust layer

---

### ADJ (Layer 4) — NAV + SYM together

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │RESET│BTCLR│O_TOG│ ___ │ ___ │   │ ___ │ F7  │ F8  │ F9  │ F12 │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │BOOT │BTNXT│ ___ │ ___ │ ___ │   │ ___ │ F4  │ F5  │ F6  │ F11 │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │BTPRV│BT 0 │BT 1 │BT 2 │   │ ___ │ F1  │ F2  │ F3  │ F10 │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ ___ │ ___ │ ___ │   │ ___ │█████│█████│
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **F-keys on right hand** — same grid positions as numbers on NUM (muscle memory transfer)
- **Bluetooth/system on left** — rarely used, safe from accidental presses
- **RESET/BOOT accessible** — for flashing firmware

---

### MOUSE (Layer 5) — Toggle: tap both inner thumbs

Both hands completely free after toggle.

```
      ┌─────┬─────┬─────┬─────┬─────┐   ┌─────┬─────┬─────┬─────┬─────┐
      │ ___ │ ___ │M_CLK│ ___ │ ___ │   │ ___ │ ___ │ ___ │ ___ │ ___ │
┌─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┐
│ ___ │ ___ │RCLK │LCLK │MCLK │ ___ │   │ ___ │ MV← │ MV↓ │ MV↑ │ MV→ │ ___ │
└─────┼─────┼─────┼─────┼─────┼─────┤   ├─────┼─────┼─────┼─────┼─────┼─────┘
      │ ___ │ ___ │ ___ │ ___ │ ___ │   │ ___ │SC_L │SC_DN│SC_UP│SC_R │
      └─────┴─────┴─────┴─────┴─────┘   └─────┴─────┴─────┴─────┴─────┘
                  ┌─────┬─────┬─────┐   ┌─────┬─────┬─────┐
                  │ TOG │ ___ │ ___ │   │LCLK │ ___ │ ___ │
                  └─────┴─────┴─────┘   └─────┴─────┴─────┘
```

**Design rationale:**
- **Mouse movement on n/e/i/o** — same home row positions as NAV arrows, zero finger movement. Same directional mapping used everywhere: n=left, e=down, i=up, o=right
- **Scroll on h/,/./slash** (bottom row, same columns) — SC_L/SC_DN/SC_UP/SC_R
- **Clicks on left home row** — LCLK on index (s position, most natural), RCLK on middle (r), MCLK on ring (t). M_CLK (middle click) also on top row for accessibility.
- **LCLK also on right inner thumb** — easy single-hand click while moving mouse with right fingers
- **TOG on left inner thumb** — same key that toggled on, tap again to exit

---

## Directional Consistency

All directional controls use n/e/i/o (right hand home row, zero finger movement):

| Layer | n | e | i | o |
|-------|---|---|---|---|
| NAV | ← | ↓ | ↑ | → |
| MOUSE | MV← | MV↓ | MV↑ | MV→ |
| Future VIM | (remap in neovim config to match) |

Scroll on MOUSE uses h/,/./slash (bottom row, same column positions):

| h | , | . | / |
|---|---|---|---|
| SC_L | SC_DN | SC_UP | SC_R |

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
| Right pinky → TMUX prefix | Eliminates Ctrl+B two-key combo |
| TAB freed from dual-role | Fixes shell autocomplete misfires |
| DEL/NUM on left inner thumb | Numpad on right hand, left hand free |
| NAV on right outer thumb | Arrows on right hand, left hand free |
| BSPC/SYM on right middle thumb | Backspace on strongest thumb position |
| Separate NUM and NAV layers | No same-hand thumb+finger contortion |
| Arrows on n/e/i/o | Home row, zero finger movement |
| MOUSE toggle via both inner thumbs | Both hands free for mouse control |
| Base layer combos for ~ \| ` _ | High-frequency symbols without layer switch |
| SYM reorganized | Terminal symbols promoted, coding symbols retained |
| Q+W → ESC combo removed | Redundant with ESC pinky key |

---

## Design Status

**Status: DESIGN COMPLETE — ready for implementation**

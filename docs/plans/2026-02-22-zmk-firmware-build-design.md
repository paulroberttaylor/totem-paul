# Totem ZMK Firmware Build Design

## Goal

Build custom ZMK firmware for the Totem split keyboard (Bluetooth, XIAO nRF52840) using GitHub Actions, with full control over keymap, advanced features, and reproducible builds.

## Approach

Fork `GEIGEIGEIST/zmk-config-totem`, strip to a clean Colemak base, customize.

## Repository Structure

```
totem-paul/
├── .github/workflows/
│   └── build.yml              # GitHub Actions — triggers on push/PR/manual
├── build.yaml                 # Build targets: left + right halves
├── config/
│   ├── west.yml               # ZMK module manifest (zmkfirmware/zmk@main)
│   ├── totem.conf             # Feature flags (mouse keys, BT, etc.)
│   ├── totem.keymap           # Custom keymap — clean Colemak with 4 layers
│   └── boards/shields/totem/  # Hardware definitions (from official repo)
│       ├── Kconfig.defconfig
│       ├── Kconfig.shield
│       ├── totem.dtsi
│       ├── totem.keymap
│       ├── totem.zmk.yml
│       ├── totem_left.conf
│       ├── totem_left.overlay
│       ├── totem_right.conf
│       └── totem_right.overlay
└── readme.md
```

## Build Pipeline

GitHub Actions workflow delegates to ZMK's official reusable workflow:

```yaml
on: [push, pull_request, workflow_dispatch]
jobs:
  build:
    uses: zmkfirmware/zmk/.github/workflows/build-user-config.yml@main
```

`build.yaml` specifies two build targets:

```yaml
include:
  - board: seeeduino_xiao_ble
    shield: totem_left
  - board: seeeduino_xiao_ble
    shield: totem_right
```

Output: Two `.uf2` files as GitHub Actions artifacts.

## Keymap Design

Clean Colemak DH with 4 layers:

| Layer | Purpose |
|-------|---------|
| BASE  | Colemak DH, home-row mods (GUI/Alt/Ctrl/Shift) |
| NAV   | Arrows, page up/down, brackets, numbers |
| SYM   | Symbols, media controls |
| ADJ   | F-keys, BT management, bootloader, reset |

Home-row mods: A=GUI, R=ALT, S=CTRL, T=SHIFT (left), mirrored on right.

Thumb cluster: DEL / TAB(hold=NAV) / SPACE (left), ENTER / ESC(hold=SYM) / BSPC (right).

## Configuration

```
CONFIG_ZMK_MOUSE=y
CONFIG_ZMK_BLE=y
CONFIG_BT_CTLR_TX_PWR_PLUS_8=y
CONFIG_ZMK_USB_LOGGING=n
```

## Flashing Workflow

1. Push keymap changes to GitHub
2. GitHub Actions builds firmware (~2-3 min)
3. Download `.uf2` artifacts from Actions run
4. Double-tap reset on each XIAO nRF52840 to enter bootloader
5. Drag `.uf2` onto USB mass storage device

## Decisions

- Shield definitions included in repo (no external module dependency)
- Mouse keys enabled from the start
- Starting from official fork for proven CI pipeline
- Colemak DH as base layout

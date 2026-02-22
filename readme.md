# totem-paul

Custom ZMK firmware for the [TOTEM](https://github.com/GEIGEIGEIST/TOTEM) split keyboard.

## Layout

Colemak DH with home-row mods and 4 layers: BASE, NAV, SYM, ADJ.

## Building

Push to GitHub. GitHub Actions builds firmware automatically.

Download `.uf2` files from the Actions tab > latest run > Artifacts.

## Flashing

1. Connect one half via USB
2. Double-tap the reset button to enter bootloader mode
3. Drag the appropriate `.uf2` file onto the USB mass storage device that appears
4. Repeat for the other half

## Features

- Mouse key emulation
- Bluetooth with increased TX power
- Home-row mods (tap-preferred, 170ms tapping term)
- ESC combo (Q+W)

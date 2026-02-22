# Totem ZMK Firmware Build — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Set up a GitHub Actions-powered ZMK firmware build for the Totem split keyboard with a clean Colemak DH keymap and mouse keys enabled.

**Architecture:** A standard zmk-config repository with shield definitions bundled inline. GitHub Actions uses ZMK's official reusable workflow to compile `.uf2` firmware for both keyboard halves on every push.

**Tech Stack:** ZMK firmware, Zephyr RTOS (via ZMK's build system), GitHub Actions, Devicetree overlays, Kconfig

---

**Note:** This is firmware configuration, not traditional software. There are no unit tests — "testing" means the GitHub Actions build compiles successfully and produces `.uf2` artifacts. Each task ends with a commit so progress is saved.

---

### Task 1: Create GitHub Actions Workflow

**Files:**
- Create: `.github/workflows/build.yml`

**Step 1: Create the workflow file**

```yaml
on: [push, pull_request, workflow_dispatch]

jobs:
  build:
    uses: zmkfirmware/zmk/.github/workflows/build-user-config.yml@main
```

This delegates entirely to ZMK's official reusable workflow. It reads `build.yaml` at the repo root to know what to compile.

**Step 2: Commit**

```bash
git add .github/workflows/build.yml
git commit -m "ci: add GitHub Actions ZMK build workflow"
```

---

### Task 2: Create Build Configuration

**Files:**
- Create: `build.yaml`

**Step 1: Create build.yaml**

```yaml
include:
  - board: seeeduino_xiao_ble
    shield: totem_left
  - board: seeeduino_xiao_ble
    shield: totem_right
```

This tells the build to produce firmware for both halves of the split keyboard, targeting the Seeed XIAO nRF52840 BLE board.

**Step 2: Commit**

```bash
git add build.yaml
git commit -m "build: add build targets for left and right halves"
```

---

### Task 3: Create West Manifest

**Files:**
- Create: `config/west.yml`

**Step 1: Create west.yml**

```yaml
manifest:
  remotes:
    - name: zmkfirmware
      url-base: https://github.com/zmkfirmware
  projects:
    - name: zmk
      remote: zmkfirmware
      revision: main
      import: app/west.yml
  self:
    path: config
```

This tells the Zephyr build system where to find ZMK firmware source.

**Step 2: Commit**

```bash
git add config/west.yml
git commit -m "build: add west manifest pointing to ZMK main"
```

---

### Task 4: Create Shield Hardware Definitions

These files define the Totem's physical key matrix, GPIO pin assignments, and board metadata. They are copied verbatim from the official GEIGEIGEIST/zmk-config-totem repository.

**Files:**
- Create: `config/boards/shields/totem/Kconfig.shield`
- Create: `config/boards/shields/totem/Kconfig.defconfig`
- Create: `config/boards/shields/totem/totem.dtsi`
- Create: `config/boards/shields/totem/totem.zmk.yml`
- Create: `config/boards/shields/totem/totem_left.overlay`
- Create: `config/boards/shields/totem/totem_right.overlay`
- Create: `config/boards/shields/totem/totem_left.conf`
- Create: `config/boards/shields/totem/totem_right.conf`

**Step 1: Create `config/boards/shields/totem/Kconfig.shield`**

```kconfig
# Copyright (c) 2022 The ZMK Contributors
# SPDX-License-Identifier: MIT

config SHIELD_TOTEM_LEFT
    def_bool $(shields_list_contains,totem_left)

config SHIELD_TOTEM_RIGHT
    def_bool $(shields_list_contains,totem_right)
```

**Step 2: Create `config/boards/shields/totem/Kconfig.defconfig`**

```kconfig
# Copyright (c) 2022 The ZMK Contributors
# SPDX-License-Identifier: MIT

if SHIELD_TOTEM_LEFT

config ZMK_KEYBOARD_NAME
    default "TOTEM"

config ZMK_SPLIT_ROLE_CENTRAL
    default y

endif

if SHIELD_TOTEM_LEFT || SHIELD_TOTEM_RIGHT

config ZMK_SPLIT
    default y

endif
```

**Step 3: Create `config/boards/shields/totem/totem.dtsi`**

```dts
/*
 * Copyright (c) 2022 The ZMK Contributors
 *
 * SPDX-License-Identifier: MIT
 */

#include <dt-bindings/zmk/matrix_transform.h>

/ {
    chosen {
        zmk,kscan = &kscan0;
        zmk,matrix_transform = &default_transform;
    };

    default_transform: keymap_transform_0 {
        compatible = "zmk,matrix-transform";
        columns = <10>;
        rows = <4>;
//             | SW01  | SW02  | SW03  | SW04  | SW05  |  | SW05  | SW04  | SW03  | SW02  | SW01  |
//             | SW06  | SW07  | SW08  | SW09  | SW10  |  | SW10  | SW09  | SW08  | SW07  | SW06  |
//      | SW16 | SW11  | SW12  | SW13  | SW14  | SW15  |  | SW15  | SW14  | SW13  | SW12  | SW11  | SW16  |
//                             | SW17  | SW18  | SW19  |  | SW19  | SW18  | SW17  |
        map = <
                RC(0,0) RC(0,1) RC(0,2) RC(0,3) RC(0,4)    RC(0,5) RC(0,6) RC(0,7) RC(0,8) RC(0,9)
                RC(1,0) RC(1,1) RC(1,2) RC(1,3) RC(1,4)    RC(1,5) RC(1,6) RC(1,7) RC(1,8) RC(1,9)
        RC(3,0) RC(2,0) RC(2,1) RC(2,2) RC(2,3) RC(2,4)    RC(2,5) RC(2,6) RC(2,7) RC(2,8) RC(2,9) RC(3,9)
                                RC(3,2) RC(3,3) RC(3,4)    RC(3,5) RC(3,6) RC(3,7)
        >;
    };


    kscan0: kscan_0 {
        compatible = "zmk,kscan-gpio-matrix";
        label = "KSCAN";

        diode-direction = "col2row";
        row-gpios
            = <&xiao_d 0 (GPIO_ACTIVE_HIGH | GPIO_PULL_DOWN)>
            , <&xiao_d 1 (GPIO_ACTIVE_HIGH | GPIO_PULL_DOWN)>
            , <&xiao_d 2 (GPIO_ACTIVE_HIGH | GPIO_PULL_DOWN)>
            , <&xiao_d 3 (GPIO_ACTIVE_HIGH | GPIO_PULL_DOWN)>
            ;
    };
};
```

**Step 4: Create `config/boards/shields/totem/totem.zmk.yml`**

```yaml
file_format: "1"
id: totem
name: TOTEM
type: shield
url: https://github.com/GEIGEIGEIST/TOTEM
requires: [seeeduino_xiao_ble]
features:
  - keys
siblings:
  - totem_left
  - totem_right
```

**Step 5: Create `config/boards/shields/totem/totem_left.overlay`**

```dts
/*
 * Copyright (c) 2022 The ZMK Contributors
 *
 * SPDX-License-Identifier: MIT
 */

#include "totem.dtsi"

&kscan0 {
    col-gpios
        = <&xiao_d 4 GPIO_ACTIVE_HIGH>
        , <&xiao_d 5 GPIO_ACTIVE_HIGH>
        , <&xiao_d 10 GPIO_ACTIVE_HIGH>
        , <&xiao_d 9 GPIO_ACTIVE_HIGH>
        , <&xiao_d 8 GPIO_ACTIVE_HIGH>
        ;
};
```

**Step 6: Create `config/boards/shields/totem/totem_right.overlay`**

```dts
/*
 * Copyright (c) 2022 The ZMK Contributors
 *
 * SPDX-License-Identifier: MIT
 */

#include "totem.dtsi"

&default_transform {
    col-offset = <5>;
};

&kscan0 {
    col-gpios
        = <&xiao_d 8 GPIO_ACTIVE_HIGH>
        , <&xiao_d 9 GPIO_ACTIVE_HIGH>
        , <&xiao_d 10 GPIO_ACTIVE_HIGH>
        , <&xiao_d 5 GPIO_ACTIVE_HIGH>
        , <&xiao_d 4 GPIO_ACTIVE_HIGH>
        ;
};
```

**Step 7: Create `config/boards/shields/totem/totem_left.conf` and `totem_right.conf`**

Both files are empty (they exist as placeholders for per-half configuration if needed later). Create them as empty files.

**Step 8: Commit**

```bash
git add config/boards/shields/totem/
git commit -m "hw: add Totem shield definitions from official repo"
```

---

### Task 5: Create Firmware Configuration

**Files:**
- Create: `config/totem.conf`

**Step 1: Create totem.conf**

```kconfig
# Mouse key emulation
CONFIG_ZMK_MOUSE=y

# Bluetooth settings
CONFIG_ZMK_BLE=y
CONFIG_BT_CTLR_TX_PWR_PLUS_8=y

# Disable USB logging (saves space, not needed for normal use)
CONFIG_ZMK_USB_LOGGING=n
```

**Step 2: Commit**

```bash
git add config/totem.conf
git commit -m "config: enable mouse keys and BLE settings"
```

---

### Task 6: Create Clean Colemak DH Keymap

**Files:**
- Create: `config/totem.keymap`

This is the main file you'll edit to customize your keyboard. It overrides any default keymap in the shield directory.

**Step 1: Create the keymap**

The Totem has 38 keys in this physical layout:

```
         0   1   2   3   4       5   6   7   8   9
        10  11  12  13  14      15  16  17  18  19
  20    21  22  23  24  25      26  27  28  29  30   31
                32  33  34      35  36  37
```

Positions 20 and 31 are the extra pinky keys on the bottom row.

```dts
#include <behaviors.dtsi>
#include <dt-bindings/zmk/keys.h>
#include <dt-bindings/zmk/bt.h>
#include <dt-bindings/zmk/outputs.h>
#include <dt-bindings/zmk/mouse.h>

#define BASE 0
#define NAV  1
#define SYM  2
#define ADJ  3

// Home-row mod settings
&mt {
    quick-tap-ms = <100>;
    global-quick-tap;
    flavor = "tap-preferred";
    tapping-term-ms = <170>;
};

/ {
    combos {
        compatible = "zmk,combos";
        combo_esc {
            timeout-ms = <50>;
            key-positions = <0 1>;
            bindings = <&kp ESC>;
        };
    };

    keymap {
        compatible = "zmk,keymap";

        base_layer {
            label = "BASE";
// Colemak DH with home-row mods
//             Q       W       F       P       G           J       L       U       Y       ;
//             A/GUI   R/ALT   S/CTRL  T/SHFT  D           H       N/SHFT  E/CTRL  I/ALT   O/GUI
//   ESC       Z       X       C       V       B           K       M       ,       .       /       \
//                             DEL     TAB/NAV SPACE       ENTER   ESC/SYM BSPC
            bindings = <
                &kp Q       &kp W       &kp F       &kp P       &kp G           &kp J       &kp L       &kp U       &kp Y       &kp SEMI
                &mt LGUI A  &mt LALT R  &mt LCTRL S &mt LSHFT T &kp D           &kp H       &mt RSHFT N &mt RCTRL E &mt RALT I  &mt RGUI O
    &kp ESC     &kp Z       &kp X       &kp C       &kp V       &kp B           &kp K       &kp M       &kp COMMA   &kp DOT     &kp FSLH    &kp BSLH
                                        &kp DEL   &lt NAV TAB   &kp SPACE       &kp RET   &lt SYM ESC   &kp BSPC
            >;
        };

        nav_layer {
            label = "NAV";
// Navigation + numbers
//             ESC     HOME    UP      END     {           }       7       8       9       +
//             SHIFT   LEFT    DOWN    RIGHT   [           ]       4       5       6       -
//   ___       DEL     PG_UP   CAPS    PG_DN   (           )       1       2       3       =       ___
//                             ___     ___     ___         ADJ     0       .
            bindings = <
                &kp ESC     &kp HOME    &kp UP      &kp END     &kp LBRC        &kp RBRC    &kp N7      &kp N8      &kp N9      &kp PLUS
                &kp LSHFT   &kp LEFT    &kp DOWN    &kp RIGHT   &kp LBKT        &kp RBKT    &kp N4      &kp N5      &kp N6      &kp MINUS
    &trans      &kp DEL     &kp PG_UP   &kp CAPS    &kp PG_DN   &kp LPAR        &kp RPAR    &kp N1      &kp N2      &kp N3      &kp EQUAL   &trans
                                        &trans      &trans      &trans          &mo ADJ     &kp N0      &kp DOT
            >;
        };

        sym_layer {
            label = "SYM";
// Symbols + media
//             !       @       #       $       %           ^       &       *       '       "
//             ~       `       _       |       TAB         MUTE    SHFT    CTRL    ALT     GUI
//   ___       ___     ___     ___     ___     ___         VOL-    VOL+    PREV    NEXT    \       ___
//                             ___     ___     ADJ         ___     PLAY    ___
            bindings = <
                &kp EXCL    &kp AT      &kp HASH    &kp DLLR    &kp PRCNT       &kp CARET   &kp AMPS    &kp ASTRK   &kp SQT     &kp DQT
                &kp TILDE   &kp GRAVE   &kp UNDER   &kp PIPE    &kp TAB         &kp C_MUTE  &kp RSHFT   &kp RCTRL   &kp RALT    &kp RGUI
    &trans      &trans      &trans      &trans      &trans      &trans          &kp C_VOL_DN &kp C_VOL_UP &kp C_PREV &kp C_NEXT  &kp BSLH    &trans
                                        &trans      &trans      &mo ADJ         &trans      &kp C_PP    &trans
            >;
        };

        adjust_layer {
            label = "ADJ";
// F-keys + Bluetooth + system
//             RESET   BT_CLR  OUT_TOG ___     ___         ___     F7      F8      F9      F12
//             BOOT    BT_NXT  ___     ___     ___         ___     F4      F5      F6      F11
//   ___       ___     BT_PRV  BT 0    BT 1    BT 2        ___     F1      F2      F3      F10     ___
//                             ___     ___     ___         ___     ___     ___
            bindings = <
                &sys_reset  &bt BT_CLR &out OUT_TOG &trans      &trans          &trans      &kp F7      &kp F8      &kp F9      &kp F12
                &bootloader &bt BT_NXT  &trans      &trans      &trans          &trans      &kp F4      &kp F5      &kp F6      &kp F11
    &trans      &trans     &bt BT_PRV &bt BT_SEL 0 &bt BT_SEL 1 &bt BT_SEL 2   &trans      &kp F1      &kp F2      &kp F3      &kp F10     &trans
                                        &trans      &trans      &trans          &trans      &trans      &trans
            >;
        };
    };
};
```

**Step 2: Commit**

```bash
git add config/totem.keymap
git commit -m "keymap: add clean Colemak DH with nav/sym/adj layers"
```

---

### Task 7: Create README

**Files:**
- Create: `readme.md`

**Step 1: Create readme.md**

```markdown
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
```

**Step 2: Commit**

```bash
git add readme.md
git commit -m "docs: add README with build and flash instructions"
```

---

### Task 8: Create GitHub Repository and Push

**Step 1: Create the GitHub repo**

```bash
gh repo create totem-paul --public --source=. --remote=origin --description "Custom ZMK firmware for TOTEM split keyboard"
```

**Step 2: Push to trigger the first build**

```bash
git push -u origin main
```

**Step 3: Verify the build started**

```bash
gh run list --limit 1
```

Expected: A workflow run should appear with status "in_progress" or "queued".

**Step 4: Wait for build and verify success**

```bash
gh run watch
```

Expected: Build completes with green checkmark. Two `.uf2` artifacts should be available.

**Step 5: Verify artifacts exist**

```bash
gh run download --dir /tmp/totem-firmware
ls /tmp/totem-firmware/
```

Expected: Two directories containing `.uf2` files — one for `totem_left` and one for `totem_right`.

---

### Task 9: Flash and Verify (Manual — User Action)

This task is manual and cannot be automated.

1. Connect the **left** half via USB
2. Double-tap the reset button — a USB mass storage device appears
3. Drag `totem_left-seeeduino_xiao_ble-zmk.uf2` onto the device
4. The device auto-ejects and reboots
5. Repeat steps 1-4 for the **right** half with `totem_right-seeeduino_xiao_ble-zmk.uf2`
6. The keyboard should now be discoverable via Bluetooth as "TOTEM"

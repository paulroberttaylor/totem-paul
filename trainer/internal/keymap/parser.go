package keymap

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseFile reads a ZMK .keymap file and extracts layers and combos.
func ParseFile(path string) (*Keymap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keymap: %w", err)
	}
	return Parse(string(data))
}

// Parse parses ZMK keymap content.
func Parse(content string) (*Keymap, error) {
	km := &Keymap{}

	// Extract combos
	km.Combos = parseCombos(content)

	// Extract layers from the keymap node
	km.Layers = parseLayers(content)

	if len(km.Layers) == 0 {
		return nil, fmt.Errorf("no layers found in keymap")
	}

	return km, nil
}

var bindingsRe = regexp.MustCompile(`bindings\s*=\s*<([^>]+)>`)
var labelRe = regexp.MustCompile(`label\s*=\s*"([^"]+)"`)
var comboNameRe = regexp.MustCompile(`(\w+)\s*\{`)
var keyPosRe = regexp.MustCompile(`key-positions\s*=\s*<([^>]+)>`)

func parseLayers(content string) []Layer {
	// Find the keymap node
	kmStart := strings.Index(content, "keymap {")
	if kmStart == -1 {
		kmStart = strings.Index(content, "keymap{")
	}
	if kmStart == -1 {
		return nil
	}

	// Find the keymap block by brace matching
	kmBlock := extractBlock(content[kmStart:])

	// Split into layer blocks — find each bindings = <...>
	var layers []Layer
	idx := 0
	for idx < len(kmBlock) {
		// Find next "bindings = <"
		bindLoc := bindingsRe.FindStringIndex(kmBlock[idx:])
		if bindLoc == nil {
			break
		}

		// Look backwards from bindings for the label
		prefix := kmBlock[idx : idx+bindLoc[0]]
		labelMatch := labelRe.FindStringSubmatch(prefix)
		name := ""
		if labelMatch != nil {
			name = labelMatch[1]
		}

		// Extract bindings
		bindMatch := bindingsRe.FindStringSubmatch(kmBlock[idx:])
		if bindMatch == nil {
			break
		}

		keys := parseBindings(bindMatch[1])
		layers = append(layers, Layer{
			Name:     name,
			Index:    len(layers),
			Bindings: keys,
		})

		idx += bindLoc[1]
	}

	return layers
}

func parseBindings(raw string) []Key {
	// Split on & to get individual bindings
	// Raw looks like: "&kp Q  &kp W  &mt LGUI A ..."
	raw = strings.TrimSpace(raw)

	var keys []Key
	parts := strings.Split(raw, "&")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}

		k := Key{Position: len(keys)}
		behavior := fields[0]

		switch behavior {
		case "kp":
			k.Type = "kp"
			if len(fields) >= 2 {
				k.Tap = fields[1]
				k.Char = ZMKCodeToChar(fields[1])
			}
		case "mt":
			k.Type = "mt"
			if len(fields) >= 3 {
				k.Hold = fields[1]
				k.Tap = fields[2]
				k.Char = ZMKCodeToChar(fields[2])
			}
		case "lt":
			k.Type = "lt"
			if len(fields) >= 3 {
				k.Hold = fields[1]
				k.Tap = fields[2]
				// layer-tap keys are layer switches, not simple typeable chars
			}
		case "trans":
			k.Type = "trans"
		case "none":
			k.Type = "none"
		case "mo":
			k.Type = "mo"
			if len(fields) >= 2 {
				k.Hold = fields[1]
			}
		case "sys_reset":
			k.Type = "sys"
			k.Tap = "RESET"
		case "bootloader":
			k.Type = "sys"
			k.Tap = "BOOT"
		case "bt":
			k.Type = "bt"
			if len(fields) >= 2 {
				k.Tap = strings.Join(fields[1:], " ")
			}
		case "out":
			k.Type = "out"
			if len(fields) >= 2 {
				k.Tap = strings.Join(fields[1:], " ")
			}
		default:
			k.Type = behavior
			if len(fields) >= 2 {
				k.Tap = strings.Join(fields[1:], " ")
			}
		}

		keys = append(keys, k)
	}

	return keys
}

func parseCombos(content string) []Combo {
	// Find the combos node
	cStart := strings.Index(content, "combos {")
	if cStart == -1 {
		cStart = strings.Index(content, "combos{")
	}
	if cStart == -1 {
		return nil
	}

	cBlock := extractBlock(content[cStart:])

	var combos []Combo

	// Find each combo sub-block (skip "compatible" property)
	lines := strings.Split(cBlock, "\n")
	var currentName string
	var currentPositions []int
	var currentBinding string
	inCombo := false
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "compatible") {
			continue
		}

		// Detect combo block start: "name {"
		if match := comboNameRe.FindStringSubmatch(trimmed); match != nil && !strings.Contains(trimmed, "combos") {
			if depth == 1 { // inside combos node
				currentName = match[1]
				inCombo = true
			}
		}

		for _, ch := range trimmed {
			if ch == '{' {
				depth++
			} else if ch == '}' {
				depth--
				if depth == 1 && inCombo {
					// End of combo block
					combos = append(combos, Combo{
						Name:      currentName,
						Positions: currentPositions,
						Output:    currentBinding,
					})
					currentName = ""
					currentPositions = nil
					currentBinding = ""
					inCombo = false
				}
			}
		}

		if inCombo {
			// Parse key-positions
			if posMatch := keyPosRe.FindStringSubmatch(trimmed); posMatch != nil {
				for _, s := range strings.Fields(posMatch[1]) {
					if n, err := strconv.Atoi(s); err == nil {
						currentPositions = append(currentPositions, n)
					}
				}
			}
			// Parse bindings
			if bindMatch := bindingsRe.FindStringSubmatch(trimmed); bindMatch != nil {
				// Extract the character from the binding
				parts := strings.Fields(strings.TrimPrefix(strings.TrimSpace(bindMatch[1]), "&"))
				if len(parts) >= 2 && parts[0] == "kp" {
					currentBinding = ZMKCodeToChar(parts[1])
				}
			}
		}
	}

	return combos
}

// extractBlock extracts content between the first { and its matching }.
func extractBlock(s string) string {
	start := strings.Index(s, "{")
	if start == -1 {
		return ""
	}

	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

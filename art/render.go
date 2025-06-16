package art

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

type Font map[rune][]string

const (
	asciiSpace  = 32
	asciiEnd    = 127
	bannerChars = 95
)

var bannerHashes = map[string]string{
	"standard.txt":   "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf",
	"shadow.txt":     "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73",
	"thinkertoy.txt": "64285e4960d199f4819323c4dc6319ba34f1f0dd9da14d07111345f5d76c3fa3",
}

func LoadBanner(filename string) (Font, error) {
	switch filename {
	case "standard", "thinkertoy", "shadow":
		filename += ".txt"
	}
	height, err := GetBannerHeight(filename)
	if err != nil {
		return nil, err
	}

	lines, err := ReadBannerFile(filename)
	if err != nil {
		return nil, err
	}

	if err := validateBannerHash(filename); err != nil {
		return nil, err
	}

	font, err := ParseBannerLines(lines, height)
	if err != nil {
		return nil, err
	}
	if err := ValidateFont(font); err != nil {
		return nil, err
	}
	return font, nil
}

func ReadBannerFile(filename string) ([]string, error) {
	file, err := os.Open("assets/" + filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open banner: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading banner: %w", err)
	}
	return lines, nil
}

func validateBannerHash(filename string) error {
	expected, ok := bannerHashes[filename]
	if !ok {
		return fmt.Errorf("no hash defined for %s", filename)
	}

	data, err := os.ReadFile("assets/" + filename)
	if err != nil {
		return fmt.Errorf("failed to read banner for hashing: %w", err)
	}

	hash := sha256.Sum256(data)
	actual := fmt.Sprintf("%x", hash[:])

	if actual != expected {
		return fmt.Errorf("invalid banner: hash mismatch for %s, wanted: %s, got: %s", filename, expected, actual)
	}
	return nil
}

func ParseBannerLines(lines []string, height int) (Font, error) {
	font := make(Font)
	startChar := rune(asciiSpace)
	char := startChar

	for i := 0; i < len(lines); {
		if i < len(lines) && lines[i] == "" {
			i++
		} else {
			return nil, fmt.Errorf("missing empty line separator after character '%c'", char)
		}

		if i+height > len(lines) {
			break
		}

		block := lines[i : i+height]
		font[char] = block

		char++
		i += height
	}

	return font, nil
}

func ValidateFont(font Font) error {
	if len(font) != bannerChars {
		return fmt.Errorf("invalid number of characters parsed: expected %d, got %d", bannerChars, len(font))
	}

	return nil
}

func GetBannerHeight(filename string) (int, error) {
	switch filename {
	case "standard.txt", "shadow.txt", "thinkertoy.txt":
		return 8, nil
	default:
		return 0, fmt.Errorf("unknown banner file: %s", filename)
	}
}
func RenderText(input string, font Font) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input text is empty")
	}

	input = strings.ReplaceAll(input, "\\n", "\n")
	lines := strings.Split(input, "\n")
	var output strings.Builder

	for _, line := range lines {
		if line == "" {
			output.WriteString("\n")
			continue
		}
		for row := 0; row < 8; row++ {
			for _, ch := range line {
				if ch < asciiSpace || ch >= asciiEnd {
					return "", fmt.Errorf("unsupported character: %q", ch)
				}
				if glyph, ok := font[ch]; ok {
					output.WriteString(glyph[row])
				} else {
					output.WriteString("       ") // fallback for unknown char
				}
			}
			output.WriteString("\n")
		}
	}
	return output.String(), nil
}

// func PrintAsciiArt(s string, banner Font) {
// 	s = strings.ReplaceAll(s, "\\n", "\n")
// 	parts := strings.Split(s, "\n")
// 	for _, part := range parts {
// 		if part == "" {
// 			fmt.Println()
// 			continue
// 		}

// 		for _, ch := range part {
// 			if ch < asciiSpace || ch >= asciiEnd {
// 				fmt.Printf("Error: input contains non-ASCII character\n")
// 				return
// 			}
// 		}

// 		for row := 0; row < 8; row++ {
// 			line := ""
// 			for _, ch := range part {
// 				if ch == ' ' {
// 					line += "      "
// 				} else if art, ok := banner[ch]; ok && row < len(art) {
// 					line += art[row]
// 				}
// 			}
// 			fmt.Println(line)
// 		}
// 	}
// }

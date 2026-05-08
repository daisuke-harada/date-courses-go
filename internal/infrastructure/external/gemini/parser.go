package gemini

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// DateSpotCandidate は Gemini が返すスポット候補です。
type DateSpotCandidate struct {
	Name        string `json:"name"`
	CityName    string `json:"city_name"`
	Description string `json:"description"`
}

var jsonArrayRe = regexp.MustCompile(`(?s)\[.*\]`)

// ParseDateSpotCandidates は Gemini のレスポンステキストから DateSpotCandidate のスライスをパースします。
func ParseDateSpotCandidates(text string) ([]DateSpotCandidate, error) {
	text = strings.TrimSpace(text)

	// コードブロック（```json ... ```）を除去
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	// JSON 配列部分だけを抽出
	match := jsonArrayRe.FindString(text)
	if match == "" {
		return nil, fmt.Errorf("gemini parser: no JSON array found in response")
	}

	var candidates []DateSpotCandidate
	if err := json.Unmarshal([]byte(match), &candidates); err != nil {
		return nil, fmt.Errorf("gemini parser: unmarshal failed: %w", err)
	}

	// name・city_name が空の候補は除外
	valid := candidates[:0]
	for _, c := range candidates {
		if strings.TrimSpace(c.Name) != "" && strings.TrimSpace(c.CityName) != "" {
			valid = append(valid, c)
		}
	}

	return valid, nil
}

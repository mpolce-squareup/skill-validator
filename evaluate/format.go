package evaluate

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/dacharyc/skill-validator/judge"
	"github.com/dacharyc/skill-validator/util"
)

// Shorthand aliases for color constants to keep format strings compact.
const (
	ColorReset  = util.ColorReset
	ColorBold   = util.ColorBold
	ColorGreen  = util.ColorGreen
	ColorYellow = util.ColorYellow
	ColorCyan   = util.ColorCyan
	ColorRed    = util.ColorRed
)

// FormatResults formats a single EvalResult in the given format.
func FormatResults(w io.Writer, results []*EvalResult, format, display string) error {
	if len(results) == 0 {
		return nil
	}
	if len(results) == 1 {
		switch format {
		case "json":
			return PrintJSON(w, results)
		case "markdown":
			PrintMarkdown(w, results[0], display)
			return nil
		default:
			PrintText(w, results[0], display)
			return nil
		}
	}
	return FormatMultiResults(w, results, format, display)
}

// FormatMultiResults formats multiple EvalResults in the given format.
func FormatMultiResults(w io.Writer, results []*EvalResult, format, display string) error {
	switch format {
	case "json":
		return PrintJSON(w, results)
	case "markdown":
		PrintMultiMarkdown(w, results, display)
		return nil
	default:
		for i, r := range results {
			if i > 0 {
				_, _ = fmt.Fprintf(w, "\n%s\n", strings.Repeat("━", 60))
			}
			PrintText(w, r, display)
		}
		return nil
	}
}

// PrintText writes a human-readable text representation of an EvalResult.
func PrintText(w io.Writer, result *EvalResult, display string) {
	_, _ = fmt.Fprintf(w, "\n%sScoring skill: %s%s\n", ColorBold, result.SkillDir, ColorReset)

	if result.SkillScores != nil {
		_, _ = fmt.Fprintf(w, "\n%sSKILL.md Scores%s\n", ColorBold, ColorReset)
		printScoredText(w, result.SkillScores)
	}

	if display == "files" && len(result.RefResults) > 0 {
		for _, ref := range result.RefResults {
			_, _ = fmt.Fprintf(w, "\n%sReference: %s%s\n", ColorBold, ref.File, ColorReset)
			printScoredText(w, ref.Scores)
		}
	}

	if result.RefAggregate != nil {
		_, _ = fmt.Fprintf(w, "\n%sReference Scores (%d file%s)%s\n", ColorBold, len(result.RefResults), util.PluralS(len(result.RefResults)), ColorReset)
		printScoredText(w, result.RefAggregate)
	}

	_, _ = fmt.Fprintln(w)
}

// printScoredText writes all dimensions, overall, assessment, and novel details for a Scored value.
func printScoredText(w io.Writer, s judge.Scored) {
	for _, d := range s.DimensionScores() {
		printDimScore(w, d.Label, d.Value)
	}
	_, _ = fmt.Fprintf(w, "  %s\n", strings.Repeat("─", 30))
	_, _ = fmt.Fprintf(w, "  %sOverall:              %.2f/5%s\n", ColorBold, s.OverallScore(), ColorReset)

	if s.Assessment() != "" {
		_, _ = fmt.Fprintf(w, "\n  %s\"%s\"%s\n", ColorCyan, s.Assessment(), ColorReset)
	}
	if s.NovelDetails() != "" {
		_, _ = fmt.Fprintf(w, "  %sNovel details: %s%s\n", ColorCyan, s.NovelDetails(), ColorReset)
	}
}

func printDimScore(w io.Writer, name string, score int) {
	color := ColorGreen
	if score <= 2 {
		color = ColorRed
	} else if score <= 3 {
		color = ColorYellow
	}
	padding := max(22-len(name), 1)
	_, _ = fmt.Fprintf(w, "  %s:%s%s%d/5%s\n", name, strings.Repeat(" ", padding), color, score, ColorReset)
}

// --- JSON output ---

// EvalJSONOutput is the top-level JSON envelope.
type EvalJSONOutput struct {
	Skills []EvalJSONSkill `json:"skills"`
}

// EvalJSONSkill is one skill entry in JSON output.
type EvalJSONSkill struct {
	SkillDir     string             `json:"skill_dir"`
	SkillScores  *judge.SkillScores `json:"skill_scores,omitempty"`
	RefScores    []EvalJSONRef      `json:"reference_scores,omitempty"`
	RefAggregate *judge.RefScores   `json:"reference_aggregate,omitempty"`
}

// EvalJSONRef is one reference file entry in JSON output.
type EvalJSONRef struct {
	File   string           `json:"file"`
	Scores *judge.RefScores `json:"scores"`
}

// PrintJSON writes results as indented JSON.
func PrintJSON(w io.Writer, results []*EvalResult) error {
	out := EvalJSONOutput{
		Skills: make([]EvalJSONSkill, len(results)),
	}
	for i, r := range results {
		skill := EvalJSONSkill{
			SkillDir:     r.SkillDir,
			SkillScores:  r.SkillScores,
			RefAggregate: r.RefAggregate,
		}
		for _, ref := range r.RefResults {
			skill.RefScores = append(skill.RefScores, EvalJSONRef(ref))
		}
		out.Skills[i] = skill
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// --- Markdown output ---

// PrintMarkdown writes a single EvalResult as Markdown.
func PrintMarkdown(w io.Writer, result *EvalResult, display string) {
	_, _ = fmt.Fprintf(w, "## Scoring skill: %s\n", result.SkillDir)

	if result.SkillScores != nil {
		_, _ = fmt.Fprintf(w, "\n### SKILL.md Scores\n\n")
		printScoredMarkdown(w, result.SkillScores)
	}

	if display == "files" && len(result.RefResults) > 0 {
		for _, ref := range result.RefResults {
			_, _ = fmt.Fprintf(w, "\n### Reference: %s\n\n", ref.File)
			printScoredMarkdown(w, ref.Scores)
		}
	}

	if result.RefAggregate != nil {
		_, _ = fmt.Fprintf(w, "\n### Reference Scores (%d file%s)\n\n", len(result.RefResults), util.PluralS(len(result.RefResults)))
		printScoredMarkdown(w, result.RefAggregate)
	}
}

// PrintMultiMarkdown writes multiple EvalResults as Markdown, separated by rules.
func PrintMultiMarkdown(w io.Writer, results []*EvalResult, display string) {
	for i, r := range results {
		if i > 0 {
			_, _ = fmt.Fprintf(w, "\n---\n\n")
		}
		PrintMarkdown(w, r, display)
	}
}

// printScoredMarkdown writes a markdown table for all dimensions plus overall, assessment, and novel details.
func printScoredMarkdown(w io.Writer, s judge.Scored) {
	_, _ = fmt.Fprintf(w, "| Dimension | Score |\n")
	_, _ = fmt.Fprintf(w, "| --- | ---: |\n")
	for _, d := range s.DimensionScores() {
		_, _ = fmt.Fprintf(w, "| %s | %d/5 |\n", d.Label, d.Value)
	}
	_, _ = fmt.Fprintf(w, "| **Overall** | **%.2f/5** |\n", s.OverallScore())

	if s.Assessment() != "" {
		_, _ = fmt.Fprintf(w, "\n> %s\n", s.Assessment())
	}
	if s.NovelDetails() != "" {
		_, _ = fmt.Fprintf(w, "\n*Novel details: %s*\n", s.NovelDetails())
	}
}

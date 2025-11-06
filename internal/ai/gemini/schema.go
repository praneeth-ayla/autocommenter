package ai

var AnalyzeFilesForCommentsSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"files": map[string]any{
			"type":  "array",
			"items": map[string]any{"type": "string"},
		},
	},
	"required": []string{"files"},
}

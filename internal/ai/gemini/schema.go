package ai

// AnalyzeFilesForCommentsSchema defines the JSON schema for the AI's response when analyzing files for comments.
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

// GenerateCommentsForFilesSchema defines the JSON schema for the AI's response when generating comments for files.
var GenerateCommentsForFilesSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"files": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]any{
						"type": "string",
					},
					"content": map[string]any{
						"type": "string",
					},
				},
				"required": []string{"path", "content"},
			},
		},
	},
	"required": []string{"files"},
}

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

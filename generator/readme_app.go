package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"gopkg.in/yaml.v2"
)

type Tag struct {
	Dir  string `yaml:"dir"`
	Name string `yaml:"name"`
}

type Config struct {
	Tags []Tag `yaml:"tags"`
}

func main() {
	data, err := os.ReadFile("../tags.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.New("readme").Parse("# {{.Name}}\n"))

	for _, tag := range config.Tags {
		// 🔴 VULNERABILITY: Path traversal via tag.Dir
		tagPath := filepath.Join("../tags", tag.Dir)

		os.MkdirAll(tagPath, 0o755)

		var buf bytes.Buffer
		tmpl.Execute(&buf, tag)

		// This will write to ../tags/../../.github/workflows/ if manipulated
		os.WriteFile(filepath.Join(tagPath, "README.md"), buf.Bytes(), 0o644)
	}
}

package docs

import (
	"embed"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const BASE_DOCS_PATH string = "terragrunt"

//go:embed terragrunt/**/*.md
var documentation embed.FS

func getDocumentation(dirName string) (docs []string, err error) {
	docsPath := fmt.Sprintf("%s/%s", BASE_DOCS_PATH, dirName)
	dirEntries, err := documentation.ReadDir(docsPath)
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot read documentation on %s", docsPath)
		log.Error(errorMsg)
		return nil, err
	}

	for _, entry := range dirEntries {
		docFile := entry.Name()
		docContent := readEmbeddedFile(fmt.Sprintf("%s/%s", docsPath, docFile))

		docs = append(docs, docContent)
	}

	return docs, nil
}

func readEmbeddedFile(fileName string) (fileContent string) {
	filePath := fmt.Sprintf(fileName)
	fileBytes, _ := documentation.ReadFile(filePath)

	return string(fileBytes)
}

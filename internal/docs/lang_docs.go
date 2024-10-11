package docs

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"

	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
)

type LangDoc struct {
	Item       string
	InsertText string
	Content    string
}

func GetLocals(tree *sitter.Tree, content []byte) []string {
	localsPattern := `(
		(identifier) @constant
		(#match? @constant "^locals$")
	)`

	q, _ := sitter.NewQuery([]byte(localsPattern), hcl.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	var locals []string
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		m = qc.FilterPredicates(m, content)
		for _, c := range m.Captures {
			localsNode := c.Node.NextSibling().NextSibling()
			for i := 0; i < int(localsNode.ChildCount()); i++ {
				locals = append(locals, localsNode.Child(i).Child(0).Content(content))
			}
		}
	}

	return locals
}

func GetBlocks() []LangDoc {
	docs, err := getDocumentation("blocks")
	if err != nil {
		log.Panic("Failed to read blocks documentation")
		panic(1)
	}

	var langDocs []LangDoc

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		blockName := docLines[0]
		content := strings.Join(docLines[1:], "\n")

		// Some blocks are named, these ones will have a $name string
		// before the opened block in the InsertText field
		var insertText string
		switch blockName {
		case "dependency":
			insertText = fmt.Sprintf("%s \"$name\" {\n\tconfig_path = $0\n}", blockName)
		case "generate", "include":
			insertText = fmt.Sprintf("%s \"$name\" {\n\tpath = $0\n}", blockName)
		case "terraform":
			insertText = fmt.Sprintf("%s {\n\tsource = $0\n}", blockName)
		default:
			insertText = fmt.Sprintf("%s {\n\t$0\n}", blockName)
		}

		langDoc := LangDoc{
			Item:       blockName,
			InsertText: insertText,
			Content:    content,
		}

		langDocs = append(langDocs, langDoc)
	}

	return langDocs
}

func GetFunctions() []LangDoc {
	docs, err := getDocumentation("functions")
	if err != nil {
		log.Panic("Failed to read functions documentation")
		panic(1)
	}

	var langDocs []LangDoc

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		functionSignature := docLines[0]
		functionName := strings.Split(functionSignature, "(")[0]
		content := strings.Join(docLines[1:], "\n")

		langDoc := LangDoc{
			Item:       functionName,
			InsertText: functionSignature,
			Content:    content,
		}

		langDocs = append(langDocs, langDoc)
	}

	return langDocs
}

func GetTopLevelAttributes() []LangDoc {
	docs, err := getDocumentation("attributes-top-level")
	if err != nil {
		log.Panic("Failed to read arguments documentation")
		panic(1)
	}

	var langDocs []LangDoc

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		attrName := docLines[0]
		content := strings.Join(docLines[1:], "\n")

		// A couple of attributes have a different type, so we build the
		// InsertText field based on attrName
		var insertText string
		switch attrName {
		case "inputs":
			insertText = fmt.Sprintf("%s = {\n\t$0\n}", attrName)
		default:
			insertText = fmt.Sprintf("%s = \"$0\"", attrName)
		}

		langDoc := LangDoc{
			Item:       attrName,
			InsertText: insertText,
			Content:    content,
		}

		langDocs = append(langDocs, langDoc)
	}

	return langDocs
}

func GetAttributes() []LangDoc {
	docs, err := getDocumentation("attributes")
	if err != nil {
		log.Panic("Failed to read arguments documentation")
		panic(1)
	}

	var langDocs []LangDoc

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		attrName := docLines[0]
		content := strings.Join(docLines[1:], "\n")

		// A couple of attributes have a different type, so we build the
		// InsertText field based on attrName
		var insertText string
		switch attrName {
		case "retryable_errors":
			insertText = fmt.Sprintf("%s = [$0]", attrName)
		case "skip", "prevent_destroy":
			insertText = fmt.Sprintf("%s = true", attrName)
		case "iam_assume_role_duration":
			insertText = fmt.Sprintf("%s = 14400", attrName)
		default:
			insertText = fmt.Sprintf("%s = \"$0\"", attrName)
		}

		langDoc := LangDoc{
			Item:       attrName,
			InsertText: insertText,
			Content:    content,
		}

		langDocs = append(langDocs, langDoc)
	}

	return langDocs
}

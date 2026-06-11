package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tejeshwarmishra/upl/pkg/parser"
)

var (
	version = "1.0.0"
	helpMsg = `
UPL - Universal Programming Language Compiler/Runtime
Version: %s (Production Stable Core)

Usage:
  upl dev <file>        Start development server on :3000 with hot reload
  upl build <file>      Compile and bundle into distribution assets
  upl convert <file>    Convert standard .js/.html to .upl format
  upl revert <file>     Convert .upl back to native .js/.html
  upl init <lang-pack>  Register/stub a language pack

Example:
  upl dev app.upl
  upl build app.upl

`
	// Hardcoded fallback configuration profile baked inside the binary memory frame
	defaultJSONConfig = `{
  "language": "hin",
  "builtins": {
    "दिखाओ": "console.log",
    "खिड़की": "window",
    "दस्तावेज़": "document",
    "गणित": "Math",
    "जेसन": "JSON",
    "दिनांक": "Date",
    "बदलोसंख्या": "parseInt",
    "बदलोदशमलव": "parseFloat"
  }
}`
)

// UPLConfig handles runtime built-in macro-substitution mappings
type UPLConfig struct {
	Language string            `json:"language"`
	Builtins map[string]string `json:"builtins"`
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf(helpMsg, version)
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "dev":
		if len(args) < 2 {
			fmt.Println("Error: dev command requires a file argument")
			os.Exit(1)
		}
		devServer(args[1])

	case "build":
		if len(args) < 2 {
			fmt.Println("Error: build command requires a file argument")
			os.Exit(1)
		}
		buildCommand(args[1])

	case "convert":
		if len(args) < 2 {
			fmt.Println("Error: convert command requires a file argument")
			os.Exit(1)
		}
		convertCommand(args[1])

	case "revert":
		if len(args) < 2 {
			fmt.Println("Error: revert command requires a file argument")
			os.Exit(1)
		}
		revertCommand(args[1])

	case "init":
		if len(args) < 2 {
			fmt.Println("Error: init requires a language pack name")
			os.Exit(1)
		}
		registerLangPack(args[1])

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Printf(helpMsg, version)
		os.Exit(1)
	}
}

// loadRuntimeConfig reads local config patterns or loads embedded safe properties
func loadRuntimeConfig() UPLConfig {
	var config UPLConfig
	file, err := os.Open("upl.config.json")
	if err != nil {
		_ = json.Unmarshal([]byte(defaultJSONConfig), &config)
		return config
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&config)
	return config
}

// devServer starts a hot-reload development server on :3000
func devServer(filePath string) {
	fmt.Printf("🚀 Starting UPL Dev Server on http://localhost:3000\n")
	fmt.Printf("📄 Watching: %s\n", filePath)

	var lastModTime time.Time
	var lastContent string

	config := loadRuntimeConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Modernized I/O Read standard
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file: %v", err), 500)
			return
		}

		stat, _ := os.Stat(filePath)
		if stat != nil {
			modTime := stat.ModTime()
			if modTime.After(lastModTime) {
				lastModTime = modTime
				fmt.Printf("♻️  Recompiling...\n")
			}
		}

		sourceCode := string(content)
		if sourceCode == lastContent {
			w.Header().Set("X-Cache", "HIT")
		} else {
			w.Header().Set("X-Cache", "MISS")
			lastContent = sourceCode
		}

		// Process structural flattening execution pass
		flattened := flattenImports(sourceCode, filePath)

		// Dictionary Interception Strategy before pipeline parsing
		for hindiBuiltin, jsBuiltin := range config.Builtins {
			flattened = strings.ReplaceAll(flattened, hindiBuiltin, jsBuiltin)
		}

		result := compile(flattened)

		html := generateHTML(result.JavaScript, result.CSS)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	// Polling File Watcher Loop
	go func() {
		var prevModTime time.Time
		for {
			stat, err := os.Stat(filePath)
			if err == nil {
				if stat.ModTime() != prevModTime && !prevModTime.IsZero() {
					fmt.Printf("\n📝 File changed, refresh your browser...\n")
				}
				prevModTime = stat.ModTime()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

// buildCommand compiles UPL into distribution-ready assets
func buildCommand(filePath string) {
	fmt.Printf("🔨 Building: %s\n", filePath)

	config := loadRuntimeConfig()

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	flattened := flattenImports(string(content), filePath)

	for hindiBuiltin, jsBuiltin := range config.Builtins {
		flattened = strings.ReplaceAll(flattened, hindiBuiltin, jsBuiltin)
	}

	result := compile(flattened)

	outputDir := "dist"
	_ = os.MkdirAll(outputDir, 0755)

	html := generateHTML(result.JavaScript, result.CSS)
	htmlFile := filepath.Join(outputDir, "index.html")
	_ = os.WriteFile(htmlFile, []byte(html), 0644)
	fmt.Printf("✓ Generated: %s\n", htmlFile)

	jsFile := filepath.Join(outputDir, "app.js")
	_ = os.WriteFile(jsFile, []byte(result.JavaScript), 0644)
	fmt.Printf("✓ Generated: %s\n", jsFile)

	cssFile := filepath.Join(outputDir, "app.css")
	_ = os.WriteFile(cssFile, []byte(result.CSS), 0644)
	fmt.Printf("✓ Generated: %s\n", cssFile)

	fmt.Printf("\n✅ Build complete! Files in: %s/\n", outputDir)
}

// convertCommand converts .js/.html to .upl format
func convertCommand(filePath string) {
	fmt.Printf("🔄 Converting: %s to UPL format\n", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	converted := convertToUPL(string(content))

	outputFile := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".upl"
	_ = os.WriteFile(outputFile, []byte(converted), 0644)
	fmt.Printf("✓ Converted: %s\n", outputFile)
}

// revertCommand converts .upl back to native .js/.html
func revertCommand(filePath string) {
	fmt.Printf("⏮️  Reverting: %s to native format\n", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	result := compile(string(content))

	outputFile := strings.TrimSuffix(filePath, ".upl") + ".js"
	_ = os.WriteFile(outputFile, []byte(result.JavaScript), 0644)
	fmt.Printf("✓ Generated: %s\n", outputFile)

	if result.CSS != "" {
		cssFile := strings.TrimSuffix(filePath, ".upl") + ".css"
		_ = os.WriteFile(cssFile, []byte(result.CSS), 0644)
		fmt.Printf("✓ Generated: %s\n", cssFile)
	}
}

// registerLangPack generates a local environment translation configuration configuration profile
func registerLangPack(packName string) {
	fmt.Printf("📦 Initializing runtime environment language configuration pack for: %s\n", packName)
	err := os.WriteFile("upl.config.json", []byte(defaultJSONConfig), 0644)
	if err != nil {
		fmt.Printf("❌ Configuration generation failure: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Local operational architecture configuration configuration profile stamped to configuration file successfully: ./upl.config.json")
}

// CompileResult holds the output of compilation
type CompileResult struct {
	JavaScript string
	CSS        string
	Errors     []*parser.ParseError
}

// compile converts UPL source to JavaScript and CSS
func compile(source string) *CompileResult {
	p := parser.New(source)
	ast := p.Parse()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(os.Stderr, "❌ Compilation Errors:\n")
		for _, err := range p.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}

	js := ast.ToJavaScript()
	css := extractCSS(ast)

	return &CompileResult{
		JavaScript: js,
		CSS:        css,
		Errors:     p.Errors(),
	}
}

// extractCSS pulls CSS blocks from the AST
func extractCSS(node parser.ASTNode) string {
	var css strings.Builder

	walkAST(node, func(n parser.ASTNode) {
		if style, ok := n.(*parser.StyleBlock); ok {
			css.WriteString(style.ToCSS())
			css.WriteString("\n")
		}
	})

	return css.String()
}

// walkAST recursively traverses the AST
func walkAST(node parser.ASTNode, fn func(parser.ASTNode)) {
	if node == nil {
		return
	}

	fn(node)

	switch n := node.(type) {
	case *parser.Block:
		for _, stmt := range n.Statements {
			walkAST(stmt, fn)
		}
	case *parser.FuncDecl:
		if n.Body != nil {
			walkAST(n.Body, fn)
		}
	case *parser.IfStmt:
		walkAST(n.Condition, fn)
		walkAST(n.Then, fn)
		if n.Else != nil {
			walkAST(n.Else, fn)
		}
	case *parser.WhileStmt:
		walkAST(n.Condition, fn)
		if n.Body != nil {
			walkAST(n.Body, fn)
		}
	case *parser.ForStmt:
		walkAST(n.Init, fn)
		walkAST(n.Condition, fn)
		walkAST(n.Update, fn)
		if n.Body != nil {
			walkAST(n.Body, fn)
		}
	case *parser.BinaryOp:
		walkAST(n.Left, fn)
		walkAST(n.Right, fn)
	case *parser.CallExpr:
		walkAST(n.Function, fn)
		for _, arg := range n.Args {
			walkAST(arg, fn)
		}
	case *parser.LogStmt:
		for _, arg := range n.Args {
			walkAST(arg, fn)
		}
	}
}

// flattenImports resolves and includes imported files (जोड़ो)
func flattenImports(source string, baseDir string) string {
	var result strings.Builder
	lines := strings.Split(source, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "जोड़ो ") || strings.HasPrefix(trimmed, "import ") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				filename := strings.Trim(parts[1], `"'`)
				importPath := filepath.Join(filepath.Dir(baseDir), filename)

				if content, err := os.ReadFile(importPath); err == nil {
					result.WriteString(string(content))
					result.WriteString("\n")
					continue
				}
			}
		}

		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}

// generateHTML wraps compiled JS/CSS into an HTML template
func generateHTML(js, css string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>UPL App</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
                'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
                sans-serif;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
            background: #f5f5f5;
        }

        #app {
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
    </style>
    <style>
%s
    </style>
</head>
<body>
    <div id="app"></div>
    
    <script>
        window.upl = {
            render: (html) => {
                document.getElementById('app').innerHTML = html;
            }
        };

        // User-Generated JavaScript
        %s
    </script>
</body>
</html>`, css, js)
}

// convertToUPL converts native JS to UPL format (basic mapping layout)
func convertToUPL(code string) string {
	converted := code
	conversions := map[string]string{
		"const ":       "स्थिर बनाओ ",
		"let ":         "बदलाव योग्य ",
		"var ":         "अपरिवर्तनीय ",
		"function ":    "कार्य बनाओ ",
		"return ":      "वापस करो ",
		"if (":         "यदि (",
		"} else {":     "} वरना {",
		"while (":      "जब तक (",
		"for (":        "के लिए (",
		"class ":       "श्रेणी ",
		"console.log": "दिखाओ",
		"async ":       "अतुल्यकालिक ",
		"await ":       "इंतजार करो ",
	}

	for js, upl := range conversions {
		converted = strings.ReplaceAll(converted, js, upl)
	}

	return converted
}
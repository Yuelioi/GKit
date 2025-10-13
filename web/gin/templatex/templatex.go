package templatex

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Yuelioi/gkit/web/gin/templatex/conf"
	"github.com/Yuelioi/gkit/web/gin/templatex/parser"
)

//go:embed templates/*
var defaultTemplatesFS embed.FS

type Generator struct {
	templatePath string // 自定义模板路径(可选)
	useEmbed     bool   // 是否使用内嵌模板
}

// NewGenerator 创建生成器
func NewGenerator(opts ...GeneratorOption) *Generator {
	g := &Generator{
		useEmbed: true, // 默认使用内嵌模板
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

type GeneratorOption func(*Generator)

// WithCustomTemplates 使用自定义模板目录
func WithCustomTemplates(path string) GeneratorOption {
	return func(g *Generator) {
		g.templatePath = path
		g.useEmbed = false
	}
}

func (g *Generator) getTemplateStr(tplPath string) (string, error) {
	if g.useEmbed {
		// 使用内嵌模板
		data, err := defaultTemplatesFS.ReadFile(path.Join("templates", tplPath))
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// 使用自定义模板路径
	p := filepath.Join(g.templatePath, tplPath)
	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (g *Generator) generateCode(modelInfo *conf.ModelInfo, config *conf.Config) error {
	outputDir := config.Output.Dir

	handlerDir := filepath.Join(outputDir, config.Output.Structure["handler"])
	routeDir := filepath.Join(outputDir, config.Output.Structure["route"])

	os.MkdirAll(handlerDir, 0755)
	os.MkdirAll(routeDir, 0755)

	// 生成 Handler
	handlerPath := filepath.Join(handlerDir, modelInfo.LowerName+"_handler.go")
	if !config.Output.Overwrite {
		if _, err := os.Stat(handlerPath); err == nil {
			fmt.Printf("⚠ Skipping %s (already exists)\n", handlerPath)
			goto generateRouter
		}
	}

	{
		handlerFile, err := os.Create(handlerPath)
		if err != nil {
			return err
		}
		defer handlerFile.Close()

		handler, err := g.getTemplateStr("handler.tmpl")
		if err != nil {
			return err
		}

		tmpl := template.Must(template.New("handler").Parse(handler))
		if err := tmpl.Execute(handlerFile, modelInfo); err != nil {
			return err
		}
		fmt.Printf("✓ Generated: %s\n", handlerPath)
	}

generateRouter:
	// 生成 Router
	routerPath := filepath.Join(routeDir, modelInfo.LowerName+"_routes.go")
	if !config.Output.Overwrite {
		if _, err := os.Stat(routerPath); err == nil {
			fmt.Printf("⚠ Skipping %s (already exists)\n", routerPath)
			return nil
		}
	}

	{
		routerFile, err := os.Create(routerPath)
		if err != nil {
			return err
		}
		defer routerFile.Close()

		router, err := g.getTemplateStr("router.tmpl")
		if err != nil {
			return err
		}

		tmpl := template.Must(template.New("router").Parse(router))
		if err := tmpl.Execute(routerFile, modelInfo); err != nil {
			return err
		}
		fmt.Printf("✓ Generated: %s\n", routerPath)
	}

	return nil
}

func (g *Generator) GenerateExampleConfig(filename string) error {
	config, err := g.getTemplateStr("config.yaml")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(config), 0644)
}

func (g *Generator) GenerateModel(modelFile, configFile string) error {
	config, err := conf.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if config.Model.File != "" {
		modelFile = config.Model.File
	}

	fmt.Printf("📝 Parsing model file: %s\n", modelFile)
	modelInfo, err := parser.ParseModel(modelFile, config)
	if err != nil {
		return fmt.Errorf("error parsing model: %w", err)
	}

	if modelInfo.Name == "" {
		return fmt.Errorf("no valid model found in file")
	}

	fmt.Printf("🚀 Generating CRUD code for model: %s\n", modelInfo.Name)
	if err := g.generateCode(modelInfo, config); err != nil {
		return fmt.Errorf("error generating code: %w", err)
	}

	fmt.Println("\n✅ Generation completed successfully!")
	return nil
}

// 便捷函数 - 使用默认内嵌模板
func GenerateModel(modelFile, configFile string) error {
	gen := NewGenerator()
	return gen.GenerateModel(modelFile, configFile)
}

func GenerateExampleConfig(filename string) error {
	gen := NewGenerator()
	return gen.GenerateExampleConfig(filename)
}

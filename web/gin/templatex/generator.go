package templatex

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Yuelioi/gkit/web/gin/templatex/conf"
	"github.com/Yuelioi/gkit/web/gin/templatex/parser"
)

//go:embed templates/*
var defaultTemplatesFS embed.FS

type Generator struct {
	templatePath string
	useEmbed     bool
	funcMap      template.FuncMap
}

func NewGenerator(opts ...GeneratorOption) *Generator {
	g := &Generator{
		useEmbed: true,
		funcMap:  createTemplateFuncMap(),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

type GeneratorOption func(*Generator)

func WithCustomTemplates(path string) GeneratorOption {
	return func(g *Generator) {
		g.templatePath = path
		g.useEmbed = false
	}
}

// 创建模板函数映射
func createTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"toLower":    strings.ToLower,
		"toUpper":    strings.ToUpper,
		"title":      strings.Title,
		"camelCase":  toCamelCase,
		"snakeCase":  toSnakeCase,
		"kebabCase":  toKebabCase,
		"plural":     toPlural,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"join":       strings.Join,
		"split":      strings.Split,
	}
}

func (g *Generator) getTemplateStr(tplPath string) (string, error) {
	if g.useEmbed {
		data, err := defaultTemplatesFS.ReadFile(path.Join("templates", tplPath))
		if err != nil {
			return "", fmt.Errorf("❌ 读取内嵌模板失败: %v", err)
		}
		return string(data), nil
	}

	p := filepath.Join(g.templatePath, tplPath)
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("❌ 读取自定义模板失败: %v", err)
	}
	return string(data), nil
}

func (g *Generator) generateCode(modelInfo *conf.ModelInfo, config *conf.Config) error {
	outputDir := config.Output.Dir

	// 创建所有需要的目录
	dirs := []string{
		filepath.Join(outputDir, config.Output.Structure["handler"]),
		filepath.Join(outputDir, config.Output.Structure["route"]),
		filepath.Join(outputDir, config.Output.Structure["service"]),
		filepath.Join(outputDir, config.Output.Structure["dto"]),
		filepath.Join(outputDir, config.Output.Structure["repository"]),
		filepath.Join(outputDir, config.Output.Structure["middleware"]),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("❌ 创建目录失败 %s: %v", dir, err)
		}
	}

	// 生成各个层的代码
	generators := []struct {
		name     string
		template string
		path     string
		enabled  bool
	}{
		{"Handler", "handler.tmpl", filepath.Join(config.Output.Structure["handler"], modelInfo.LowerName+"_handler.go"), true},
		{"Router", "router.tmpl", filepath.Join(config.Output.Structure["route"], modelInfo.LowerName+"_routes.go"), true},
		{"Service", "service.tmpl", filepath.Join(config.Output.Structure["service"], modelInfo.LowerName+"_service.go"), config.Features.Service},
		{"DTO", "dto.tmpl", filepath.Join(config.Output.Structure["dto"], modelInfo.LowerName+"_dto.go"), config.Features.DTO},
		{"Repository", "repository.tmpl", filepath.Join(config.Output.Structure["repository"], modelInfo.LowerName+"_repository.go"), config.Features.Repository},
		{"Test", "handler_test.tmpl", filepath.Join(config.Output.Structure["handler"], modelInfo.LowerName+"_handler_test.go"), config.Features.Test},
		{"Mock", "mock.tmpl", filepath.Join(config.Output.Structure["handler"], "mock_"+modelInfo.LowerName+".go"), config.Features.Mock},
	}

	for _, gen := range generators {
		if !gen.enabled {
			continue
		}

		fullPath := filepath.Join(outputDir, gen.path)

		// 检查是否覆盖
		if !config.Output.Overwrite {
			if _, err := os.Stat(fullPath); err == nil {
				fmt.Printf("⚠️  已跳过：%s (文件已存在)\n", fullPath)
				continue
			}
		}

		if err := g.generateFile(gen.template, fullPath, modelInfo); err != nil {
			return fmt.Errorf("❌ 生成 %s 失败: %v", gen.name, err)
		}
		fmt.Printf("✅ 已生成 %s: %s\n", gen.name, fullPath)
	}

	return nil
}

func (g *Generator) generateFile(templateName, outputPath string, data interface{}) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	templateStr, err := g.getTemplateStr(templateName)
	if err != nil {
		return err
	}

	tmpl, err := template.New(templateName).Funcs(g.funcMap).Parse(templateStr)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, data)
}

func (g *Generator) GenerateExampleConfig(filename string) error {
	config, err := g.getTemplateStr("config.yaml")
	if err != nil {
		return fmt.Errorf("❌ 读取示例配置模板失败: %v", err)
	}
	if err := os.WriteFile(filename, []byte(config), 0644); err != nil {
		return fmt.Errorf("❌ 写入示例配置文件失败: %v", err)
	}
	fmt.Printf("📄 已生成示例配置文件: %s\n", filename)
	return nil
}

func (g *Generator) GenerateModel(configFile string) error {
	config, err := conf.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("❌ 加载配置文件失败: %v", err)
	}

	if len(config.Model.Files) == 0 {
		return fmt.Errorf("⚠️  模型文件列表为空，请检查配置文件！")
	}

	fmt.Println("🚀 开始生成代码...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, f := range config.Model.Files {
		fmt.Printf("\n📝 正在解析模型文件: %s\n", f)
		modelInfo, err := parser.ParseModel(f, config)
		if err != nil {
			return fmt.Errorf("❌ 解析模型失败: %v", err)
		}

		if modelInfo.Name == "" {
			fmt.Printf("⚠️  未找到有效模型，跳过文件: %s\n", f)
			continue
		}

		fmt.Printf("🔨 为模型 [%s] 生成代码...\n", modelInfo.Name)
		if err := g.generateCode(modelInfo, config); err != nil {
			return fmt.Errorf("❌ 生成代码出错: %v", err)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 代码生成完成！")

	// 生成使用说明
	g.printUsageGuide(config)

	return nil
}

func (g *Generator) printUsageGuide(config *conf.Config) {
	fmt.Println("\n📖 使用指南:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("1. 在 main.go 中注册路由:\n")
	fmt.Printf("   routes.Setup{ModelName}Routes(r, db)\n\n")

	if config.Features.Service {
		fmt.Printf("2. Service 层已生成，包含业务逻辑\n")
	}

	if config.Features.Repository {
		fmt.Printf("3. Repository 层已生成，包含数据访问逻辑\n")
	}

	if config.Features.Test {
		fmt.Printf("4. 运行测试: go test ./...\n")
	}

	fmt.Println("\n💡 提示:")
	fmt.Println("   - 可以在配置文件中自定义更多选项")
	fmt.Println("   - 支持自定义模板，使用 WithCustomTemplates() 选项")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// 便捷函数
func GenerateModel(configFile string) error {
	gen := NewGenerator()
	return gen.GenerateModel(configFile)
}

func GenerateExampleConfig(filename string) error {
	gen := NewGenerator()
	return gen.GenerateExampleConfig(filename)
}

// 工具函数
func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}

func toPlural(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

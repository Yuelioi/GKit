package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/Yuelioi/gkit/web/gin/templatex/conf"
)

func ParseModel(filePath string, config *conf.Config) (*conf.ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var modelInfo conf.ModelInfo
	modelInfo.PkgName = node.Name.Name
	modelInfo.Config = config

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		modelName := typeSpec.Name.Name

		// 检查是否排除
		for _, exclude := range config.Model.Exclude {
			if modelName == exclude {
				return false
			}
		}

		// 检查是否包含
		if len(config.Model.Include) > 0 {
			found := false
			for _, include := range config.Model.Include {
				if modelName == include {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		modelInfo.Name = modelName
		modelInfo.LowerName = strings.ToLower(string(modelInfo.Name[0])) + modelInfo.Name[1:]
		modelInfo.PluralName = modelInfo.Name + "s"

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			fieldInfo := conf.FieldInfo{
				Name: field.Names[0].Name,
				Type: getTypeString(field.Type),
			}

			if field.Tag != nil {
				tag := field.Tag.Value
				fieldInfo.JsonTag = extractTag(tag, "json")
				fieldInfo.GormTag = extractTag(tag, "gorm")
				if strings.Contains(fieldInfo.GormTag, "primaryKey") {
					fieldInfo.IsID = true
				}
			}

			// 检查是否是可搜索字段
			if config.Features.Search {
				for _, searchField := range config.Features.SearchFields {
					if fieldInfo.JsonTag == searchField || fieldInfo.Name == searchField {
						fieldInfo.Searchable = true
						break
					}
				}
			}

			// 检查时间字段
			if strings.Contains(fieldInfo.Type, "time.Time") {
				modelInfo.HasTimeField = true
			}

			modelInfo.Fields = append(modelInfo.Fields, fieldInfo)
		}

		return false
	})

	return &modelInfo, nil
}

func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

func extractTag(tag, key string) string {
	tag = strings.Trim(tag, "`")
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, key+":") {
			value := strings.TrimPrefix(part, key+":")
			value = strings.Trim(value, "\"")
			if strings.Contains(value, ",") {
				return strings.Split(value, ",")[0]
			}
			return value
		}
	}
	return ""
}

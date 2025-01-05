package util

import (
	"text/template"
	"text/template/parse"

	"github.com/samber/lo"
)

// ExtractVariablesFromContent 从模板中提取变量
//
//	param content string
//	return []string
//	author centonhuang
//	update 2024-12-09 17:22:41
func ExtractVariablesFromContent(content string) []string {
	// 创建模板并解析内容
	tmpl, err := template.New("prompt").Parse(content)
	if err != nil {
		return []string{}
	}

	// 用于存储找到的变量
	var variables []string

	// 遍历语法树以找到所有变量
	var extractVars func(node parse.Node)
	extractVars = func(node parse.Node) {
		if node == nil {
			return
		}

		switch n := node.(type) {
		case *parse.ActionNode:
			// 处理 {{.VarName}} 形式的节点
			if len(n.Pipe.Cmds) > 0 && len(n.Pipe.Cmds[0].Args) > 0 {
				if field, ok := n.Pipe.Cmds[0].Args[0].(*parse.FieldNode); ok {
					if len(field.Ident) > 0 {
						variables = append(variables, field.Ident[0])
					}
				}
			}
		case *parse.ListNode:
			for _, n := range n.Nodes {
				extractVars(n)
			}
		}
	}

	extractVars(tmpl.Tree.Root)

	return lo.Uniq(variables)
}

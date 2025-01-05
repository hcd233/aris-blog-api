package document

import (
	"fmt"

	"github.com/samber/lo"
)

// Relation 关系
//
//	author centonhuang
//	update 2024-12-07 14:23:01
type Relation struct {
	TextDocument
	Source Document
	Target Document
}

// Map 转换为map
//
//	receiver r *Relation
//	return map
//	author centonhuang
//	update 2024-12-07 14:38:45
func (r *Relation) Map() map[string]interface{} {
	return lo.Assign(map[string]interface{}{
		"source": r.Source.Map(),
		"target": r.Target.Map(),
	}, r.TextDocument.Map())
}

// String 转换为字符串
//
//	receiver r *Relation
//	return string
//	author centonhuang
//	update 2024-12-07 14:39:15
func (r *Relation) String() string {
	str := r.TextDocument.String()
	str += fmt.Sprintf("\t<source>\n\t%s\n\t</source>\n", r.Source.String())
	str += fmt.Sprintf("\t<target>\n\t%s\n\t</target>\n", r.Target.String())
	return str
}

// GraphDocument 图谱文档
//
//	author centonhuang
//	update 2024-12-07 14:23:07
type GraphDocument struct {
	Name      string
	Relations []Relation
}

// Map 转换为map
//
//	receiver gd *GraphDocument
//	return map
//	author centonhuang
//	update 2024-12-07 14:30:44
func (gd *GraphDocument) Map() map[string]interface{} {
	return map[string]interface{}{
		"name": gd.Name,
		"relations": lo.Map(gd.Relations, func(relation Relation, _ int) map[string]interface{} {
			return relation.Map()
		}),
	}
}

func (gd *GraphDocument) String() string {
	str := fmt.Sprintf("<name>\n%s\n</name>\n", gd.Name)
	str += "<relations>\n"
	for idx, relation := range gd.Relations {
		str += fmt.Sprintf("\t<relation_%d>\n\t%s\n\t</relation_%d>\n", idx, relation.String(), idx)
	}
	str += "</relations>\n"
	return str
}

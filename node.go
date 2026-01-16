package querydsl

// NodeType 定义节点类型
type NodeType string

const (
	NodeAnd  NodeType = "and"  // 逻辑与
	NodeOr   NodeType = "or"   // 逻辑或
	NodeCond NodeType = "cond" // 条件节点，只有条件节点才有 Payload （Condition）
)

// Node 表示 Query DSL 的节点
type Node struct {
	Type     NodeType
	Children []*Node
	Payload  any // Condition
}

func Cond(field string, op Operator, value any) *Node {
	return &Node{
		Type: NodeCond,
		Payload: Condition{
			Field: field,
			Op:    op,
			Value: value,
		},
	}
}

func Or(children ...*Node) *Node {
	return &Node{
		Type:     NodeOr,
		Children: children,
	}
}

func And(children ...*Node) *Node {
	return &Node{
		Type:     NodeAnd,
		Children: children,
	}
}

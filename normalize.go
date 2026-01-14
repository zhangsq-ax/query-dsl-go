package querydsl

// Normalize 会：
// 1. 合并同类型 AND / OR 子节点
// 2. 删除空节点
func Normalize(n *Node) *Node {
	if n == nil {
		return nil
	}

	// leaf
	if n.Type == NodeCond {
		return n
	}

	var children []*Node
	for _, c := range n.Children {
		cn := Normalize(c)
		if cn == nil {
			continue
		}

		// 扁平化同 Op
		if cn.Type == n.Type {
			children = append(children, cn.Children...)
		} else {
			children = append(children, cn)
		}
	}

	if len(children) == 0 {
		return nil
	}
	if len(children) == 1 {
		return children[0]
	}

	n.Children = children
	return n
}

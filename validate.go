package querydsl

import "errors"

func Validate(n *Node) error {
	if n == nil {
		return nil
	}

	switch n.Type {
	case NodeAnd, NodeOr:
		if len(n.Children) < 1 {
			return errors.New("logic node must have children")
		}
		for _, c := range n.Children {
			if err := Validate(c); err != nil {
				return err
			}
		}
	case NodeCond:
		if n.Payload == nil {
			return errors.New("condition payload is required")
		}
		if _, ok := n.Payload.(Condition); !ok {
			return errors.New("payload must be querydsl.Condition")
		}
	default:
		return errors.New("unknown node type")
	}

	return nil
}

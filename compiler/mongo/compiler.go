package mongoCompiler

import (
	"errors"

	"github.com/zhangsq-ax/query-dsl-go"
	"go.mongodb.org/mongo-driver/bson"
)

type Compiler struct{}

func New() *Compiler {
	return &Compiler{}
}

// Compile 编译 Query DSL IR 节点为 MongoDB 查询过滤器
func (c *Compiler) Compile(node *querydsl.Node) (bson.M, error) {
	if node == nil {
		return bson.M{}, nil
	}

	switch node.Type {
	case querydsl.NodeAnd:
		return c.compileLogic("$and", node.Children)
	case querydsl.NodeOr:
		return c.compileLogic("$or", node.Children)
	case querydsl.NodeCond:
		return c.compileCondition(node.Payload)
	default:
		return nil, errors.New("unknown node type")
	}
}

func (c *Compiler) compileLogic(op string, children []*querydsl.Node) (bson.M, error) {
	var clauses []bson.M

	for _, child := range children {
		m, err := c.Compile(child)
		if err != nil || len(m) == 0 {
			continue
		}
		clauses = append(clauses, m)
	}

	if len(clauses) == 0 {
		return bson.M{}, nil
	}
	if len(clauses) == 1 {
		return clauses[0], nil
	}

	return bson.M{op: clauses}, nil
}

func (c *Compiler) compileCondition(payload any) (bson.M, error) {
	cond, ok := payload.(querydsl.Condition)
	if !ok {
		return nil, errors.New("invalid condition payload")
	}

	field := cond.Field

	switch cond.Op {

	case querydsl.OpEq:
		return bson.M{field: cond.Value}, nil

	case querydsl.OpNe:
		return bson.M{
			field: bson.M{"$ne": cond.Value},
		}, nil

	case querydsl.OpIn:
		return bson.M{
			field: bson.M{"$in": cond.Value},
		}, nil

	case querydsl.OpLike:
		s, ok := cond.Value.(string)
		if !ok {
			return nil, errors.New("like value must be string")
		}
		return bson.M{
			field: bson.M{
				"$regex": ".*?" + s + ".*?",
			},
		}, nil

	case querydsl.OpLikeI:
		s, ok := cond.Value.(string)
		if !ok {
			return nil, errors.New("like value must be string")
		}
		return bson.M{
			field: bson.M{
				"$regex":   ".*?" + s + ".*?",
				"$options": "i",
			},
		}, nil

	case querydsl.OpPrefixLike:
		s, ok := cond.Value.(string)
		if !ok {
			return nil, errors.New("like value must be string")
		}
		return bson.M{
			field: bson.M{
				"$regex": "^" + s + ".*?",
			},
		}, nil

	case querydsl.OpExists:
		return bson.M{
			field: bson.M{
				"$exists": cond.Value,
			},
		}, nil

	case querydsl.OpRange:
		val, ok := cond.Value.([2]any)
		if !ok {
			return nil, errors.New("range value must be [2]any")
		}

		rangeExpr := bson.M{}
		if val[0] != nil {
			rangeExpr["$gte"] = val[0]
		}
		if val[1] != nil {
			rangeExpr["$lte"] = val[1]
		}

		return bson.M{field: rangeExpr}, nil
	}

	return nil, errors.New("unsupported operator")
}

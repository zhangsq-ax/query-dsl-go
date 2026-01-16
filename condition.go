package querydsl

// Operator 条件操作符
type Operator string

const (
	OpEq         Operator = "="           // 等于
	OpNe         Operator = "!="          // 不等于
	OpIn         Operator = "in"          // 包含于
	OpRange      Operator = "range"       // 范围
	OpLike       Operator = "like"        // 模糊匹配
	OpLikeI      Operator = "like_i"      // 不区分大小写模糊匹配
	OpPrefixLike Operator = "prefix_like" // 前缀匹配
	OpExists     Operator = "exists"      // 存在
)

// Condition 查询条件
type Condition struct {
	Field string   // 字段名
	Op    Operator // 操作符
	Value any      // 值
}

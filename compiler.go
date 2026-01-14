package querydsl

// Compiler 将 Query DSL IR 编译为某种后端结构
// T 由具体 compiler 决定（bson.M / SQL string / ES DSL 等）
type Compiler[T any] interface {
	// Compile 编译节点
	Compile(node *Node) (T, error)
}



# querydsl-go

> **A backend-agnostic query DSL for building complex list queries**
> 一个可复用、可编译、可解释的查询表达式内核，用于构建复杂列表查询。

------

## ✨ 背景 & 动机

在实际业务中，**列表查询**往往会不断演进：

- 条件越来越多（AND / OR / 嵌套）
- 查询来源多样（HTTP API / WebSocket / 内部 RPC）
- 底层存储可能变化（MongoDB / SQL / ES）
- 需要 Debug、Explain、审计、复用

如果直接在业务中拼装数据库查询语句，常见问题包括：

- 查询逻辑分散在各个接口中
- 难以统一校验、优化、调试
- 强耦合某一种数据库
- 无法复用到其他业务

**querydsl 的目标**是解决这些问题：

> 👉 **把“查询逻辑”抽象成一棵 DSL 表达树，再由不同 Compiler 编译为具体存储的查询语句**

------

## 🎯 设计目标

- ✅ 后端无关（Mongo / SQL / ES 可共存）
- ✅ 强结构化（可校验、可规范化）
- ✅ Explain-friendly（可读、可调试）
- ✅ 适合复杂列表查询
- ✅ 不侵入业务 DTO

------

## 🧩 核心概念

### 1️⃣ Query Node（IR）

所有查询最终都会转化为一棵**中间表示树（IR）**：

```text
        AND
       /   \
   COND     OR
           / \
       COND  COND
```

对应结构：

```go
type Node struct {
	Type     NodeType      // And / Or / Cond
	Payload  Condition     // 仅 Cond 节点有
	Children []*Node
}
```

------

### 2️⃣ Condition

最小查询单元：

```go
type Condition struct {
	Field string
	Op    Operator
	Value any
}
```

示例：

```go
Condition{
	Field: "status",
	Op:    OpIn,
	Value: []string{"online", "offline"},
}
```

------

### 3️⃣ Operator（内置）

| Operator   | 含义     | Mongo 示例              |
| ---------- | -------- | ----------------------- |
| `OpEq`     | 等于     | `{field: value}`        |
| `OpIn`     | 包含     | `{field: {$in: [...]}}` |
| `OpRange`  | 区间     | `{field: {$gte, $lte}}` |
| `OpLike`   | 模糊匹配 | `$regex`                |
| `OpExists` | 字段存在 | `$exists`               |

> Operator 本身不绑定数据库语义，由 Compiler 解释。

------

## 🔧 Normalize（规范化）

在编译前，querydsl 会对 Node 进行规范化：

- 合并同类逻辑节点（AND / OR）
- 消除空节点
- 单子节点折叠
- 保证输出结构 explain-friendly

```go
node = querydsl.Normalize(node)
```

**Normalize 后的查询更稳定、可预测、便于调试**

------

## 🛡 Validate（校验）

querydsl 提供严格校验能力：

- 非法 Operator
- 缺失 Field / Value
- 逻辑节点为空
- 不合理的 Range / In

```go
if err := querydsl.Validate(node); err != nil {
	return err
}
```

------

## 🧠 Compiler 抽象

querydsl 本身**不关心数据库**，只定义 Compiler 协议：

```go
type Compiler[T any] interface {
	Compile(node *Node) (T, error)
}
```

------

### Mongo Compiler 示例

```go
compiler := mongocompiler.New()
filter, err := compiler.Compile(node)
// filter => bson.M
```

支持：

- `$and / $or`
- `$in`
- `$regex`
- `$exists`
- range（`$gte / $lte`）

------

## 🚀 一个完整示例

### 构建查询

```go
node := &querydsl.Node{
  Type: querydsl.NodeAnd,
  Children: []*querydsl.Node{
    {
      Type: querydsl.NodeCond,
      Payload: querydsl.Condition{
        Field: "status",
        Op: querydsl.OpIn,
        Value: []string{"online", "offline"},
    },
    {
      Type: querydsl.NodeOr,
      Children: []*querydal.Node{
        {
          Type: querydsl.NodeCond,
          Payload: querydsl.Condition{
            Field: "name",
            Op: querydsl.OpLike,
            Value: "robot",
          },
        },
        {
          Type: querydsl.NodeCond,
          Payload: querydsl.Condition{
            Field: "sn",
            Op: querydsl.OpEq,
            Value: "R123",
          },
        },
      }
    },    
  }
}
```

### Normalize + Compile

```go
node = querydsl.Normalize(node)

compiler := mongocompiler.New()
filter, _ := compiler.Compile(node)
```

### 输出 Mongo 查询

```json
{
  "$and": [
    { "status": { "$in": ["online", "offline"] } },
    {
      "$or": [
        { "name": { "$regex": "robot", "$options": "i" } },
        { "sn": "R123" }
      ]
    }
  ]
}
```

------

## 🧱 推荐使用方式（最佳实践）

### ✔ API 层

- 使用 **业务 DTO / Expr**
- 不直接接触 querydsl

### ✔ Service 层

```text
Expr → querydsl.Node → Normalize → Compile
```

### ✔ Storage 层

- 只接收已编译的查询结构
- 不包含业务逻辑

------

## ❌ 不推荐的使用方式

- 在 Handler 中直接拼 Node
- 在 core 中引入 bson / sql
- 用 querydsl 替代 ORM

------

## 🧩 模块结构（推荐）

```text
querydsl/
├── node.go
├── condition.go
├── operator.go
├── normalize.go
├── validate.go
├── compiler.go        // Compiler[T] interface
│
├── compiler/
│   ├── mongo/
│   ├── sql/
│   └── es/
```

------

## 🧪 适用场景

- 列表查询 / 搜索接口
- 权限可控的复杂过滤
- 动态查询条件
- Explain / Debug 需求强
- 多存储后端共存

------

## 📌 非目标（明确不做）

- ❌ ORM
- ❌ 自动索引选择
- ❌ 执行计划优化
- ❌ SQL 生成器替代品

------

## 🧭 总结

> **querydsl 是“查询逻辑的中间语言”，不是数据库工具**

如果你需要：

- 统一复杂查询逻辑
- 降低业务与数据库的耦合
- 提升可维护性与可调试性

那么 querydsl 会是一个**长期收益极高的基础模块**。

------

如果你愿意，我可以继续帮你补：

- `README.zh-CN` / `README.en`
- Compiler 扩展（Explain / Cost）
- SQL / PostgreSQL Compiler 示例
- 一套真实业务的落地范式说明
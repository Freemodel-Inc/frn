frn
------------------------------------

Freemodel Resource Name (frn) - go package to manage universal ids

### Tags

| annotation                   | example                              | description                                      |
|:-----------------------------|:-------------------------------------|:-------------------------------------------------|
| `frn=project`                | fm:crm:project:1                     | require unary form                               |
| `frn=project/`               | fm:crm:project:1:contract:2          | require compound form; unary form not acceptable |
| `frn=/contract`              | fm:crm:project:1:contract:2          | require child type                               |
| `frn=/entity/account`        | fm:crm:project:1:entity:2/account/ar | require child type with path head, key           |
| `frn=project/contract`       | fm:crm:project:1:contract:2          | require parent and child                         |
| `frn=project/entity/account` | fm:crm:project:1:entity:2/account/ar | require parent and child and path head, key      |

### Example

```go
type Input struct {
  ID frn.ID `validate:"required,frn=project"` // require id that must be a project id
}
```
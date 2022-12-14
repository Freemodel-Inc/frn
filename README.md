frn
------------------------------------

Freemodel Resource Name (frn) - go package to manage universal ids

### Tags

| annotation             | example                                                    | description                                      |
|:-----------------------|:-----------------------------------------------------------|:-------------------------------------------------|
| `frn=project`          | {namespace}:{service}:project:{id}                         | require unary form                               |
| `frn=project/`         | {namespace}:{service}:project:{id}:{child-type}:{child-id} | require compound form; unary form not acceptable |
| `frn=/contract`        | {namespace}:{service}:{type}:{id}:contract:{child-id}      | require child type                               |
| `frn=project/contract` | {namespace}:{service}:project:{id}:contract:{child-id      | require parent and child                         |


### Example

```go
type Input struct {
	ID frn.ID `validate:"required,frn=project"` // require id that must be a project id
}
```
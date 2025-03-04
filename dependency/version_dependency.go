package dependency

import "lucy/lucytypes"

type VersionOperator uint8

const (
	Equal VersionOperator = iota
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)

// Dependency can describe a dependency relationship. You MUST NOT use the
// Id's PackageId.Version field. Instead, you should use the Value and Operator.
type Dependency struct {
	Id       lucytypes.PackageId
	Value    SemanticVersion
	Operator VersionOperator
}

func (d Dependency) Satisfy(
	id lucytypes.PackageId,
	v SemanticVersion,
) bool {
	if (id.Platform != d.Id.Platform) || (id.Name != d.Id.Name) {
		return false
	}
	switch d.Operator {
	case Equal:
		return v.Eq(d.Value)
	case NotEqual:
		return v.Neq(d.Value)
	case GreaterThan:
		return v.Gt(d.Value)
	case GreaterThanOrEqual:
		return v.Gte(d.Value)
	case LessThan:
		return v.Lt(d.Value)
	case LessThanOrEqual:
		return v.Lte(d.Value)
	default:
		return false
	}
}

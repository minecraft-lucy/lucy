package lucytypes

import "lucy/dependency"

// Dependency can describe a dependency relationship. You MUST NOT use the
// Id's PackageId.Version field. Instead, you should use the Value and Operator.
type Dependency struct {
	Id       PackageId
	Value    dependency.SemanticVersion
	Operator dependency.VersionOperator
}

func (d Dependency) Satisfy(
id PackageId,
v dependency.SemanticVersion,
) bool {
	if (id.Platform != d.Id.Platform) || (id.Name != d.Id.Name) {
		return false
	}
	switch d.Operator {
	case dependency.Equal:
		return v.Eq(d.Value)
	case dependency.NotEqual:
		return v.Neq(d.Value)
	case dependency.GreaterThan:
		return v.Gt(d.Value)
	case dependency.GreaterThanOrEqual:
		return v.Gte(d.Value)
	case dependency.LessThan:
		return v.Lt(d.Value)
	case dependency.LessThanOrEqual:
		return v.Lte(d.Value)
	default:
		return false
	}
}

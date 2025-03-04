package dependency

func (p1 SemanticVersion) Eq(p2 SemanticVersion) bool {
	// If the labels are different, the versions are not comparable.
	if p1.Label != p2.Label {
		return false
	}
	return p1.Major == p2.Major && p1.Minor == p2.Minor && p1.Patch == p2.Patch
}

func (p1 SemanticVersion) Neq(p2 SemanticVersion) bool {
	return !p1.Eq(p2)
}

func (p1 SemanticVersion) StrictEq(p2 SemanticVersion) bool {
	// Even in strict equality, we ignore the build.
	if p1.Label != p2.Label {
		return false
	}
	return p1.Major == p2.Major && p1.Minor == p2.Minor && p1.Patch == p2.Patch && p1.Prerelease == p2.Prerelease
}

func (p1 SemanticVersion) Lt(p2 SemanticVersion) bool {
	if p1.Major < p2.Major {
		return true
	}
	if p1.Major > p2.Major {
		return false
	}
	if p1.Minor < p2.Minor {
		return true
	}
	if p1.Minor > p2.Minor {
		return false
	}
	return p1.Patch < p2.Patch
}

func (p1 SemanticVersion) Gt(p2 SemanticVersion) bool {
	if p1.Major > p2.Major {
		return true
	}
	if p1.Major < p2.Major {
		return false
	}
	if p1.Minor > p2.Minor {
		return true
	}
	if p1.Minor < p2.Minor {
		return false
	}
	return p1.Patch > p2.Patch
}

func (p1 SemanticVersion) Lte(p2 SemanticVersion) bool {
	return p1.Lt(p2) || p1.Eq(p2)
}

func (p1 SemanticVersion) Gte(p2 SemanticVersion) bool {
	return p1.Gt(p2) || p1.Eq(p2)
}

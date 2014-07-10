package spdx

import "strings"

type AnyLicence interface {
	Value
	LicenceId() string
}

type Licence struct{ ValueStr }

func (l Licence) LicenceId() string    { return l.V() }
func (l Licence) Equal(b Licence) bool { return l.ValueStr.Equal(b.ValueStr) }

// Returns whether the licence is a reference or a is supposed to be in the SPDX Licence List.
// Does not check if the licence actually is in the licence list (use InList() for that).
func (l Licence) IsReference() bool {
	return isLicIdRef(l.V())
}

// Checks whether the licence is in the SPDX Licence List.
// It always looks up the SPDX Licence List index.
func (l Licence) InList() bool {
	return CheckLicence(l.V())
}

// Creates a new Licence.
func NewLicence(id string, m *Meta) Licence {
	return Licence{Str(id, m)}
}

type ExtractedLicence struct {
	Id             ValueStr
	Name           []ValueStr // conditional. one required if the licence is not in the SPDX Licence List
	Text           ValueStr
	CrossReference []ValueStr //optional
	Comment        ValueStr   //optional
}

func (l *ExtractedLicence) LicenceId() string { return l.Id.V() }
func (l *ExtractedLicence) V() string         { return l.LicenceId() }
func (l *ExtractedLicence) M() *Meta          { return l.Id.M() }

// DisjunctiveLicenceSet is a AnyLicence
type ConjunctiveLicenceSet []AnyLicence

func (c ConjunctiveLicenceSet) LicenceId() string { return join(c, " and ") }
func (c ConjunctiveLicenceSet) V() string         { return c.LicenceId() }
func (c ConjunctiveLicenceSet) M() *Meta {
	for _, k := range c {
		if k != nil {
			return k.M()
		}
	}
	return nil
}

// DisjunctiveLicenceSet is a AnyLicence
type DisjunctiveLicenceSet []AnyLicence

func (c DisjunctiveLicenceSet) LicenceId() string { return join(c, " or ") }
func (c DisjunctiveLicenceSet) V() string         { return c.LicenceId() }
func (c DisjunctiveLicenceSet) M() *Meta {
	for _, k := range c {
		if k != nil {
			return k.M()
		}
	}
	return nil
}

// Useful functions for working with licences

// Join the IDs for given licences by separator. Similar
// to strings.Join but for []AnyLicence.
func join(list []AnyLicence, separator string) string {
	if len(list) == 0 {
		return "()"
	}
	res := "(" + list[0].LicenceId()
	for i := 1; i < len(list); i++ {
		res += separator + list[i].LicenceId()
	}
	res += ")"
	return res
}

// Returns whether the given ID is a Licence Reference ID (starts with LicenseRef-).
// Does not check if the string after "LicenseRef-" satisfies the requirements of any SPDX version.
// It is case-insensitive.
func isLicIdRef(id string) bool {
	return strings.HasPrefix(strings.ToLower(id), "licenseref-")
}

// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package common

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ElementID represents the identifier string portion of an SPDX element
// identifier. DocElementID should be used for any attributes which can
// contain identifiers defined in a different SPDX document.
// ElementIDs should NOT contain the mandatory 'SPDXRef-' portion.
type ElementID string

// DocElementID represents an SPDX element identifier that could be defined
// in a different SPDX document, and therefore could have a "DocumentRef-"
// portion, such as Relationships and Annotations.
// ElementID is used for attributes in which a "DocumentRef-" portion cannot
// appear, such as a Package or File definition (since it is necessarily
// being defined in the present document).
// DocumentRefID will be the empty string for elements defined in the
// present document.
// DocElementIDs should NOT contain the mandatory 'DocumentRef-' or
// 'SPDXRef-' portions.
// SpecialID is used ONLY if the DocElementID matches a defined set of
// permitted special values for a particular field, e.g. "NONE" or
// "NOASSERTION" for the right-hand side of Relationships. If SpecialID
// is set, DocumentRefID and ElementRefID should be empty (and vice versa).
type DocElementID struct {
	DocumentRefID string
	ElementRefID  ElementID
	SpecialID     string
}

// UnmarshalJSON takes a SPDX Identifier string parses it into a DocElementID struct.
// This function is also used when unmarshalling YAML
func (d *DocElementID) UnmarshalJSON(data []byte) error {
	// SPDX identifier will simply be a string
	idStr := string(data)
	idStr = strings.Trim(idStr, "\"")

	// handle special cases
	if idStr == "NONE" || idStr == "NOASSERTION" {
		d.SpecialID = idStr
		return nil
	}

	var idFields []string
	// handle DocumentRef- if present
	if strings.HasPrefix(idStr, "DocumentRef-") {
		// strip out the "DocumentRef-" so we can get the value
		idFields = strings.SplitN(idStr, "DocumentRef-", 2)
		idStr = idFields[1]

		// an SPDXRef can appear after a DocumentRef, separated by a colon
		idFields = strings.SplitN(idStr, ":", 2)
		d.DocumentRefID = idFields[0]

		if len(idFields) == 2 {
			idStr = idFields[1]
		} else {
			return nil
		}
	}

	// handle SPDXRef-
	idFields = strings.SplitN(idStr, "SPDXRef-", 2)
	if len(idFields) != 2 {
		return fmt.Errorf("failed to parse SPDX Identifier '%s'", idStr)
	}

	d.ElementRefID = ElementID(idFields[1])

	return nil
}

// MarshalJSON converts the receiver into a slice of bytes representing a DocElementID in string form.
// This function is also used when marshalling to YAML
func (d DocElementID) MarshalJSON() ([]byte, error) {
	if d.DocumentRefID != "" && d.ElementRefID != "" {
		return json.Marshal(fmt.Sprintf("DocumentRef-%s:SPDXRef-%s", d.DocumentRefID, d.ElementRefID))
	} else if d.ElementRefID != "" {
		return json.Marshal(fmt.Sprintf("SPDXRef-%s", d.ElementRefID))
	} else if d.SpecialID != "" {
		return json.Marshal(d.SpecialID)
	}

	return []byte{}, fmt.Errorf("failed to marshal empty DocElementID")
}

// TODO: add equivalents for LicenseRef- identifiers

// MakeDocElementID takes strings (without prefixes) for the DocumentRef-
// and SPDXRef- identifiers, and returns a DocElementID. An empty string
// should be used for the DocumentRef- portion if it is referring to the
// present document.
func MakeDocElementID(docRef string, eltRef string) DocElementID {
	return DocElementID{
		DocumentRefID: docRef,
		ElementRefID:  ElementID(eltRef),
	}
}

// MakeDocElementSpecial takes a "special" string (e.g. "NONE" or
// "NOASSERTION" for the right side of a Relationship), nd returns
// a DocElementID with it in the SpecialID field. Other fields will
// be empty.
func MakeDocElementSpecial(specialID string) DocElementID {
	return DocElementID{SpecialID: specialID}
}

// RenderElementID takes an ElementID and returns the string equivalent,
// with the SPDXRef- prefix reinserted.
func RenderElementID(eID ElementID) string {
	return "SPDXRef-" + string(eID)
}

// RenderDocElementID takes a DocElementID and returns the string equivalent,
// with the SPDXRef- prefix (and, if applicable, the DocumentRef- prefix)
// reinserted. If a SpecialID is present, it will be rendered verbatim and
// DocumentRefID and ElementRefID will be ignored.
func RenderDocElementID(deID DocElementID) string {
	if deID.SpecialID != "" {
		return deID.SpecialID
	}
	prefix := ""
	if deID.DocumentRefID != "" {
		prefix = "DocumentRef-" + deID.DocumentRefID + ":"
	}
	return prefix + "SPDXRef-" + string(deID.ElementRefID)
}

package i18nlevel

import (
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/responses"
)

// The COMPARATORS response. See RFC 5255 Section 4.8.
type Comparators struct {
	Active  string
	Matched []string
}

const responseName = "COMPARATOR"

/*
	ABNF:

	comparator-data   = "COMPARATOR" SP comp-sel-quoted [SP "("
	                    comp-id-quoted *(SP comp-id-quoted) ")"]
	comp-id-quoted    = astring
*/

func (rs *Comparators) Parse(fields []interface{}) error {
	var available []string
	switch len(fields) {
	case 2:
		cmps, ok := fields[1].([]interface{})
		if !ok {
			return errors.New("Second argment must be a list")
		}
		for _, cmp := range cmps {
			s, ok := cmp.(string)
			if !ok {
				return errors.New("Comparator ID must be a string")
			}

			available = append(available, s)
		}

		fallthrough
	case 1:
		cmp, ok := fields[0].(string)
		if !ok {
			return errors.New("Comparator ID must be a string")
		}

		rs.Active = cmp
		rs.Matched = available
	case 0:
		return errors.New("At least one argument is required for the COMPARATOR command")
	default:
		return errors.New("Too many arguments for the COMPARATOR command")
	}

	return nil
}

func (rs *Comparators) Format() (fields []interface{}) {
	fields = append(fields, rs.Active)

	availableF := make([]interface{}, 0, len(rs.Matched))
	for _, cmp := range rs.Matched {
		availableF = append(availableF, cmp)
	}
	fields = append(fields, availableF)
	return fields
}

func (rs *Comparators) WriteTo(w *imap.Writer) error {
	fields := []interface{}{imap.RawString(responseName)}
	fields = append(fields, rs.Format()...)
	return imap.NewUntaggedResp(fields).WriteTo(w)
}

func (rs *Comparators) Handle(resp imap.Resp) error {
	name, fields, ok := imap.ParseNamedResp(resp)
	if !ok || name != responseName {
		return responses.ErrUnhandled
	}
	return rs.Parse(fields)
}

package i18nlevel

import (
	"errors"

	"github.com/emersion/go-imap"
)

const (
	comparatorCmdName = "COMPARATOR"
)

// The COMPARATOR command. See RFC 5255 Section 4.7.
type ComparatorCmd struct {
	Comparators []string
}

func (cmd *ComparatorCmd) Command() *imap.Command {
	var args []interface{}
	for _, cmp := range cmd.Comparators {
		args = append(args, cmp)
	}

	return &imap.Command{
		Name:      comparatorCmdName,
		Arguments: args,
	}
}

func (cmd *ComparatorCmd) Parse(fields []interface{}) error {
	var cmps []string
	for _, arg := range fields {
		s, ok := arg.(string)
		if !ok {
			return errors.New("COMPARATOR argument must be a string")
		}

		cmps = append(cmps, s)
	}
	cmd.Comparators = cmps
	return nil
}

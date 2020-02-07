package i18nlevel

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/server"
)

var (
	ErrUnsupportedBackend = errors.New("i18nlevel: backend not supported")
)

type Backend interface {
	// I18NLevel returns the internationalization level supported by backend.
	//
	// See RFC 5255 4.2, 4.3 and 4.6 sections for requirements of each level.
	//
	// If this method returns 1, the User objects do not have to implement the
	// User interface.
	I18NLevel() int
}

type User interface {
	// UseComparator changes active comparator for the session.
	//
	// cmps is a non-empty list of comparators client wishes to use in the
	// preference order. The backend selects the first supported comparator and
	// changes the active session comparator to it. The return value contains
	// the identifier of selected comparator. And a list of  comparators which
	// matches any of the arguments to the COMPARATOR command and is present
	// only if more than one match is found.
	//
	// Special value "default" refers to the comparator that is used when no
	// UseComparator is called. Note that with I18NLevel() == 2, the default
	// comparator SHOULD be "i;unicode-casemap".
	UseComparator(cmps []string) (string, []string, error)

	// ActiveComparator returns the currently active comparator for the
	// session.
	ActiveComparator() string
}

type ComparatorHandler struct {
	ComparatorCmd
}

func (h *ComparatorHandler) Handle(conn server.Conn) error {
	if conn.Context().User == nil {
		return server.ErrNotAuthenticated
	}

	be, ok := conn.Server().Backend.(Backend)
	if !ok {
		// Backend is not compatible with extension, no-op.
		return ErrUnsupportedBackend
	}

	if be.I18NLevel() <= 1 {
		return errors.New("COMPARATOR command is not supported")
	}

	u, ok := conn.Context().User.(User)
	if !ok {
		return ErrUnsupportedBackend
	}

	// No change requested, return active comparator.
	if len(h.Comparators) == 0 {
		return conn.WriteResp(&Comparators{
			Active: u.ActiveComparator(),
		})
	}

	// Comparator change.
	active, matched, err := u.UseComparator(h.Comparators)
	if err != nil {
		return server.ErrStatusResp(&imap.StatusResp{
			Type: imap.StatusRespNo,
			Code: "BADCOMPARATOR",
			Info: "Failed to activate any of selected comparators",
		})
	}

	return conn.WriteResp(&Comparators{
		Active:  active,
		Matched: matched,
	})
}

type ext struct{}

func NewExtension() server.Extension {
	return &ext{}
}

func (ext *ext) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState == 0 {
		return nil
	}

	be, ok := c.Server().Backend.(Backend)
	if !ok {
		// Backend is not compatible with extension, no-op.
		return nil
	}

	switch be.I18NLevel() {
	case 0:
		return nil
	case 1:
		return []string{"I18NLEVEL=1"}
	case 2:
		return []string{"I18NLEVEL=2"}
	default:
		return []string{"I18NLEVEL=2"}
	}
}

func (ext *ext) Command(name string) server.HandlerFactory {
	if name == comparatorCmdName {
		return func() server.Handler {
			return &ComparatorHandler{}
		}
	}
	return nil
}

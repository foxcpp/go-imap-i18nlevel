package i18nlevel

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	c *client.Client
}

func NewClient(c *client.Client) *Client {
	return &Client{c: c}
}

// I18NLevel returns the internationalization level supported by the server.
//
// If server does not support the I18NLEVEL extension, 0 is returned.
func (c *Client) I18NLevel() (int, error) {
	ok, err := c.c.Support("I18NLEVEL=1")
	if err != nil {
		return -1, nil
	}
	if ok {
		return 1, nil
	}

	ok, err = c.c.Support("I18NLEVEL=2")
	if err != nil {
		return -1, nil
	}
	if ok {
		return 2, nil
	}

	return 0, nil
}

// ActiveComparator returns the active comparator used.
//
// This command is valid only if I18NLevel() returns 2.
// See RFC 5255 for details.
func (c *Client) ActiveComparator() (string, error) {
	if c.c.State()&imap.AuthenticatedState == 0 {
		return "", client.ErrNotLoggedIn
	}

	res := &Comparators{}
	status, err := c.c.Execute(&ComparatorCmd{}, res)
	if err != nil {
		return "", err
	}
	if err := status.Err(); err != nil {
		return "", err
	}
	return res.Active, nil
}

// ActiveComparator changes the active comparator to the first comparator
// listed in cmps and supported by the server.
//
// This command is valid only if I18NLevel() returns 2.
// See RFC 5255 for details.
func (c *Client) UseComparator(cmps []string) (string, []string, error) {
	if c.c.State()&imap.AuthenticatedState == 0 {
		return "", nil, client.ErrNotLoggedIn
	}

	res := &Comparators{}
	status, err := c.c.Execute(&ComparatorCmd{
		Comparators: cmps,
	}, res)
	if err != nil {
		return "", nil, err
	}
	if err := status.Err(); err != nil {
		return "", nil, err
	}
	return res.Active, res.Matched, nil
}

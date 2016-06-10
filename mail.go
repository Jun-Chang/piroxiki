package piroxiki

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
	"github.com/mxk/go-imap/imap"
	"github.com/pkg/errors"
)

const (
	MailRfcHeader = "RFC822.HEADER"
	MailRfcText   = "RFC822.TEXT"
)

type Mail struct {
	Server   string
	User     string
	Password string
	MailBox  string
	filter   []imap.Field
	message  chan Message
	errc     chan<- error
}

func (m *Mail) HandleInput(policy Policy, filter Filter, errc chan<- error) (<-chan Message, error) {
	m.errc = errc
	m.message = make(chan Message)

	if m.MailBox == "" {
		m.MailBox = "inbox"
	}

	if err := m.setFilter(filter); err != nil {
		return nil, err
	}

	go func() {
		tick := time.Tick(time.Second * time.Duration(policy.IntervalSeconds))
		errCount := 0
		for {
			<-tick
			err := m.input()
			if err != nil {
				errCount++
				if errCount >= policy.Count {
					errc <- err
				}
			} else {
				errCount = 0
			}
		}
	}()

	return m.message, nil
}

func (*Mail) HandleOutput(Policy, <-chan Message, chan<- error) error {
	// TODO:
	return nil
}

func (m *Mail) input() error {
	fmt.Println("mail input start")

	// dial
	c, err := imap.DialTLS(m.Server, nil)
	if err != nil {
		return err
	}
	defer c.Close(false)

	if _, err := c.Login(m.User, m.Password); err != nil {
		return err
	}
	if _, err := c.Select(m.MailBox, false); err != nil {
		return err
	}

	spec := append(m.filter, "UNSEEN")
	cmd, err := c.Search(spec...)
	if err != nil {
		return err
	}
	if _, err := cmd.Result(imap.OK); err != nil {
		return err
	}

	res := strings.Split(cmd.Data[0].String(), " ")
	if len(res) <= 2 {
		return nil
	}

	for _, r := range res[2:] {
		fmt.Println("id", r)
		seq, err := imap.NewSeqSet(r)
		if err != nil {
			return err
		}
		cmd, err = c.Fetch(seq, MailRfcHeader, MailRfcText)
		if err != nil {
			return err
		}
		if _, err := cmd.Result(imap.OK); err != nil {
			return err
		}

		var header, body string
		for _, f := range cmd.Data[0].Fields {
			imf, ok := f.([]imap.Field)
			if ok {
				for i, im := range imf {
					lt, ok := im.(imap.Literal)
					if ok {
						nm, ok := imf[i-1].(string)
						if !ok {
							return err
						}
						buf := &bytes.Buffer{}
						lt.WriteTo(buf)
						switch nm {
						case MailRfcHeader:
							header = buf.String()
						case MailRfcText:
							body, _ = html2text.FromReader(buf)
						default:
							return err
						}
					}
				}
			}
		}
		if header != "" && body != "" {
			cmd, err := c.UIDStore(seq, "+FLAGS.SILENT", imap.NewFlagSet(`\SEEN`))
			if err != nil {
				return err
			}
			if _, err := cmd.Result(imap.OK); err != nil {
				return err
			}

			m.message <- m.format(header, body)
		}
	}

	return nil
}

func (m *Mail) setFilter(filter Filter) error {
	m.filter = make([]imap.Field, 0, 3)

	set := func(nm string) error {
		if from, ok := filter[nm]; ok {
			s, ok := from.(string)
			if !ok {
				return errors.New(nm + " is not string")
			}
			m.filter = append(m.filter, []imap.Field{nm, s})
		}
		return nil
	}

	// from
	if err := set("from"); err != nil {
		return err
	}
	// subject
	if err := set("subject"); err != nil {
		return err
	}
	// body
	if err := set("body"); err != nil {
		return err
	}

	fmt.Println(m.filter)

	return nil
}

func (m *Mail) format(header string, body string) map[string]interface{} {
	//TODO: encode
	return map[string]interface{}{
		"body": fmt.Sprintf("HEADER\n\n%s\nBODY\n\n%s", "TODO", body),
	}
}

package sendmail

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"testing"
)

const domain = "example.com"

func maddr(name, address string) *mail.Address {
	return &mail.Address{Name: name, Address: address + domain}
}

func TestSend(tc *testing.T) {
	tc.Run("debug:true", func(t *testing.T) {
		t.Parallel()
		testSend(t, true)
	})
	tc.Run("debug:false", func(t *testing.T) {
		t.Parallel()
		testSend(t, false)
	})
}

func testSend(t *testing.T, withDebug bool) {
	sm := Mail{
		Subject: "Cześć",
		From:    maddr("Michał", "me@"),
		To: []*mail.Address{
			maddr("Ktoś", "info@"),
			maddr("Ktoś2", "info2@"),
		},
	}
	sm.SetSendmail("/bin/true").SetDebug(withDebug)

	io.WriteString(&sm.Text, ":)\r\n")
	if err := sm.Send(); err != nil {
		t.Errorf("(debug=%v) %v", withDebug, err)
	}
	subject := sm.Header.Get("Subject")
	if subject != "=?utf-8?q?Cze=C5=9B=C4=87?=" {
		t.Errorf("(debug=%v) Wrong `Subject` encoding: %s", withDebug, subject)
	}
	from := sm.Header.Get("From")
	if from != fmt.Sprintf("=?utf-8?q?Micha=C5=82?= <me@%s>", domain) {
		t.Errorf("(debug=%v) Wrong `From` encoding: %s", withDebug, from)
	}
}

func TestFromError(t *testing.T) {
	sm := Mail{
		To: []*mail.Address{maddr("Ktoś", "info@")},
	}
	if sm.Send() == nil {
		t.Errorf("Expected an error because of missing `From` addresses")
	}
}

func TestToError(t *testing.T) {
	sm := Mail{
		From: maddr("Michał", "me@"),
	}
	if sm.Send() == nil {
		t.Errorf("Expected an error because of missing `To` addresses")
	}
}

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	m := New(
		Subject("Test subject"),
		From("Dominik", "dominik@example.org"),
		To("Dominik2", "dominik2@example.org"),
		DebugOutput(&buf),
		Sendmail("/bin/true"),
	)

	if m.Subject != "Test subject" {
		t.Errorf("Expected subject to be %q, got %q", "Test subject", m.Subject)
	}
	if len(m.To) != 1 {
		t.Errorf("Expected len(To) to be 1, got %d: %+v", len(m.To), m.To)
	}
	if m.From == nil || m.From.Address != "dominik@example.org" {
		expected := mail.Address{Name: "Dominik", Address: "dominik@example.org"}
		t.Errorf("Expected From address to be %s, got %s", expected, m.From)
	}
	if m.sendmail != "/bin/true" {
		t.Errorf("Expected sendmail to be %q, got %q", "/bin/true", m.sendmail)
	}
	if m.debugOut != &buf {
		t.Errorf("Expected debugOut to be %T (buf), got %T", &buf, m.debugOut)
	}
}

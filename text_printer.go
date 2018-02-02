package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type EntryTextPrinter struct {
	mu  sync.Mutex
	buf []byte
}

func (w *EntryTextPrinter) Print(entry *Entry, wr io.Writer) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	buf := w.buf[0:0]
	buf = append(buf, '[')
	buf = append(buf, entry.Level.String()...)
	buf = append(buf, ']')
	buf = append(buf, '\t')
	w.writeTime(&buf, entry.Time, entry.Flags)
	if len(entry.File) > 0 {
		buf = append(buf, entry.File...)
		if len(entry.Function) > 0 {
			buf = append(buf, '(')
			buf = append(buf, entry.Function...)
			buf = append(buf, ')')
		}
	} else if len(entry.Function) > 0 {
		buf = append(buf, entry.Function...)
	}

	if len(entry.File) > 0 || len(entry.Function) > 0 {
		buf = append(buf, ':')
		itoa(&buf, entry.Line, -1)
		buf = append(buf, '\t')
		buf = append(buf, '|')
		buf = append(buf, ' ')
	}

	for _, f := range entry.Fields {
		buf = append(buf, f.Key...)
		buf = append(buf, '=')
		buf = append(buf, fmt.Sprint(f.Value)...)
		buf = append(buf, ' ')
	}

	if len(entry.Fields) > 0 {
		buf = append(buf, '\t')
		buf = append(buf, '|')
		buf = append(buf, ' ')
	}

	buf = append(buf, entry.Message...)
	if buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}

	n, err := wr.Write(buf)
	for n < len(buf) && err == nil {
		buf = buf[n:]
		n, err = wr.Write(buf)
	}

	return err
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (w *EntryTextPrinter) writeTime(buf *[]byte, t time.Time, flags int) {
	if t.Unix() == 0 || flags&(Ldate|Ltime|Lmicroseconds) == 0 {
		return
	}

	if flags&Ldate != 0 {
		year, month, day := t.Date()
		itoa(buf, year, 4)
		*buf = append(*buf, '/')
		itoa(buf, int(month), 2)
		*buf = append(*buf, '/')
		itoa(buf, day, 2)
		*buf = append(*buf, ' ')
	}

	if flags&(Ltime|Lmicroseconds) != 0 {
		hour, min, sec := t.Clock()
		itoa(buf, hour, 2)
		*buf = append(*buf, ':')
		itoa(buf, min, 2)
		*buf = append(*buf, ':')
		itoa(buf, sec, 2)
		if flags&Lmicroseconds != 0 {
			*buf = append(*buf, '.')
			itoa(buf, t.Nanosecond()/1e3, 6)
		}
		*buf = append(*buf, ' ')
	}
}

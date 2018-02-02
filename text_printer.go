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
	w.buf = w.buf[0:0]
	w.buf = append(w.buf, '[')
	w.buf = append(w.buf, entry.Level.String()...)
	w.buf = append(w.buf, ']')
	w.buf = append(w.buf, '\t')
	w.writeTime(entry.Time, entry.Flags)
	if len(entry.File) > 0 {
		w.buf = append(w.buf, entry.File...)
		if len(entry.Function) > 0 {
			w.buf = append(w.buf, '(')
			w.buf = append(w.buf, entry.Function...)
			w.buf = append(w.buf, ')')
		}
	} else if len(entry.Function) > 0 {
		w.buf = append(w.buf, entry.Function...)
	}

	if len(entry.File) > 0 || len(entry.Function) > 0 {
		w.buf = append(w.buf, ':')
		itoa(&w.buf, entry.Line, -1)
		w.buf = append(w.buf, '\t')
		w.buf = append(w.buf, '|')
		w.buf = append(w.buf, ' ')
	}

	for _, f := range entry.Fields {
		w.buf = append(w.buf, f.Key...)
		w.buf = append(w.buf, '=')
		w.buf = append(w.buf, fmt.Sprint(f.Value)...)
		w.buf = append(w.buf, ' ')
	}

	if len(entry.Fields) > 0 {
		w.buf = append(w.buf, '\t')
		w.buf = append(w.buf, '|')
		w.buf = append(w.buf, ' ')
	}

	w.buf = append(w.buf, entry.Message...)
	if w.buf[len(w.buf)-1] != '\n' {
		w.buf = append(w.buf, '\n')
	}

	n, err := wr.Write(w.buf)
	for n < len(w.buf) && err == nil {
		w.buf = w.buf[n:]
		n, err = wr.Write(w.buf)
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

func (w *EntryTextPrinter) writeTime(t time.Time, flags int) {
	if t.Unix() == 0 || flags&(Ldate|Ltime|Lmicroseconds) == 0 {
		return
	}

	if flags&Ldate != 0 {
		year, month, day := t.Date()
		itoa(&w.buf, year, 4)
		w.buf = append(w.buf, '/')
		itoa(&w.buf, int(month), 2)
		w.buf = append(w.buf, '/')
		itoa(&w.buf, day, 2)
		w.buf = append(w.buf, ' ')
	}

	if flags&(Ltime|Lmicroseconds) != 0 {
		hour, min, sec := t.Clock()
		itoa(&w.buf, hour, 2)
		w.buf = append(w.buf, ':')
		itoa(&w.buf, min, 2)
		w.buf = append(w.buf, ':')
		itoa(&w.buf, sec, 2)
		if flags&Lmicroseconds != 0 {
			w.buf = append(w.buf, '.')
			itoa(&w.buf, t.Nanosecond()/1e3, 6)
		}
		w.buf = append(w.buf, ' ')
	}
}

package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

const logLevelWidth = 8

//var (
//	greenColor   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
//	whiteColor   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
//	yellowColor  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
//	redColor     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
//	blueColor    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
//	magentaColor = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
//	cyanColor    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
//	resetColor   = string([]byte{27, 91, 48, 109})
//)

type EntryTextPrinter struct {
	mu  sync.Mutex
	buf []byte
}

func (w *EntryTextPrinter) Print(entry *Entry, wr io.Writer) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	buf := w.buf[0:0]
	w.writeTime(&buf, entry.Time, entry.Flags)

	buf = append(buf, '[')
	buf = append(buf, entry.Level.String()...)
	buf = append(buf, ']')
	for i := len(entry.Level.String()) + 2; i < logLevelWidth; i++ {
		buf = append(buf, ' ')
	}

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
		buf = append(buf, fmt.Sprintf("%v", f.Value)...)
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
	if t.Unix() == 0 || flags&(Ldate|Ltime|Lmillisecond|Lmicroseconds) == 0 {
		return
	}

	if flags&Ldate != 0 {
		year, month, day := t.Date()
		itoa(buf, year, 4)
		*buf = append(*buf, '-')
		itoa(buf, int(month), 2)
		*buf = append(*buf, '-')
		itoa(buf, day, 2)
		*buf = append(*buf, ' ')
	}

	if flags&(Ltime|Lmillisecond|Lmicroseconds) != 0 {
		hour, min, sec := t.Clock()
		itoa(buf, hour, 2)
		*buf = append(*buf, ':')
		itoa(buf, min, 2)
		*buf = append(*buf, ':')
		itoa(buf, sec, 2)
		if flags&Lmicroseconds != 0 {
			*buf = append(*buf, '.')
			itoa(buf, t.Nanosecond()/1e3, 6)
		} else if flags&Lmillisecond != 0 {
			*buf = append(*buf, '.')
			itoa(buf, t.Nanosecond()/1e6, 3)
		}
		_, offset := t.Zone()
		// e.g. UTC+0800, offset is 28800 seconds, +0800 = offset/3600*100 = offset/36
		offset = offset / 36
		*buf = append(*buf, '+')
		itoa(buf, offset, 4)
		*buf = append(*buf, ' ')
	}
}

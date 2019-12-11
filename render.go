package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type render struct {
	wr  io.Writer
	mu  sync.Mutex
	buf []byte
}

func newRender(wr io.Writer) *render {
	return &render{
		wr:  wr,
		buf: make([]byte, 0, 2048), // 2048 bytes should be enough for most Log entry
	}
}

func (r *render) SetWriter(wr io.Writer) {
	r.mu.Lock()
	r.wr = wr
	r.mu.Unlock()
}

func (r *render) Render(e *entry) error {
	r.mu.Lock()

	r.buf = r.buf[0:0]
	renderEntry(&r.buf, e)

	// flush buffer to writer
	buf := r.buf
	n, err := r.wr.Write(buf)
	for n < len(buf) && err == nil {
		buf = buf[n:]
		n, err = r.wr.Write(buf)
	}

	r.mu.Unlock()
	return err
}

// RenderString is only called by Log.Panic[f], it's ok to use local buffer
func (r *render) RenderString(e *entry) string {
	r.mu.Lock()
	r.buf = r.buf[0:0]
	renderEntry(&r.buf, e)
	str := string(r.buf)
	r.mu.Unlock()
	return str
}

func renderEntry(buf *[]byte, e *entry) {
	writeTime(buf, e.Time, e.Flags)

	*buf = append(*buf, '[')
	*buf = append(*buf, e.Level.String()...)
	*buf = append(*buf, ']', ' ')

	if len(e.Name) > 0 {
		*buf = append(*buf, '[')
		*buf = append(*buf, e.Name...)
		*buf = append(*buf, ']', ' ')
	}

	if len(e.File) > 0 {
		*buf = append(*buf, e.File...)
		if len(e.Function) > 0 {
			*buf = append(*buf, '(')
			*buf = append(*buf, e.Function...)
			*buf = append(*buf, ')')
		}
	} else if len(e.Function) > 0 {
		*buf = append(*buf, e.Function...)
	}

	if len(e.File) > 0 || len(e.Function) > 0 {
		*buf = append(*buf, ':')
		itoa(&*buf, e.Line, -1)
		*buf = append(*buf, '\t')
		*buf = append(*buf, '|')
		*buf = append(*buf, ' ')
	}

	for _, f := range e.Fields {
		*buf = append(*buf, f.Key...)
		*buf = append(*buf, '=')
		*buf = append(*buf, fmt.Sprintf("%+v", f.Value)...)
		*buf = append(*buf, ' ')
	}

	if len(e.Fields) > 0 {
		*buf = append(*buf, '\t')
		*buf = append(*buf, '|')
		*buf = append(*buf, ' ')
	}

	*buf = append(*buf, e.Message...)
	if (*buf)[len(*buf)-1] != '\n' {
		*buf = append(*buf, '\n')
	}
}

func writeTime(buf *[]byte, t time.Time, flags int) {
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

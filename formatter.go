package log

type Formatter interface {
	Format(entry *Entry) string
}

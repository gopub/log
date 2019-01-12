package log

const (
	Ldate         = 1 << iota                                             // the date in the local time zone: 2009-01-23
	Ltime                                                                 // the time in the local time zone: 01:23:23
	Lmillisecond                                                          // microsecond resolution: 01:23:23.123.  assumes Ltime.
	Lmicroseconds                                                         // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                                             // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                                            // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                                                  // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lfunction                                                             // function name and line number: print:23. overrides Llongfile, Lshortfile
	Lname                                                                 // logger's name
	LstdFlags     = Ldate | Lmillisecond | Lshortfile | Lfunction | Lname // initial values for the standard logger
)

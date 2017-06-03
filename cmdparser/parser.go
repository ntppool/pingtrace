package cmdparser

type Parser interface {
	Add(string)
	Close()
	Read() ParserOutput
}

type ParserOutput interface {
	Bytes() []byte
	JSON() []byte
	String() string
	Error() error
}

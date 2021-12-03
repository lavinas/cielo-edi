package ports

type StringParser interface {
	Parse(interface{}, string) error
}

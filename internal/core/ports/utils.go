package ports

type StringParser interface {
	Marshal(interface{}, string)
}

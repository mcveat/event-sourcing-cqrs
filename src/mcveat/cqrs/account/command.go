package account

type Command interface{}

type OpenAccount struct {
	InitialBalance int
}

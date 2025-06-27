package commands

const (
	GET_COMMAND    = 0
	SET_COMMAND    = 1
	DELETE_COMMAND = 2
)

type Command struct {
	Identifier int
	Key        string
	Value      string
}

func CreateGetCommand(key string) Command {
	return Command{GET_COMMAND, key, ""}
}

func CreateSetCommand(key, value string) Command {
	return Command{SET_COMMAND, key, value}
}

func CreateDeleteCommand(key string) Command {
	return Command{DELETE_COMMAND, key, ""}
}

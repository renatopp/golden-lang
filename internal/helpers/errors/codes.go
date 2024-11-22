package errors

const (
	InternalError ErrorCode = iota
	NotImplemented
	InvalidFileError
	InvalidFolderError
	CircularReferenceError
	ParserError
	TypeError
	NameNotFound
	NameAlreadyDefined
)

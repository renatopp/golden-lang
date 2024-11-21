package errors

const (
	InternalError ErrorCode = iota
	InvalidFileError
	InvalidFolderError
	ParserError
	CircularReferenceError
)

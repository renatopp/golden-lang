package naming

import "regexp"

var typeRegex = regexp.MustCompile(`^_*[A-Z][a-zA-Z0-9_]*$`)
var valueRegex = regexp.MustCompile(`^_*[a-z][a-zA-Z0-9_]*$`)
var publicRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
var privateRegex = regexp.MustCompile(`^_[a-zA-Z0-9_]*$`)
var wildcardRegex = regexp.MustCompile(`^_$`)

func IsTypeName(s string) bool    { return typeRegex.MatchString(s) }
func IsValueName(s string) bool   { return valueRegex.MatchString(s) }
func IsPublicName(s string) bool  { return publicRegex.MatchString(s) }
func IsPrivateName(s string) bool { return privateRegex.MatchString(s) }
func IsWildcard(s string) bool    { return wildcardRegex.MatchString(s) }

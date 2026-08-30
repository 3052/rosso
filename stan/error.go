package stan

import "strings"

type ApiError struct {
   Code string
}

type ApiErrors []ApiError

func (e ApiErrors) Error() string {
   var builder strings.Builder
   builder.WriteString("stan: ")
   for i, err := range e {
      if i > 0 {
         builder.WriteString(", ")
      }
      builder.WriteString(err.Code)
   }
   return builder.String()
}

// error.go

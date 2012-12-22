package godata

import (
	"fmt"
)

type ModuleExistsError string

func (err ModuleExistsError) Error() string {
	return fmt.Sprintf("Module %v has already been registered", string(err))
}

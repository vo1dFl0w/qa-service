package usecase

import "fmt"

func validateID(id int) error {
	if id < 0 {
		return fmt.Errorf("error id < 0")
	}

	return nil
}
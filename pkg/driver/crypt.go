package driver

import "fmt"

func ClevisBind(path, server string) error {
	template := `{"url":"http://%s"}`

	tang := fmt.Sprintf(template, server)

	cmd := "clevis"
	arg := []string{"luks", "bind", "-d", path, "tang", tang}

	err := run(cmd, arg)
	if err != nil {
		return fmt.Errorf("problem Format: %v", err.Error())
	}

	return nil
}

func CheckClevisBind(path string) error {
	cmd := "clevis"
	arg := []string{"luks", "list", "-d", path}
	err := run(cmd, arg)
	if err != nil {
		return fmt.Errorf("clevis not binding: %v", err.Error())
	}

	return nil
}

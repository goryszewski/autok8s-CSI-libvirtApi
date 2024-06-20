package driver

import (
	"fmt"
	"os"
	"os/exec"
)

func mount(source, target string, args *[]string) error {
	err := os.MkdirAll(target, 0777)
	if err != nil {
		return fmt.Errorf("MKDIR problem create")
	}

	cmd := "mount"
	arg := []string{}
	arg = append(arg, *args...)
	arg = append(arg, source)
	arg = append(arg, target)

	_, err = exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("[DEBUG][mount][LookPath][ERROR] %#+v \n", err)
		return fmt.Errorf("problem path cmd Mount: %v", err.Error())
	}

	cmdmount := exec.Command(cmd, arg...)

	_, err = cmdmount.Output()
	if err != nil {

		return fmt.Errorf("problem Mount: %v", err.Error())
	}

	return nil
}

func Umount(path string) error {
	cmd := "umount"
	arg := append([]string{}, path)
	cmdmount := exec.Command(cmd, arg...)

	_, err := cmdmount.Output()
	if err != nil {

		return fmt.Errorf("problem Mount: %v", err.Error())
	}

	return nil
}

func Formater(fstype, source string) error {
	cmd := fmt.Sprintf("mkfs.%s", fstype)
	arg := []string{"-F", source}

	_, err := exec.LookPath(cmd)
	if err != nil {
		return fmt.Errorf("problem  path cmd Format: %v", err.Error())
	}

	_, err = exec.Command(cmd, arg...).CombinedOutput()
	if err != nil {

		return fmt.Errorf("problem Format: %v", err.Error())
	}

	return nil
}
func isNotFormated(source string) (bool, error) {
	cmd := "dumpe2fs"
	arg := []string{source}

	err := run(cmd, arg)
	if err != nil {
		return true, fmt.Errorf("error run cmd")
	}

	return false, nil
}

func run(cmd string, arg []string) error {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return fmt.Errorf("problem path cmd Mount: %v", err.Error())
	}

	cmdmount := exec.Command(cmd, arg...)

	_, err = cmdmount.Output()
	if err != nil {

		return fmt.Errorf("problem %s: %v", cmd, err.Error())
	}

	return nil
}

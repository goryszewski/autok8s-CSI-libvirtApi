package driver

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func mount(source, target string, args *[]string) error {
	err := os.MkdirAll(target, 0777)
	if err != nil {
		return fmt.Errorf("mkdir: %v problem create: %v", target, err.Error())
	}

	cmd := "mount"
	arg := []string{}
	arg = append(arg, *args...)
	arg = append(arg, source)
	arg = append(arg, target)

	err = run(cmd, arg)
	if err != nil {
		return fmt.Errorf("problem Mount: %v", err.Error())
	}

	return nil
}

func Umount(path string) error {
	cmd := "umount"
	arg := append([]string{}, path)

	err := run(cmd, arg)
	if err != nil {
		return fmt.Errorf("problem Mount: %v", err.Error())
	}

	return nil
}

func Formater(fstype, source string) error {
	cmd := fmt.Sprintf("mkfs.%s", fstype)
	arg := []string{"-F", source}

	err := run(cmd, arg)
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
		return true, fmt.Errorf("error exec: %s", err.Error())
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

func GetIDNode() (string, error) {
	// DOTO endpoint metadata
	nodeID, err := os.ReadFile("/id")
	if err != nil {
		return "", fmt.Errorf("failed read is %s \n", err.Error())
	}
	id := strings.ReplaceAll(string(nodeID), "\\n", "")
	id = strings.ReplaceAll(id, "\n", "")
	return id, nil
}

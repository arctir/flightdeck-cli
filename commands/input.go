package commands

import (
	"os"
	"os/exec"
)

const defaultEditor = "vim"

func openFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func captureInputFromEditor(content []byte) ([]byte, error) {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	if content != nil {
		_, err := file.Write(content)
		if err != nil {
			return nil, err
		}
	}

	filename := file.Name()
	defer os.Remove(filename)
	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = openFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

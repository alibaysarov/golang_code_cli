package files

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("ошибка создания директории:", err)
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	l, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return err
	}
	fmt.Println(l, "Изменения применены!")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

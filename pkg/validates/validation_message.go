package validates

import "fmt"

func GetValidationMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("Kolom %s wajib diisi", field)
	case "email":
		return fmt.Sprintf("Kolom %s harus berupa email yang valid", field)
	case "min":
		return fmt.Sprintf("Kolom %s minimal %s karakter", field, param)
	case "max":
		return fmt.Sprintf("Kolom %s maksimal %s karakter", field, param)
	case "len":
		return fmt.Sprintf("Kolom %s harus terdiri dari %s karakter", field, param)
	case "numeric":
		return fmt.Sprintf("Kolom %s harus berupa angka", field)
	default:
		return fmt.Sprintf("Kolom %s tidak valid", field)
	}
}

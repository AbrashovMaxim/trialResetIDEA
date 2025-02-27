package main

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Функция для поиска папки, начинающейся на определенное название, в указанной директории
func findFolder(directory string, name string) (string, error) {
	var result string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		// Проверяем, что это директория и имя папки начинается с названия
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasPrefix(info.Name(), name) {
			result = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", fmt.Errorf("Папка, начинающаяся на " + name + " - не найдена!")
	}
	return result, nil
}

// Удаляем ключ в регистре
func deleteKey(key registry.Key, path string) error {
	parentKey, err := registry.OpenKey(key, path, registry.READ)
	if errors.Is(err, registry.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	defer parentKey.Close()

	if len(path) > 0 && path[len(path)-1] != '\\' {
		path = path + `\`
	}

	subKeyNames, err := parentKey.ReadSubKeyNames(-1)
	if err != nil {
		return err
	}

	for _, name := range subKeyNames {
		subKeyPath := path + name
		err := deleteKey(key, subKeyPath)
		if err != nil {
			return err
		}
	}

	return registry.DeleteKey(key, path)
}

func main() {
	names := map[int]string{1: "WebStorm", 2: "IntelliJ", 3: "GoLand", 4: "PyCharm", 5: "PhpStorm", 6: "CLion", 7: "Rider", 8: "RustRover", 9: "DataGrip", 10: "RubyMine"}
	keyNames := map[int]string{1: "webstorm", 2: "idea", 3: "goland", 4: "pycharm", 5: "phpstorm", 6: "clion", 7: "rider", 8: "rustrover", 9: "datagrip", 10: "rubymine"}
	idFilesNames := [4]string{"bl", "crl", "PermanentDeviceId", "PermanentUserId"}

	fmt.Println(
		"=======================================\n" +
			"\tСброс JetBrains IDEA v0.2\n" +
			"=======================================\n",
	)

	// Сортируем ключи, потому что, иногда, они идут не по порядку
	keys := make([]int, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Выводим информацию по ключам в консоль
	printString := "Выбери приложения, которые ты хочешь сбросить:\n"
	for _, k := range keys {
		value := strconv.Itoa(k)
		printString += "\t[" + value + "] " + names[k] + "\n"
	}
	printString += "\n!!! ВАЖНО !!!\nВводи номера через запятую, например: \"1,2,3\"\n"

	// Получаем строку из чисел (возможно)
	fmt.Println(printString)
	fmt.Print("Введите значения: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	selectedReset := scanner.Text()
	fmt.Println("")

	// Удаляем лишние пробелы и делим строку по запятым
	selectedReset = strings.ReplaceAll(selectedReset, " ", "")
	splitValues := strings.Split(selectedReset, ",")

	// Получаем путь до AppData\Roaming
	path, _ := os.UserConfigDir()

	// Проходимся циклом по веденным индексам
	for _, element := range splitValues {
		i, err := strconv.Atoi(element)

		if err != nil {
			continue
		}
		valName, _ := names[i]
		valKey, ok := keyNames[i]

		// Если такой индекс существует, то действуем
		if ok {
			// Проверяем, существует ли папка
			result, err := findFolder(path+"\\JetBrains", valName)
			if err != nil {
				fmt.Println("[K] Папка " + valName + " - не найдена!")
				continue
			}

			// Удаляем ключ в папке
			err = os.Remove(result + "\\" + valKey + ".key")
			if err != nil {
				fmt.Println("[K] Ключ " + valName + " - не найден!")
				continue
			}

			fmt.Println("[K] Ключ " + valName + " - удален!")
		} else {
			fmt.Println("[K] Значение: \"", element, "\" - не найдено!")
		}
	}

	// Удаляем файлы с ID
	for _, element := range idFilesNames {
		err := os.Remove(path + "\\JetBrains\\" + element)
		if err != nil {
			fmt.Println("[I] Файл ID - " + element + " - не найден! ")
		} else {
			fmt.Println("[I] Файл ID - " + element + " - удален! ")
		}
	}

	// Путь к ключу реестра, который нужно удалить
	keyPath := `Software\JavaSoft`

	// Открываем ключ реестра HKEY_CURRENT_USER
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("[R] Ключи регистра - не найдены!")
	} else {
		defer key.Close()

		// Удаляем ключ
		err = deleteKey(registry.CURRENT_USER, keyPath)
		if err != nil {
			fmt.Printf("[R] Ошибка при удалении ключа реестра: %v\n", err)
		} else {
			fmt.Println("[R] Ключи регистра - удалены!")
		}

	}

	// Ждем нажатие ENTER
	fmt.Print("\nPress 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}

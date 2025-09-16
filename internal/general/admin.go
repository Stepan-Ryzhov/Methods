package general

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

	models "methodi_razrabotki/internal/models"
	rep "methodi_razrabotki/internal/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateStoreMan(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	AdminSidebar(user, app, window, content)
	z := canvas.NewText("Зарегистрировать кладовщика:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите имя")

	lastnameEntry := widget.NewEntry()
	lastnameEntry.SetPlaceHolder("Введите фамилию")

	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Введите e-mail")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Придумайте пароль")

	passwordEntry2 := widget.NewPasswordEntry()
	passwordEntry2.SetPlaceHolder("Подтвердите пароль")

	loginBtn := widget.NewButton("Зарегистрировать", func() {
		if passwordEntry.Text != passwordEntry2.Text {
			dialog.NewError(errors.New("Введенные пароли не совпадают"), window).Show()
		}
		registerRequest := &models.RegisterRequest{
			FirstName: nameEntry.Text,
			LastName:  lastnameEntry.Text,
			Email:     loginEntry.Text,
			Password:  passwordEntry.Text,
		}
		err := rep.CreateStoreMan(registerRequest)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			dialog.NewInformation("Успех", "Вы успешно зарегистрировались! Воспользуйтесь меню входа для доступа к приложению.", window).Show()
		}
	})
	loginBtn.Resize(fyne.NewSize(200, 50))

	otstup := widget.NewLabel(" ")

	loginform := container.NewVBox(nameEntry, lastnameEntry, loginEntry, passwordEntry, passwordEntry2, otstup, loginBtn)
	loginform.Resize(fyne.NewSize(200, 200))
	loginform.Move(fyne.NewPos(600, 200))

	content.Add(loginform)
}

func UpdateProduct(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	AdminSidebar(user, app, window, content)
	z := canvas.NewText("Редактировать товар:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	products, err := rep.GetProducts()
	if err != nil {
		dialog.NewError(err, window).Show()
	} else {
		getSize := func() (int, int) {
			return len(products) + 1, 5
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if cellID.Row == 0 {
				headers := []string{"ID", "Название товара", "Цена", "Кол-во", "Категория"}
				label.SetText(headers[cellID.Col])
			} else {
				emp := products[cellID.Row-1]

				switch cellID.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", emp.ID))
				case 1:
					label.SetText(emp.Name)
				case 2:
					label.SetText(fmt.Sprintf("%.2f", emp.Price))
				case 3:
					label.SetText(fmt.Sprintf("%d", emp.Stock))
				case 4:
					label.SetText(emp.Category.Name)
				}
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 230)
		table.SetColumnWidth(2, 100)
		table.SetColumnWidth(3, 60)
		table.SetColumnWidth(4, 230)
		table.Resize(fyne.NewSize(700, 250))
		table.Move(fyne.NewPos(370, 100))
		content.Add(table)
	}

	noteLabel := widget.NewLabel("Обязательно введите ID товара, остальные поля заполните только если хотите изменить их")
	noteLabel.TextStyle.Bold = true
	noteLabel.Resize(fyne.NewSize(700, 20))
	noteLabel.Move(fyne.NewPos(295, 60))
	content.Add(noteLabel)

	IdEntry := widget.NewEntry()
	IdEntry.SetPlaceHolder("Введите ID товара")
	IdEntry.Resize(fyne.NewSize(300, 100))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите название товара")
	nameEntry.Resize(fyne.NewSize(300, 100))

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Введите описание товара не выходя\nза границы поля по ширине\nНе более 12 строчек")
	descriptionEntry.Resize(fyne.NewSize(300, 300))

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Введите цену товара")
	priceEntry.Resize(fyne.NewSize(300, 100))

	categoryNameEntry := widget.NewEntry()
	categoryNameEntry.SetPlaceHolder("Введите название категории")
	categoryNameEntry.Resize(fyne.NewSize(300, 100))

	stockEntry := widget.NewEntry()
	stockEntry.SetPlaceHolder("Введите количество товара")
	stockEntry.Resize(fyne.NewSize(300, 100))

	otstup := widget.NewLabel(" ")

	okBtn := widget.NewButton("Обновить", func() {
		val, err := strconv.ParseUint(IdEntry.Text, 10, 64)
		if err != nil {
			dialog.NewError(errors.New("Некорректный ID"), window).Show()
			return
		}
		uval := uint(val)

		var categoryID uint
		if categoryNameEntry.Text != "" {
			category, err := rep.FindCategory(categoryNameEntry.Text)
			if err != nil {
				dialog.NewError(err, window).Show()
				return
			}
			categoryID = category.ID
		} else {
			categoryID = uint(products[val-1].CategoryID)
		}

		var floatPrice float64
		if priceEntry.Text != "" {
			floatPrice, err = strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.NewError(errors.New("Некорректная цена"), window).Show()
				return
			}
		}
		var intStock int
		if stockEntry.Text != "" {
			intStock, err = strconv.Atoi(stockEntry.Text)
			if err != nil {
				dialog.NewError(errors.New("Некорректное количество"), window).Show()
				return
			}
		}

		newProduct := &models.Product{
			ID:          uval,
			Name:        nameEntry.Text,
			Description: descriptionEntry.Text,
			Price:       floatPrice,
			CategoryID:  categoryID,
			Stock:       intStock,
		}

		err = rep.UpdateProduct(newProduct)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			dialog.NewInformation("Успех", "Товар "+newProduct.Name+" успешно обновлен!", window).Show()
			UpdateProduct(user, app, window, content)
		}
	})

	okBtn.Resize(fyne.NewSize(200, 50))

	Categoryform := container.NewVBox(IdEntry, nameEntry, descriptionEntry, priceEntry, categoryNameEntry, stockEntry, otstup, okBtn)
	Categoryform.Resize(fyne.NewSize(300, 300))
	Categoryform.Move(fyne.NewPos(500, 400))

	content.Add(Categoryform)
}

func DeleteProduct(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	AdminSidebar(user, app, window, content)
	z := canvas.NewText("Удалить товар:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	products, err := rep.GetProducts()
	if err != nil {
		dialog.NewError(err, window).Show()
	} else {
		getSize := func() (int, int) {
			return len(products) + 1, 5
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if cellID.Row == 0 {
				headers := []string{"ID", "Название товара", "Цена", "Кол-во", "Категория"}
				label.SetText(headers[cellID.Col])
			} else {
				emp := products[cellID.Row-1]

				switch cellID.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", emp.ID))
				case 1:
					label.SetText(emp.Name)
				case 2:
					label.SetText(fmt.Sprintf("%.2f", emp.Price))
				case 3:
					label.SetText(fmt.Sprintf("%d", emp.Stock))
				case 4:
					label.SetText(emp.Category.Name)
				}
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 230)
		table.SetColumnWidth(2, 100)
		table.SetColumnWidth(3, 60)
		table.SetColumnWidth(4, 230)
		table.Resize(fyne.NewSize(700, 250))
		table.Move(fyne.NewPos(370, 100))
		content.Add(table)
	}

	productIndex := widget.NewEntry()
	productIndex.SetPlaceHolder("Введите название товара")
	productIndex.Resize(fyne.NewSize(300, 100))

	okBtn := widget.NewButton("Удалить", func() {
		if len(products) == 0 {
			dialog.NewInformation("Ошибка", "Нет товаров для удаления", window).Show()
		} else {
			err := rep.DeleteProduct(productIndex.Text)
			if err != nil {
				dialog.NewError(err, window).Show()
			} else {
				dialog.NewInformation("Успех", "Товар "+productIndex.Text+" успешно удален!", window).Show()
				DeleteProduct(user, app, window, content)
			}
		}
	})
	okBtn.Resize(fyne.NewSize(200, 50))

	otstup := widget.NewLabel(" ")

	Categoryform := container.NewVBox(productIndex, otstup, okBtn)
	Categoryform.Resize(fyne.NewSize(300, 300))
	Categoryform.Move(fyne.NewPos(500, 400))

	content.Add(Categoryform)

}

func CreateProduct(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	AdminSidebar(user, app, window, content)
	z := canvas.NewText("Новый товар:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	products, err := rep.GetProducts()
	if err != nil {
		dialog.NewError(err, window).Show()
	} else {
		getSize := func() (int, int) {
			return len(products) + 1, 5
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if cellID.Row == 0 {
				headers := []string{"ID", "Название товара", "Цена", "Кол-во", "Категория"}
				label.SetText(headers[cellID.Col])
			} else {
				emp := products[cellID.Row-1]

				switch cellID.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", emp.ID))
				case 1:
					label.SetText(emp.Name)
				case 2:
					label.SetText(fmt.Sprintf("%.2f", emp.Price))
				case 3:
					label.SetText(fmt.Sprintf("%d", emp.Stock))
				case 4:
					label.SetText(emp.Category.Name)
				}
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 230)
		table.SetColumnWidth(2, 100)
		table.SetColumnWidth(3, 60)
		table.SetColumnWidth(4, 230)
		table.Resize(fyne.NewSize(700, 250))
		table.Move(fyne.NewPos(370, 100))
		content.Add(table)
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите название товара")
	nameEntry.Resize(fyne.NewSize(300, 100))

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Введите описание товара не выходя\nза границы поля по ширине\nНе более 12 строчек")
	descriptionEntry.Resize(fyne.NewSize(300, 500))

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Введите цену товара")
	priceEntry.Resize(fyne.NewSize(300, 100))

	categoryNameEntry := widget.NewEntry()
	categoryNameEntry.SetPlaceHolder("Введите название категории")
	categoryNameEntry.Resize(fyne.NewSize(300, 100))

	stockEntry := widget.NewEntry()
	stockEntry.SetPlaceHolder("Введите количество товара")
	stockEntry.Resize(fyne.NewSize(300, 100))

	otstup := widget.NewLabel(" ")

	okBtn := widget.NewButton("Создать", func() {
		category, err := rep.FindCategory(categoryNameEntry.Text)
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}
		floatPrice, err := strconv.ParseFloat(priceEntry.Text, 64)
		if err != nil {
			dialog.NewError(errors.New("Некорректная цена"), window).Show()
			return
		}

		intStock, err := strconv.Atoi(stockEntry.Text)
		if err != nil {
			dialog.NewError(errors.New("Некорректное количество"), window).Show()
			return
		}

		newProduct := &models.Product{
			Name:        nameEntry.Text,
			Description: descriptionEntry.Text,
			Price:       floatPrice,
			CategoryID:  category.ID,
			Stock:       intStock,
		}

		product, err := rep.CreateProduct(newProduct)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			dialog.NewInformation("Успех", "Товар "+product.Name+" успешно создан!", window).Show()
			CreateProduct(user, app, window, content)
		}
	})

	okBtn.Resize(fyne.NewSize(200, 50))

	Categoryform := container.NewVBox(nameEntry, descriptionEntry, priceEntry, categoryNameEntry, stockEntry, otstup, okBtn)
	Categoryform.Resize(fyne.NewSize(300, 800))
	Categoryform.Move(fyne.NewPos(500, 400))

	content.Add(Categoryform)
}

func DeleteCategory(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	AdminSidebar(user, app, window, content)
	z := canvas.NewText("Удалить категорию:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	categories, err := rep.GetCategories()
	if err != nil {
		dialog.NewError(err, window).Show()
	} else {
		getSize := func() (int, int) {
			return len(categories) + 1, 2
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if cellID.Row == 0 {
				headers := []string{"ID", "Название категории"}
				label.SetText(headers[cellID.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				emp := categories[cellID.Row-1]

				switch cellID.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", emp.ID))
				case 1:
					label.SetText(emp.Name)
				}
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 230)
		table.Resize(fyne.NewSize(300, 250))
		table.Move(fyne.NewPos(500, 100))
		content.Add(table)
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите название категории")
	nameEntry.Resize(fyne.NewSize(300, 100))

	otstup := widget.NewLabel(" ")

	okBtn := widget.NewButton("Удалить", func() {
		err := rep.DeleteCategory(nameEntry.Text)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			dialog.NewInformation("Успех", "Категория успешно удалена!", window).Show()
			DeleteCategory(user, app, window, content)
		}
	})
	okBtn.Resize(fyne.NewSize(200, 50))

	Categoryform := container.NewVBox(nameEntry, otstup, okBtn)
	Categoryform.Resize(fyne.NewSize(300, 300))
	Categoryform.Move(fyne.NewPos(500, 400))

	content.Add(Categoryform)
}

func CreateCategory(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	AdminSidebar(user, app, window, content)
	z := canvas.NewText("Новая категория:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	categories, err := rep.GetCategories()

	if err != nil {
		dialog.NewError(err, window).Show()
	} else {
		getSize := func() (int, int) {
			return len(categories) + 1, 2
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if cellID.Row == 0 {
				headers := []string{"ID", "Название категории"}
				label.SetText(headers[cellID.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				emp := categories[cellID.Row-1]

				switch cellID.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", emp.ID))
				case 1:
					label.SetText(emp.Name)
				}
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 230)
		table.Resize(fyne.NewSize(300, 250))
		table.Move(fyne.NewPos(500, 100))
		content.Add(table)
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите название категории")
	nameEntry.Resize(fyne.NewSize(300, 100))

	otstup := widget.NewLabel(" ")

	okBtn := widget.NewButton("Создать", func() {
		category, err := rep.CreateCategory(nameEntry.Text)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			dialog.NewInformation("Успех", "Категория "+category.Name+" успешно создана!", window).Show()
			CreateCategory(user, app, window, content)
		}
	})
	okBtn.Resize(fyne.NewSize(200, 50))

	Categoryform := container.NewVBox(nameEntry, otstup, okBtn)
	Categoryform.Resize(fyne.NewSize(300, 300))
	Categoryform.Move(fyne.NewPos(500, 400))

	content.Add(Categoryform)
}

func AdminSidebar(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	img := canvas.NewImageFromFile("1.jpg")
	img.Resize(fyne.NewSize(200, 200))
	img.Move(fyne.NewPos(50, 600))
	content.Add(img)

	profileBtn := widget.NewButton("Профиль", func() { Profile(user, app, window, content) })
	profileBtn.Resize(fyne.NewSize(200, 50))
	profileBtn.Move(fyne.NewPos(50, 20))
	content.Add(profileBtn)

	newCatogoryBtn := widget.NewButton("Новая категория", func() { CreateCategory(user, app, window, content) })
	newCatogoryBtn.Resize(fyne.NewSize(200, 50))
	newCatogoryBtn.Move(fyne.NewPos(50, 100))
	content.Add(newCatogoryBtn)

	deleteCatogoryBtn := widget.NewButton("Удалить категорию", func() { DeleteCategory(user, app, window, content) })
	deleteCatogoryBtn.Resize(fyne.NewSize(200, 50))
	deleteCatogoryBtn.Move(fyne.NewPos(50, 180))
	content.Add(deleteCatogoryBtn)

	newProductBtn := widget.NewButton("Новый товар", func() { CreateProduct(user, app, window, content) })
	newProductBtn.Resize(fyne.NewSize(200, 50))
	newProductBtn.Move(fyne.NewPos(50, 260))
	content.Add(newProductBtn)

	updateProductBtn := widget.NewButton("Изменить товар", func() { UpdateProduct(user, app, window, content) })
	updateProductBtn.Resize(fyne.NewSize(200, 50))
	updateProductBtn.Move(fyne.NewPos(50, 340))
	content.Add(updateProductBtn)

	deleteProductBtn := widget.NewButton("Удалить товар", func() { DeleteProduct(user, app, window, content) })
	deleteProductBtn.Resize(fyne.NewSize(200, 50))
	deleteProductBtn.Move(fyne.NewPos(50, 420))
	content.Add(deleteProductBtn)

	createStoreBth := widget.NewButton("Новый кладовщик", func() { CreateStoreMan(user, app, window, content) })
	createStoreBth.Resize(fyne.NewSize(200, 50))
	createStoreBth.Move(fyne.NewPos(50, 500))
	content.Add(createStoreBth)
}

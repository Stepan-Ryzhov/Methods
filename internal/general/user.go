package general

import (
	"errors"
	"image/color"

	//"strconv"
	"fmt"

	models "methodi_razrabotki/internal/models"
	rep "methodi_razrabotki/internal/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	// "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var selectedProductIndex int = -1

func UserOrders(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Мои заказы:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
}

func Cart(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Корзина:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
}

func TableWidget(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container, selected string, productSelect *widget.Select) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Каталог:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	headerLabel := widget.NewLabel("Название товара\t\t\t        Цена\t\t    Описание\t\t\t\t\t\tКатегория")
	headerLabel.Resize(fyne.NewSize(850, 50))
	headerLabel.Move(fyne.NewPos(300, 100))
	content.Add(headerLabel)

	products, err := rep.GetProductsUser()
	if err != nil {
		dialog.NewError(err, window).Show()
		return
	}

	// ФИЛЬТРАЦИЯ ДАННЫХ ДО создания таблицы
	var filteredProducts []models.Product
	if selected != "Все категории" {
		for _, product := range products {
			if product.Category.Name == selected {
				filteredProducts = append(filteredProducts, product)
			}
		}
	} else {
		filteredProducts = products
	}

	getSize := func() (int, int) {
		return len(filteredProducts), 4 // Используем отфильтрованные данные
	}

	createCell := func() fyne.CanvasObject {
		return widget.NewLabel("")
	}

	updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		emp := filteredProducts[cellID.Row] // Используем отфильтрованные данные

		switch cellID.Col {
		case 0:
			label.SetText(emp.Name)
		case 1:
			label.SetText(fmt.Sprintf("%.2f", emp.Price))
		case 2:
			label.SetText(emp.Description)
			label.Wrapping = fyne.TextWrapWord
		case 3:
			label.SetText(emp.Category.Name)
		}
	}

	table := widget.NewTable(getSize, createCell, updateCell)
	table.SetColumnWidth(0, 230)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 300)
	table.SetColumnWidth(3, 190)

	for i := 0; i < len(filteredProducts); i++ {
		table.SetRowHeight(i, widget.NewLabel(filteredProducts[i].Description).MinSize().Height+20)
	}

	table.Resize(fyne.NewSize(850, 400))
	table.Move(fyne.NewPos(300, 150))
	content.Add(table)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(filteredProducts) {
			selectedProductIndex = findProductIndex(products, filteredProducts[id.Row].ID)
		}
	}

	productSelect.PlaceHolder = selected
	productSelect.Resize(fyne.NewSize(200, 50))
	productSelect.Move(fyne.NewPos(950, 600))
	content.Add(productSelect)

	addToCartBtn := widget.NewButton("Добавить в корзину", func() {
		if selectedProductIndex < 0 || selectedProductIndex >= len(products) {
			dialog.NewError(errors.New("Выберите товар из таблицы"), window).Show()
			return
		}

		selectedProduct := products[selectedProductIndex].ID
		fmt.Println(selectedProduct)
		dialog.NewInformation("Успех", fmt.Sprintf("Товар %s добавлен в корзину", products[selectedProductIndex].Name), window).Show()
	})
	addToCartBtn.Resize(fyne.NewSize(200, 50))
	addToCartBtn.Move(fyne.NewPos(600, 700))
	content.Add(addToCartBtn)
}

func findProductIndex(products []models.Product, productID uint) int {
	for i, product := range products {
		if product.ID == productID {
			return i
		}
	}
	return -1
}

func getUniqueCategories(products []models.Product) []string {
	uniqueMap := make(map[string]bool)
	categories := []string{"Все категории"}

	for _, product := range products {
		if product.Category.Name != "" && !uniqueMap[product.Category.Name] {
			uniqueMap[product.Category.Name] = true
			categories = append(categories, product.Category.Name)
		}
	}

	return categories
}

func Catalog(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Каталог:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	headerLabel := widget.NewLabel("Название товара\t\t\t        Цена\t\t    Описание\t\t\t\t\t\tКатегория")
	headerLabel.Resize(fyne.NewSize(850, 50))
	headerLabel.Move(fyne.NewPos(300, 100))
	content.Add(headerLabel)

	products, err := rep.GetProductsUser()
	categories := getUniqueCategories(products)
	if err != nil {
		dialog.NewError(err, window).Show()
	} else {
		getSize := func() (int, int) {
			return len(products), 4
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			emp := products[cellID.Row]

			switch cellID.Col {
			case 0:
				label.SetText(emp.Name)
			case 1:
				label.SetText(fmt.Sprintf("%.2f", emp.Price))
			case 2:
				label.SetText(emp.Description)
				label.Wrapping = fyne.TextWrapWord
			case 3:
				label.SetText(emp.Category.Name)
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 230)
		table.SetColumnWidth(1, 100)
		table.SetColumnWidth(2, 300)
		table.SetColumnWidth(3, 190)
		table.SetRowHeight(0, 250)
		table.SetRowHeight(1, 250)
		table.SetRowHeight(2, 250)
		table.SetRowHeight(3, 250)
		table.Resize(fyne.NewSize(850, 400))
		table.Move(fyne.NewPos(300, 150))
		content.Add(table)

		var productSelect *widget.Select
		productSelect = widget.NewSelect(categories, func(selected string) {
			TableWidget(user, app, window, content, selected, productSelect)
		})
		productSelect.PlaceHolder = "Все категории"
		productSelect.Resize(fyne.NewSize(200, 50))
		productSelect.Move(fyne.NewPos(950, 600))
		content.Add(productSelect)

		addToCartBtn := widget.NewButton("Добавить в корзину", func() {
			if selectedProductIndex < 0 || selectedProductIndex >= len(products) {
				dialog.NewError(errors.New("Выберите товар из таблицы"), window).Show()
				return
			}

			selectedProduct := products[selectedProductIndex].ID
			fmt.Println(selectedProduct)
			dialog.NewInformation("Успех", fmt.Sprintf("Товар %s добавлен в корзину", products[selectedProductIndex].Name), window).Show()
		})
		addToCartBtn.Resize(fyne.NewSize(200, 50))
		addToCartBtn.Move(fyne.NewPos(600, 700))
		content.Add(addToCartBtn)

		table.OnSelected = func(id widget.TableCellID) {
			selectedProductIndex = id.Row
		}
	}
}

func UserSidebar(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	img := canvas.NewImageFromFile("1.jpg")
	img.Resize(fyne.NewSize(200, 200))
	img.Move(fyne.NewPos(50, 600))
	content.Add(img)

	profileBtn := widget.NewButton("Профиль", func() { Profile(user, app, window, content) })
	profileBtn.Resize(fyne.NewSize(200, 50))
	profileBtn.Move(fyne.NewPos(50, 20))
	content.Add(profileBtn)

	cartBtn := widget.NewButton("Корзина", func() { Cart(user, app, window, content) })
	cartBtn.Resize(fyne.NewSize(200, 50))
	cartBtn.Move(fyne.NewPos(50, 120))
	content.Add(cartBtn)

	catalogBtn := widget.NewButton("Каталог", func() { Catalog(user, app, window, content) })
	catalogBtn.Resize(fyne.NewSize(200, 50))
	catalogBtn.Move(fyne.NewPos(50, 220))
	content.Add(catalogBtn)

	orderBtn := widget.NewButton("Мои заказы", func() { UserOrders(user, app, window, content) })
	orderBtn.Resize(fyne.NewSize(200, 50))
	orderBtn.Move(fyne.NewPos(50, 320))
	content.Add(orderBtn)
}

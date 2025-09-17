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

	db "methodi_razrabotki/internal/database"

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

	orders, err := rep.GetOrdersByUserID(user.ID)
	if err != nil {
		dialog.NewError(err, window).Show()
		return
	}

	if len(orders) == 0 {
		z := canvas.NewText("Заказов нет", color.Black)
		z.TextSize = 26
		z.Move(fyne.NewPos(600, 200))
		content.Add(z)
		return
	}

	yPos := 100

	for _, order := range orders {
		orderHeader := widget.NewLabel(fmt.Sprintf("Заказ №%d | Сумма: %.2f | Статус: %s", order.ID, order.Total, order.Status))
		orderHeader.TextStyle.Bold = true
		orderHeader.Resize(fyne.NewSize(800, 40))
		orderHeader.Move(fyne.NewPos(300, float32(yPos)))
		content.Add(orderHeader)
		yPos += 40

		if len(order.Items) > 0 {
			headerLabel := widget.NewLabel("Название товара\t\tКоличество\tЦена")
			headerLabel.Resize(fyne.NewSize(800, 30))
			headerLabel.Move(fyne.NewPos(300, float32(yPos)))
			content.Add(headerLabel)
			yPos += 30

			getSize := func() (int, int) {
				return len(order.Items), 3
			}

			createCell := func() fyne.CanvasObject {
				return widget.NewLabel("")
			}

			updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
				label := obj.(*widget.Label)
				item := order.Items[cellID.Row]

				switch cellID.Col {
				case 0:
					product, err := rep.GetProductByID(item.ProductID)
					if err != nil {
						dialog.NewError(err, window).Show()
					}
					label.SetText(product.Name)
				case 1:
					label.SetText(fmt.Sprintf("%d", item.Quantity))
				case 2:
					label.SetText(fmt.Sprintf("%.2f", item.Price))
				}
			}

			table := widget.NewTable(getSize, createCell, updateCell)
			table.SetColumnWidth(0, 300)
			table.SetColumnWidth(1, 100)
			table.SetColumnWidth(2, 100)
			table.Resize(fyne.NewSize(550, float32(len(order.Items))*50))
			table.Move(fyne.NewPos(300, float32(yPos)))
			content.Add(table)
			yPos += int(table.Size().Height) + 20
		} else {
			noItemsLabel := widget.NewLabel("Товаров в этом заказе не найдено")
			noItemsLabel.Resize(fyne.NewSize(500, 30))
			noItemsLabel.Move(fyne.NewPos(300, float32(yPos)))
			content.Add(noItemsLabel)
			yPos += 30
		}
		cancelBtn := widget.NewButton("Отменить заказ", func() {
			if order.ID == 0 {
				dialog.NewError(errors.New("Некорректный идентификатор заказа"), window).Show()
				return
			}

			if err := rep.DeleteOrder(order.ID, user.ID); err != nil {
				dialog.NewError(err, window).Show()
				return
			}

			dialog.NewInformation("Успех", "Заказ отменён", window).Show()
			UserOrders(user, app, window, content)
		})

		cancelBtn.Resize(fyne.NewSize(200, 40))
		cancelBtn.Move(fyne.NewPos(900, float32(yPos-40)))
		content.Add(cancelBtn)
		if order.Status == "Оформлен" || order.Status == "Сформирован" {
			payBtn := widget.NewButton("Оплатить", func() {
				if order.ID == 0 {
					dialog.NewError(errors.New("Некорректный идентификатор заказа"), window).Show()
					return
				} else {
					if err := rep.UpdateOrderStatus(order.ID, "Оплачен"); err != nil {
						dialog.NewError(err, window).Show()
					} else {
						dialog.NewInformation("Успех", "Заказ оплачен", window).Show()
						UserOrders(user, app, window, content)
					}

				}
			})

			payBtn.Resize(fyne.NewSize(200, 40))
			payBtn.Move(fyne.NewPos(900, float32(yPos-90)))
			content.Add(payBtn)
		}
		if order.Status == "Отправлен" {
			takeBtn := widget.NewButton("Получить заказ", func() {
				if order.ID == 0 {
					dialog.NewError(errors.New("Некорректный идентификатор заказа"), window).Show()
					return
				}

				if err := rep.DeleteOrder(order.ID, user.ID); err != nil {
					dialog.NewError(err, window).Show()
					return
				}

				dialog.NewInformation("Успех", "Заказ получен", window).Show()
				UserOrders(user, app, window, content)
			})
			takeBtn.Resize(fyne.NewSize(200, 40))
			takeBtn.Move(fyne.NewPos(900, float32(yPos-40)))
			content.Add(takeBtn)
		}
	}
}

func Cart(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	UserSidebar(user, app, window, content)

	z := canvas.NewText("Корзина:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)
	cart_id, err := rep.GetCartID(user.ID)
	if err != nil {
		dialog.NewError(err, window).Show()
	}
	cart, err := rep.GetCartByID(cart_id)
	if err != nil {
		dialog.NewError(err, window).Show()
	}
	if len(cart.Items) == 0 {
		z := canvas.NewText("Корзина пуста", color.Black)
		z.TextSize = 26
		z.Move(fyne.NewPos(600, 200))
		content.Add(z)
	} else {
		headerLabel := widget.NewLabel("Название товара\t\t\t        Цена\t\t   Количество       Стоимость")
		headerLabel.Resize(fyne.NewSize(850, 50))
		headerLabel.Move(fyne.NewPos(300, 100))
		content.Add(headerLabel)

		var total float64 = 0

		products := []string{}
		for _, item := range cart.Items {
			product, err := rep.GetProductByID(item.ProductID)
			if err != nil {
				dialog.NewError(err, window).Show()
			} else {
				products = append(products, product.Name)
				total += item.Total
			}
		}
		getSize := func() (int, int) {
			return len(products), 4
		}

		createCell := func() fyne.CanvasObject {
			return widget.NewLabel("")
		}

		updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			emp2 := cart.Items[cellID.Row]
			switch cellID.Col {
			case 0:
				label.SetText(products[cellID.Row])
			case 1:
				label.SetText(fmt.Sprintf("%.2f", emp2.Price))
			case 2:
				label.SetText(fmt.Sprintf("%d", emp2.Quantity))
			case 3:
				label.SetText(fmt.Sprintf("%.2f", emp2.Total))
			}
		}

		table := widget.NewTable(getSize, createCell, updateCell)
		table.SetColumnWidth(0, 230)
		table.SetColumnWidth(1, 100)
		table.SetColumnWidth(2, 100)
		table.SetColumnWidth(3, 100)

		table.Resize(fyne.NewSize(550, 400))
		table.Move(fyne.NewPos(300, 150))
		content.Add(table)

		orderBtn := widget.NewButton("Оформить заказ", func() {
			orderItems := make([]models.OrderItem, len(cart.Items))
			for i, item := range cart.Items {
				orderItems[i] = models.OrderItem{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
					Price:     item.Price,
					Total:     item.Total,
				}
			}

			order := models.Order{
				UserID: user.ID,
				Total:  total,
				Status: "Оформлен",
				Items:  orderItems,
			}

			if err := rep.CreateOrder(&order); err != nil {
				dialog.NewError(err, window).Show()
			} else {
				dialog.NewInformation("Успех", "Заказ оформлен", window).Show()
				for _, item := range cart.Items {
					if err := rep.RemoveFromCart(cart_id, item.ProductID); err != nil {
						dialog.NewError(err, window).Show()
					}
				}

				if err := db.GetDB().Save(&cart).Error; err != nil {
					dialog.NewError(errors.New("Ошибка очистки корзины"), window).Show()
					return
				}

				Cart(user, app, window, content)
			}
		})
		orderBtn.Resize(fyne.NewSize(200, 50))
		orderBtn.Move(fyne.NewPos(900, 700))
		content.Add(orderBtn)

		z1 := canvas.NewText(fmt.Sprintf("Итого: %.2f", total), color.Black)
		z1.TextSize = 28
		z1.Move(fyne.NewPos(300, 700))
		content.Add(z1)

		uprLabel := widget.NewLabel("Выберите товар из корзины для\nуправления его количеством")
		uprLabel.Resize(fyne.NewSize(200, 50))
		uprLabel.Move(fyne.NewPos(900, 100))
		content.Add(uprLabel)

		table.OnSelected = func(id widget.TableCellID) {
			selectedProductIndex = id.Row

			plusOneBtn := widget.NewButton("+1", func() {
				if selectedProductIndex >= 0 && selectedProductIndex < len(cart.Items) {
					item := cart.Items[selectedProductIndex]
					err := rep.IncrementItem(cart_id, item.ProductID)
					if err != nil {
						dialog.NewError(err, window).Show()
					} else {
						dialog.ShowInformation("Успех", "Количество увеличено", window)
						Cart(user, app, window, content)
					}
				}
			})

			minusOneBtn := widget.NewButton("-1", func() {
				if selectedProductIndex >= 0 && selectedProductIndex < len(cart.Items) {
					item := cart.Items[selectedProductIndex]
					err := rep.DecrementItem(cart_id, item.ProductID)
					if err != nil {
						dialog.NewError(err, window).Show()
					} else {
						dialog.ShowInformation("Успех", "Количество уменьшено", window)
						Cart(user, app, window, content)
					}
				}
			})

			removeBtn := widget.NewButton("Удалить товар", func() {
				if selectedProductIndex >= 0 && selectedProductIndex < len(cart.Items) {
					item := cart.Items[selectedProductIndex]
					err := rep.RemoveFromCart(cart_id, item.ProductID)
					if err != nil {
						dialog.NewError(err, window).Show()
					} else {
						dialog.ShowInformation("Успех", "Товар удален", window)
						Cart(user, app, window, content)
					}
				}
			})

			plusOneBtn.Resize(fyne.NewSize(200, 50))
			plusOneBtn.Move(fyne.NewPos(900, 200))
			content.Add(plusOneBtn)

			minusOneBtn.Resize(fyne.NewSize(200, 50))
			minusOneBtn.Move(fyne.NewPos(900, 300))
			content.Add(minusOneBtn)

			removeBtn.Resize(fyne.NewSize(200, 50))
			removeBtn.Move(fyne.NewPos(900, 400))
			content.Add(removeBtn)

		}
	}
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
		return len(filteredProducts), 4
	}

	createCell := func() fyne.CanvasObject {
		return widget.NewLabel("")
	}

	updateCell := func(cellID widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		emp := filteredProducts[cellID.Row]

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
			quantityEntry := widget.NewEntry()
			quantityEntry.PlaceHolder = "Укажите количество"
			quantityEntry.Resize(fyne.NewSize(200, 40))
			quantityEntry.Move(fyne.NewPos(600, 650))
			content.Add(quantityEntry)

			addToCartBtn := widget.NewButton("Добавить в корзину", func() {
				if selectedProductIndex < 0 || selectedProductIndex >= len(products) {
					dialog.NewError(errors.New("Выберите товар из таблицы"), window).Show()
					return
				}

				selectedProduct := products[selectedProductIndex].ID
				cart_id, err := rep.GetCartID(user.ID)
				if err != nil {
					dialog.NewError(err, window).Show()
				}
				quantity, err := strconv.Atoi(quantityEntry.Text)
				if err != nil {
					dialog.NewError(errors.New("Некорректное количество"), window).Show()
					return
				}
				item := models.CartItem{
					CartID:    cart_id,
					ProductID: selectedProduct,
					Quantity:  quantity,
					Price:     products[selectedProductIndex].Price,
					Total:     products[selectedProductIndex].Price * float64(quantity),
				}
				if err := rep.AddItemToCart(&item); err != nil {
					dialog.NewError(err, window).Show()
				} else {
					dialog.NewInformation("Успех", fmt.Sprintf("Товар %s добавлен в корзину", products[selectedProductIndex].Name), window).Show()
				}
			})
			addToCartBtn.Resize(fyne.NewSize(200, 50))
			addToCartBtn.Move(fyne.NewPos(600, 700))
			content.Add(addToCartBtn)
		}
	}

	productSelect.PlaceHolder = selected
	productSelect.Resize(fyne.NewSize(200, 50))
	productSelect.Move(fyne.NewPos(950, 550))
	content.Add(productSelect)
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

		for i := 0; i < len(products); i++ {
			table.SetRowHeight(i, widget.NewLabel(products[i].Description).MinSize().Height+20)
		}
		table.Resize(fyne.NewSize(850, 400))
		table.Move(fyne.NewPos(300, 150))
		content.Add(table)

		var productSelect *widget.Select
		productSelect = widget.NewSelect(categories, func(selected string) {
			TableWidget(user, app, window, content, selected, productSelect)
		})
		productSelect.PlaceHolder = "Все категории"
		productSelect.Resize(fyne.NewSize(200, 50))
		productSelect.Move(fyne.NewPos(950, 550))
		content.Add(productSelect)

		table.OnSelected = func(id widget.TableCellID) {
			selectedProductIndex = id.Row
			quantityEntry := widget.NewEntry()
			quantityEntry.PlaceHolder = "Укажите количество"
			quantityEntry.Resize(fyne.NewSize(200, 40))
			quantityEntry.Move(fyne.NewPos(600, 650))
			content.Add(quantityEntry)

			addToCartBtn := widget.NewButton("Добавить в корзину", func() {
				if selectedProductIndex < 0 || selectedProductIndex >= len(products) {
					dialog.NewError(errors.New("Выберите товар из таблицы"), window).Show()
					return
				}

				selectedProduct := products[selectedProductIndex].ID
				cart_id, err := rep.GetCartID(user.ID)
				if err != nil {
					dialog.NewError(err, window).Show()
				}
				quantity, err := strconv.Atoi(quantityEntry.Text)
				if err != nil {
					dialog.NewError(errors.New("Некорректное количество"), window).Show()
					return
				}
				item := models.CartItem{
					CartID:    cart_id,
					ProductID: selectedProduct,
					Quantity:  quantity,
					Price:     products[selectedProductIndex].Price,
					Total:     products[selectedProductIndex].Price * float64(quantity),
				}
				if err := rep.AddItemToCart(&item); err != nil {
					dialog.NewError(err, window).Show()
				} else {
					dialog.NewInformation("Успех", fmt.Sprintf("Товар %s добавлен в корзину", products[selectedProductIndex].Name), window).Show()
				}
			})
			addToCartBtn.Resize(fyne.NewSize(200, 50))
			addToCartBtn.Move(fyne.NewPos(600, 700))
			content.Add(addToCartBtn)
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

package general

import (
	"errors"
	"fmt"
	"image/color"

	//"strconv"

	models "methodi_razrabotki/internal/models"
	rep "methodi_razrabotki/internal/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	//"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Forming(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	StoreSidebar(user, app, window, content)

	z := canvas.NewText("Формирование заказов:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	orders, err := rep.GetOrders()
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
		if order.Status != "Отправлен" {
			if order.Status != "Сформирован" && order.Status != "Оплачен" {
				cancelBtn := widget.NewButton("Сформировать", func() {
					if order.ID == 0 {
						dialog.NewError(errors.New("Некорректный идентификатор заказа"), window).Show()
						return
					}

					if err := rep.UpdateOrderStatus(order.ID, "Сформирован"); err != nil {
						dialog.NewError(err, window).Show()
					} else {
						dialog.NewInformation("Успех", "Заказ сформирован", window).Show()
						Forming(user, app, window, content)
					}
				})
				cancelBtn.Resize(fyne.NewSize(200, 40))
				cancelBtn.Move(fyne.NewPos(900, float32(yPos-40)))
				content.Add(cancelBtn)
			}

			if order.Status == "Оплачен" {
				payBtn := widget.NewButton("Сформировать и отправить", func() {
					if order.ID == 0 {
						dialog.NewError(errors.New("Некорректный идентификатор заказа"), window).Show()
						return
					} else {
						if err := rep.UpdateOrderStatus(order.ID, "Отправлен"); err != nil {
							dialog.NewError(err, window).Show()
						} else {
							dialog.NewInformation("Успех", "Заказ отправлен", window).Show()
							Forming(user, app, window, content)
						}
					}
				})
				payBtn.Resize(fyne.NewSize(200, 40))
				payBtn.Move(fyne.NewPos(900, float32(yPos-40)))
				content.Add(payBtn)
			}
		}
	}
}

func StoreSidebar(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	img := canvas.NewImageFromFile("1.jpg")
	img.Resize(fyne.NewSize(200, 200))
	img.Move(fyne.NewPos(50, 600))
	content.Add(img)

	profileBtn := widget.NewButton("Профиль", func() { Profile(user, app, window, content) })
	profileBtn.Resize(fyne.NewSize(200, 50))
	profileBtn.Move(fyne.NewPos(50, 20))
	content.Add(profileBtn)

	formedBtn := widget.NewButton("Формирование", func() { Forming(user, app, window, content) })
	formedBtn.Resize(fyne.NewSize(200, 50))
	formedBtn.Move(fyne.NewPos(50, 120))
	content.Add(formedBtn)
}

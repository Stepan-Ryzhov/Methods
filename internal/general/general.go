package general

import (
	"image/color"

	models "methodi_razrabotki/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func Profile(user *models.User, app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()

	z := canvas.NewText("Профиль пользователя:", color.Black)
	z.TextSize = 32
	z.Move(fyne.NewPos(300, 20))
	content.Add(z)

	name := canvas.NewText("Имя: "+user.FirstName, color.Black)
	name.TextSize = 24
	name.Move(fyne.NewPos(300, 100))
	content.Add(name)

	lastname := canvas.NewText("Фамилия: "+user.LastName, color.Black)
	lastname.TextSize = 24
	lastname.Move(fyne.NewPos(300, 150))
	content.Add(lastname)

	email := canvas.NewText("Электронная почта: "+user.Email, color.Black)
	email.TextSize = 24
	email.Move(fyne.NewPos(300, 200))
	content.Add(email)

	role := canvas.NewText("Роль: "+user.Role, color.Black)
	if user.Role == "admin" {
		role.Text = "Роль: Администратор"
		AdminSidebar(user, app, window, content)
	} else if user.Role == "user" {
		role.Text = "Роль: Пользователь"
		UserSidebar(user, app, window, content)
	} else {
		role.Text = "Роль: Кладовщик"
		StoreSidebar(user, app, window, content)
	}

	role.TextSize = 24
	role.Move(fyne.NewPos(300, 250))
	content.Add(role)

	createdAt := canvas.NewText("Дата регистрации: "+user.CreatedAt.String(), color.Black)
	createdAt.TextSize = 24
	createdAt.Move(fyne.NewPos(300, 300))
	content.Add(createdAt)
}

package login

import (
	gn "methodi_razrabotki/internal/general"

	models "methodi_razrabotki/internal/models"
	rep "methodi_razrabotki/internal/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func logo(content *fyne.Container) {
	img := canvas.NewImageFromFile("1.jpg")
	img.Resize(fyne.NewSize(400, 400))
	img.Move(fyne.NewPos(400, 50))
	content.Add(img)
}

func RegisterFront(app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	logo(content)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Введите имя")

	lastnameEntry := widget.NewEntry()
	lastnameEntry.SetPlaceHolder("Введите фамилию")

	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Введите e-mail")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Придумайте пароль")

	loginBtn := widget.NewButton("Зарегистрироваться", func() {
		registerRequest := &models.RegisterRequest{
			FirstName: nameEntry.Text,
			LastName:  lastnameEntry.Text,
			Email:     loginEntry.Text,
			Password:  passwordEntry.Text,
		}
		err := rep.Register(registerRequest)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			dialog.NewInformation("Успех", "Вы успешно зарегистрировались! Воспользуйтесь меню входа для доступа к приложению.", window).Show()
			Start(app, window, content)
		}
	})
	loginBtn.Resize(fyne.NewSize(200, 50))

	otstup := widget.NewLabel(" ")

	loginform := container.NewVBox(nameEntry, lastnameEntry, loginEntry, passwordEntry, otstup, loginBtn)
	loginform.Resize(fyne.NewSize(200, 200))
	loginform.Move(fyne.NewPos(500, 450))

	content.Add(loginform)
}

func LoginFront(app fyne.App, window fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	logo(content)

	loginEntry := widget.NewEntry()
	loginEntry.SetPlaceHolder("Введите e-mail")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	loginBtn := widget.NewButton("Войти", func() {
		LoginRequest := models.LoginRequest{
			Email:    loginEntry.Text,
			Password: passwordEntry.Text,
		}

		user, err := rep.Login(LoginRequest)
		if err != nil {
			dialog.NewError(err, window).Show()
		} else {
			gn.Profile(user, app, window, content)
		}
	})
	loginBtn.Resize(fyne.NewSize(200, 50))

	otstup := widget.NewLabel(" ")

	loginform := container.NewVBox(loginEntry, passwordEntry, otstup, loginBtn)
	loginform.Resize(fyne.NewSize(200, 200))
	loginform.Move(fyne.NewPos(500, 450))

	content.Add(loginform)
}

func Start(myApp fyne.App, myWindow fyne.Window, content *fyne.Container) {
	content.RemoveAll()
	logo(content)

	loginBtn := widget.NewButton("Войти", func() {
		LoginFront(myApp, myWindow, content)
	})
	loginBtn.Resize(fyne.NewSize(200, 50))
	loginBtn.Move(fyne.NewPos(390, 450))

	registerBtn := widget.NewButton("Зарегистрироваться", func() {
		RegisterFront(myApp, myWindow, content)
	})

	registerBtn.Resize(fyne.NewSize(200, 50))
	registerBtn.Move(fyne.NewPos(610, 450))

	content.Add(loginBtn)
	content.Add(registerBtn)
}

package password

import (
	"net/http"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"
	echo "github.com/labstack/echo/v4"
)

// Custom errors
var (
	ErrIncorrectPassword = echo.NewHTTPError(http.StatusBadRequest, "Senha incorrecta")
	ErrInsecurePassword  = echo.NewHTTPError(http.StatusBadRequest, "Senha insegura")
)

// Change changes user's password
func (p *Password) Update(c echo.Context, userID string, oldPass, newPass string) error {
	if err := p.rbac.EnforceUser(c, userID); err != nil {
		return err
	}

	u, err := p.udb.View(p.db, userID)
	if err = zaplog.ZLog(err); err != nil {
		return err
	}

	au := p.rbac.User(c)

	//Only operators and higher can change passwords without knowing old password.
	//Users can only change passwords for roles below them
	if !p.sec.HashMatchesPassword(u.Password, oldPass) && (au.Role > model.OperatorRole || au.Role >= u.Role) {
		return ErrIncorrectPassword
	}

	if !p.sec.Password(newPass, u.Name, u.Username, u.Email) {
		return ErrInsecurePassword
	}

	u.ChangePassword(p.sec.Hash(newPass))

	return p.udb.Update(p.db, &model.User{Base: model.Base{ID: u.ID}, Password: u.Password})
}

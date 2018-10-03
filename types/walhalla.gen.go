// The code was generated. Dont edit the one
//
package types

import (
	"errors"
	"mime/multipart"
)

// NOTES:
// 1. insead of '&lt;'(less) the template generates '&it;'. I don't know why but it is.

func (ct *APIUser) Validate() bool {

	if !(len(ct.Password) >= 6) {
		return false
	}

	if len(ct.Password) > 100 {
		return false
	}

	return true
}

func (ct *APIUser) UnmarshalForm(form *multipart.Form) error {

	if values, ok := form.Value["username"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.Username = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if values, ok := form.Value["password"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.Password = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	return nil
}

func (ct *APISignUp) Validate() bool {

	return true
}

func (ct *APISignUp) UnmarshalForm(form *multipart.Form) error {

	if values, ok := form.Value["username"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.Username = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if values, ok := form.Value["password"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.Password = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if files, ok := form.File["avatar"]; len(files) != 0 && ok {
		if len(files) == 1 {
			file := files[0]
			data := make([]byte, file.Size)
			entery, err := file.Open()
			if err != nil {
				return err
			}

			entery.Read(data)
			ct.Avatar = data
		} else {
			return errors.New("multipl files with the same name were detected")
		}
	}

	return nil
}

func (ct *APIProfile) Validate() bool {

	return true
}

func (ct *APIProfile) UnmarshalForm(form *multipart.Form) error {

	if values, ok := form.Value["username"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.Username = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if values, ok := form.Value["avatarSource"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.AvatarURI = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	return nil
}

func (ct *APIEditProfile) Validate() bool {

	return true
}

func (ct *APIEditProfile) UnmarshalForm(form *multipart.Form) error {

	if values, ok := form.Value["username"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.Username = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if values, ok := form.Value["curPassword"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.CurPassword = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if values, ok := form.Value["newPassword"]; len(values) != 0 && ok {
		if len(values) == 1 {

			ct.NewPassword = values[0]

		} else {
			return errors.New("multipl values with the same name were detected")
		}

	}

	if files, ok := form.File["avatar"]; len(files) != 0 && ok {
		if len(files) == 1 {
			file := files[0]
			data := make([]byte, file.Size)
			entery, err := file.Open()
			if err != nil {
				return err
			}

			entery.Read(data)
			ct.Avatar = data
		} else {
			return errors.New("multipl files with the same name were detected")
		}
	}

	return nil
}

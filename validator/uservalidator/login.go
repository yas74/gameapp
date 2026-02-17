package uservalidator

import (
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/pkg/errmsg"
	"gocasts/gameapp/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateLoginRequest(req dto.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.doesPhoneNumberExist)),

		validation.Field(&req.Password, validation.Required),
	); err != nil {
		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {

			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithErr(err)
	}

	return nil, nil

}

func (v Validator) doesPhoneNumberExist(value interface{}) error {
	phonenumber := value.(string)

	if _, err := v.repo.GetUserByPhoneNumber(phonenumber); err != nil {
		return fmt.Errorf(errmsg.ErrorMsgNotFound)
	}

	return nil
}

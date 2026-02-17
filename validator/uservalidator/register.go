package uservalidator

import (
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/pkg/errmsg"
	"gocasts/gameapp/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		// TODO - add 3 to config
		validation.Field(&req.Name,
			validation.Required,
			validation.Length(3, 50)),

		validation.Field(&req.Password,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9@#%^&*]{8,}$`))),

		validation.Field(&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),
			validation.By(v.checkPhoneNumberUniqness)),
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

func (v Validator) checkPhoneNumberUniqness(value interface{}) error {
	phonenumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phonenumber); err != nil || !isUnique {
		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
		}
	}

	return nil
}

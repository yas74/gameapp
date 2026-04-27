package matchingvalidator

import (
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/entity"
	"gocasts/gameapp/pkg/errmsg"
	"gocasts/gameapp/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateAddToWaitingListRequest(req dto.AddToWaitingListRequest) (map[string]string, error) {
	op := richerror.Op("matchingvalidator.ValidateAddToWaitingListRequest")

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Category,
			validation.Required,
			validation.By(v.isCategoryValid)),
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

func (v Validator) isCategoryValid(value interface{}) error {
	category := value.(entity.Category)

	if !category.IsValid() {
		return fmt.Errorf(errmsg.ErrorMsgCategoryIsNotValid)
	}

	return nil
}

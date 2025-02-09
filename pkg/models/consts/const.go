package consts

const (
	ErrInternal               = "content.api.errors.payApi.internalError"          // Внутренняя ошибка
	ErrBadRequest             = "content.api.errors.payApi.badRequest"             // Плохой запрос
	ErrMethodNotAllowed       = "content.api.errors.payApi.methodNotAllowed"       // Метод не поддерживается
	ErrForbidden              = "content.api.errors.payApi.forbidden"              // Доступ запрещен
	ErrInvalidRequest         = "content.api.errors.payApi.invalidRequest"         // Неправильный запрос
	ErrAccessDenied           = "content.api.errors.payApi.accessDenied"           // Отказано в доступе
	ErrFailedGetOldSupplierID = "content.api.errors.payApi.failedGetOldSupplierId" //"failed get old supplierID"
	ErrNotEnoughMoney         = "content.api.errors.payApi.notEnoughMoney"         // На балансе недостаточно средств 402
	ErrUnknownCountryCode     = "content.api.errors.payApi.unknownCountryCode"     // Неизвестный код страны
)


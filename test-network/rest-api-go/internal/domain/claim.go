package domain

import (
	"mime/multipart"
)

type ClaimServiceInterface interface {
	StoreClaim(file *multipart.FileHeader) error
}

type ClaimRepositoryInterface interface {
	SaveFile(file *multipart.FileHeader, filename string) error
}

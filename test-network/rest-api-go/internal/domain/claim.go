package domain

import (
	"mime/multipart"
)

type ClaimServiceInterface interface {
	StoreClaim(file *multipart.FileHeader, uploadDir string) error
}

type ClaimRepositoryInterface interface {
	SaveFile(file *multipart.FileHeader, uploadDir, filename string) error
}

package dto

import "github.com/oliveirabalsa/go-globalhitss-be/app/model"

type PaginationResponse struct {
	Data       []*model.Client `json:"data"`
	Page       int             `json:"page"`
	NextPage   int             `json:"next_page"`
	TotalPages int             `json:"total_pages"`
}

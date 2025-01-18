package models

type Item struct {
    ShortDescription string `json:"shortDescription" binding:"required"`
    Price           string `json:"price" binding:"required"`
}

type Receipt struct {
    Retailer     string `json:"retailer" binding:"required"`
    PurchaseDate string `json:"purchaseDate" binding:"required"`
    PurchaseTime string `json:"purchaseTime" binding:"required"`
    Items        []Item `json:"items" binding:"required,min=1"`
    Total        string `json:"total" binding:"required"`
}

type ProcessResponse struct {
    ID string `json:"id"`
}

type PointsResponse struct {
    Points int64 `json:"points"`
}
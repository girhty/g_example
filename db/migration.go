package db

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dimensions struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
	Depth  float32 `json:"depth"`
}

type Review struct {
	ID            uint `gorm:"primaryKey"`
	Review_ID     uint
	Rating        int8      `json:"rating"`
	Comment       string    `json:"comment"`
	Date          time.Time `json:"date"`
	ReviewerName  string    `json:"reviewerName"`
	ReviewerEmail string    `json:"reviewerEmail"`
}

type Meta struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	Barcode   string    `json:"barcode"`
	QrCode    string    `json:"qrCode"`
}

type Image struct {
	gorm.Model
	Image_id int16
	Image    string
}
type Tag struct {
	gorm.Model
	Tag_id int16
	Tag    string
}

type Product struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Category             string     `json:"category"`
	Price                float32    `json:"price"`
	DiscountPercentage   float32    `json:"discountPercentage"`
	Rating               float32    `json:"rating"`
	Brand                string     `json:"brand"`
	Sku                  string     `json:"sku"`
	Weight               int16      `json:"weight"`
	Dimensions           Dimensions `gorm:"embedded" json:"dimensions"`
	WarrantyInformation  string     `json:"warrantyInformation"`
	ShippingInformation  string     `json:"shippingInformation"`
	AvailablityStatus    string     `json:"availablityStatus"`
	Reviews              []Review   `gorm:"foreignKey:Review_ID" json:"reviews"`
	ReturnPolicy         string     `json:"returnPolicy"`
	MinimumorderQuantity int16      `json:"minimumorderQuantity"`
	Meta                 Meta       `gorm:"embedded" json:"meta"`
	Thumbnail            string     `json:"thumbnail"`
}

type ProductsContainer struct {
	Items []Product `json:"products"`
	Total int16     `json:"total"`
	Skip  int16     `json:"skip"`
	Limit int16     `json:"limit"`
}

func Create_db() error {
	if req, err := http.Get("https://dummyjson.com/products"); err != nil {
		log.Fatal(err.Error())
	} else {
		body, b_err := io.ReadAll(req.Body)
		if b_err != nil {
			log.Fatal(b_err.Error())
			return b_err
		}
		data := &ProductsContainer{}
		j_err := json.Unmarshal(body, data)
		if j_err != nil {
			log.Fatal(j_err.Error())
			return j_err
		}

		sql_db, db_conn__err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if db_conn__err != nil {
			log.Fatalf("failed to connect database : %s\n", db_conn__err.Error())
			return db_conn__err
		}

		migrate_error := sql_db.AutoMigrate(&Product{}, &Review{})
		if migrate_error != nil {
			log.Fatal(migrate_error.Error())
			return migrate_error
		}
		ctx := context.Background()
		for _, item := range data.Items {
			for _, review := range item.Reviews {
				sql_err := gorm.G[Review](sql_db).Create(ctx, &Review{Review_ID: item.ID,
					Rating:        review.Rating,
					Comment:       review.Comment,
					Date:          review.Date,
					ReviewerName:  review.ReviewerName,
					ReviewerEmail: review.ReviewerEmail})
				if sql_err != nil {
					return sql_err
				}
			}
			sql_err := gorm.G[Product](sql_db).Create(ctx, &item)
			if sql_err != nil {
				log.Fatal(err)
				return sql_err
			}
		}
	}
	return nil
}

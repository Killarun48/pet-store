package models

type Category struct {
	ID   int    `json:"id" example:"4"`
	Name string `json:"name" example:"rabbit"`
}

type Tag struct {
	ID   int    `json:"id" example:"3"`
	Name string `json:"name" example:"gift"`
}

type Pet struct {
	ID        int      `json:"id" example:"1"`
	Category  Category `json:"category"`
	Name      string   `json:"name" example:"Daisy"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status" example:"available"`
}

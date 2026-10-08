package depop

type SearchResponse struct {
	Products []Product `json:"products"`
}

type Product struct {
	ID           int           `json:"id"`
	Slug         string        `json:"slug"`
	Title        *string       `json:"title,omitempty"`
	// User mentioned slugUsername, could be inside another struct or top level
	SlugUsername string        `json:"slugUsername,omitempty"`
	// Sometimes seller is an object
	Seller       Seller        `json:"seller,omitempty"`
	Price        Price         `json:"price"`
	Sizes        []string      `json:"sizes,omitempty"`
	// User said images[0][0].url. Let's try both Pictures and Images
	Pictures     [][]Image     `json:"pictures,omitempty"`
	Images       []Image       `json:"images,omitempty"`
}

type Seller struct {
    Username string `json:"username,omitempty"`
    Slug     string `json:"slug,omitempty"`
}

type Price struct {
	PriceString  string `json:"priceString,omitempty"`
	Amount       string `json:"amount,omitempty"`
	CurrencyCode string `json:"currencyCode,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
}

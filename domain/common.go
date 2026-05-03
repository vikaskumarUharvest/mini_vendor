package domain

type PaginatedResponse struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Pages int         `json:"pages"`
}

/*
Using interface{} there is basically Go’s way of saying: “this field can hold anything.”
/users → []User
/products → []Product
/orders → []Order
Instead of creating a new struct for each case, interface{} lets you reuse one response type:
*/

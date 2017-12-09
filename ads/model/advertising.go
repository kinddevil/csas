package model

type Advertising struct {
	// imageName, imageLink, schoolIds, province, city, title string, displayPages int
	Id           string `json:"id"`
	ImageName    string `json:"imageName"`
	ImageLink    string `json:"imageLink"`
	SchoolIds    string `json:"schoolIds"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Title        string `json:"title"`
	DisplayPages int    `json:"displayPages"`
}

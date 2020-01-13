package main

import (
	"fmt"
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	var pages string

	c := colly.NewCollector(
		colly.Async(true),
	)

	fileName := "craigslist.csv"
	file, err := os.Create(fileName)
	
	if err != nil {
		log.Fatalf("Could not create %s", fileName)
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Product Name", "Price"})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Current", r.URL)
	})

	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
		Parallelism: 4,
	})
	
	c.OnHTML("form.search-form", func(e *colly.HTMLElement) {
		pages = e.ChildText("span.totalcount")
		fmt.Printf("Count: %s \n", pages)
		e.ForEach("li.result-row", func(_ int, e *colly.HTMLElement) {
			var price, title string

			price = e.ChildText("span.result-price")
			title = e.ChildText("a.result-title.hdrlnk")
		
			new := strings.Split(price, "$")

			writer.Write([]string{
				title,
				new[1],
			})

			
			fmt.Printf("Product Name: %s \nPrice %s \n", title, new[1])
		})
	})


	
	c.Visit("https://vancouver.craigslist.org/d/video-gaming/search/vga")
	// for i := 120; i <= count; i=i+120 {
	// 	otherURLS := fmt.Sprintf("https://vancouver.craigslist.org/search/vga?s=%d", i)
	// 	c.Visit(fullURL)
	// }
	
}

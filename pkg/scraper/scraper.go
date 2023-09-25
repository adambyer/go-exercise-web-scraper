package scraper

import (
	"fmt"
	"net/http"
	"sync"
	"web-scraper/internal/models"

	"github.com/PuerkitoBio/goquery"
)

type FetchedArticles struct {
	Url string
	Articles []models.Article
}

func FetchArticles(urls []string) (map[string][]models.Article, error) {
	// Create a channel to receive articles.
	articleChannel := make(chan FetchedArticles, len(urls))

	// Create a wait group to wait for all Goroutines to finish.
	var waitGroup sync.WaitGroup

	for _, url := range urls {
		// Increment the wait group counter
		waitGroup.Add(1)

		// Start a Goroutine to fetch articles.
		go func(u string) {
			defer waitGroup.Done()

			articles, error := FetchArticlesFromUrl(u)
			if error != nil {
				fmt.Print("Error fetching articles")
			}

			articleChannel <- FetchedArticles{Url: u, Articles: articles}
		}(url)
	}

	go func() {
		waitGroup.Wait()
		close(articleChannel)
	}()

	// Declare the output structure.
	output := make(map[string][]models.Article)

	// Collect the articles from the channel.
	for fetchedArticles := range articleChannel {
		output[fetchedArticles.Url] = fetchedArticles.Articles
	}

	return output, nil
}

func FetchArticlesFromUrl(url string) ([]models.Article, error) {
	var articles []models.Article

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching URL: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error creating document: %w", err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		link, exists := s.Attr("href")

		if exists {
			articles = append(articles, models.Article{
				Title: "Title12: " + title,
				Link:  link,
			})
		}
	})

	return articles, nil
}

package alerter

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"depop-alerter/config"
	"depop-alerter/depop"
	"depop-alerter/notifier"
)

type Alerter struct {
	cfg      *config.Config
	client   *depop.Client
	notifier *notifier.TelegramNotifier
	cache    *Cache
	wg       sync.WaitGroup
	quit     chan struct{}
}

func NewAlerter(cfg *config.Config) (*Alerter, error) {
	client, err := depop.NewClient(cfg.Proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to init depop client: %w", err)
	}

	tgNotifier, err := notifier.NewTelegramNotifier(cfg.TelegramToken, cfg.ChatID)
	if err != nil {
		return nil, fmt.Errorf("failed to init telegram notifier: %w", err)
	}

	return &Alerter{
		cfg:      cfg,
		client:   client,
		notifier: tgNotifier,
		cache:    NewCache(500), // Max 500 items per query
		quit:     make(chan struct{}),
	}, nil
}

func (a *Alerter) Start() {
	log.Println("Starting Depop Alerter...")
	
	// Create a random seed
	rand.Seed(time.Now().UnixNano())

	for _, query := range a.cfg.Queries {
		a.wg.Add(1)
		go a.worker(query)
	}
}

func (a *Alerter) Stop() {
	close(a.quit)
	a.wg.Wait()
	log.Println("Alerter stopped gracefully.")
}

func (a *Alerter) worker(query string) {
	defer a.wg.Done()
	
	log.Printf("[Worker - %s] Started", query)

	for {
		select {
		case <-a.quit:
			log.Printf("[Worker - %s] Stopped", query)
			return
		default:
			err := a.checkNewItems(query)
			if err != nil {
				log.Printf("[Worker - %s] Error: %v", query, err)
				
				// Basic error handling for 403 / 429
				if err.Error() == "403" || err.Error() == "429" || fmt.Sprintf("%v", err) != "" {
					// We could check if it contains 403 or 429
					if containsRateLimitError(err) {
						errMsg := fmt.Sprintf("Rate limited or Forbidden (403/429) for query '%s'. Sleeping for 5 minutes...", query)
						log.Println(errMsg)
						a.notifier.NotifyError(errMsg)
						time.Sleep(5 * time.Minute)
						continue
					}
				}
			}

			// Sleep for interval + jitter (-2 to +5 seconds)
			jitter := rand.Intn(8) - 2
			sleepDuration := time.Duration(a.cfg.IntervalSeconds + jitter) * time.Second
			if sleepDuration < time.Second {
				sleepDuration = time.Second // minimum 1 second
			}
			
			time.Sleep(sleepDuration)
		}
	}
}

func (a *Alerter) checkNewItems(query string) error {
	products, err := a.client.FetchNewProducts(query)
	if err != nil {
		return err
	}

	isColdStart := a.cache.IsColdStart(query)
	
	newItemsFound := 0

	// We iterate backwards to process the oldest of the 'newlyListed' first
	// (usually index 0 is the newest, index 19 is the oldest of the batch)
	for i := len(products) - 1; i >= 0; i-- {
		p := products[i]
		
		isNew := a.cache.IsNew(query, p.ID)
		if isNew {
			if !isColdStart {
				newItemsFound++
				log.Printf("[Worker - %s] Found NEW item: %s (ID: %d)", query, p.Slug, p.ID)
				
				err := a.notifier.NotifyNewProduct(query, p)
				if err != nil {
					log.Printf("[Worker - %s] Failed to send notification for item %d: %v", query, p.ID, err)
				}
			}
		}
	}

	if isColdStart {
		log.Printf("[Worker - %s] Cold start complete. Cached %d initial items.", query, len(products))
	} else if newItemsFound > 0 {
		log.Printf("[Worker - %s] Cycle complete. Processed %d new items.", query, newItemsFound)
	}

	return nil
}

func containsRateLimitError(err error) bool {
	errMsg := err.Error()
	// Check for 403 or 429 inside the error string
	// TLS client sometimes returns "unexpected status code 403"
	// Or we can just use strings.Contains
	importStrings := false
	if importStrings {
		// we don't have strings imported in this snippet, let's just do a basic check
	}
	// We will implement a quick string contains
	for i := 0; i < len(errMsg)-2; i++ {
		if errMsg[i:i+3] == "403" || errMsg[i:i+3] == "429" {
			return true
		}
	}
	return false
}

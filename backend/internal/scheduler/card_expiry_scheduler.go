package scheduler

import (
	"fmt"
	"gym-admin/internal/service"
	"log"
	"time"
)

// CardExpiryScheduler handles card expiry checks and notifications
type CardExpiryScheduler struct {
	cardService *service.CardService
	ticker      *time.Ticker
	stopChan    chan bool
}

func NewCardExpiryScheduler() *CardExpiryScheduler {
	return &CardExpiryScheduler{
		cardService: service.NewCardService(),
		stopChan:    make(chan bool),
	}
}

// Start starts the scheduler
func (s *CardExpiryScheduler) Start() {
	go func() {
		// Run immediately on start
		s.runTasks()

		// Calculate time until next 9:00 AM
		now := time.Now()
		next9AM := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		if now.After(next9AM) {
			next9AM = next9AM.Add(24 * time.Hour)
		}

		// Wait until 9:00 AM
		time.Sleep(time.Until(next9AM))

		// Run every day at 9:00 AM
		s.ticker = time.NewTicker(24 * time.Hour)

		for {
			select {
			case <-s.ticker.C:
				s.runTasks()
			case <-s.stopChan:
				return
			}
		}
	}()

	log.Println("Card expiry scheduler started")
}

// Stop stops the scheduler
func (s *CardExpiryScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.stopChan <- true
	log.Println("Card expiry scheduler stopped")
}

// runTasks runs all scheduled tasks
func (s *CardExpiryScheduler) runTasks() {
	log.Println("Running card expiry tasks...")

	// Task 1: Update expired cards status
	if count, err := s.cardService.CheckAndUpdateExpiredCards(); err != nil {
		log.Printf("Error updating expired cards: %v", err)
	} else {
		log.Printf("Updated %d expired cards", count)
	}

	// Task 2: Send 7-day expiry notifications
	if count, err := s.cardService.SendExpiryNotifications(7); err != nil {
		log.Printf("Error sending 7-day expiry notifications: %v", err)
	} else {
		log.Printf("Sent %d expiry notifications (7 days)", count)
	}

	// Task 3: Send 3-day expiry notifications
	if count, err := s.cardService.SendExpiryNotifications(3); err != nil {
		log.Printf("Error sending 3-day expiry notifications: %v", err)
	} else {
		log.Printf("Sent %d expiry notifications (3 days)", count)
	}

	// Task 4: Send 1-day expiry notifications
	if count, err := s.cardService.SendExpiryNotifications(1); err != nil {
		log.Printf("Error sending 1-day expiry notifications: %v", err)
	} else {
		log.Printf("Sent %d expiry notifications (1 day)", count)
	}

	log.Println("Card expiry tasks completed")
}

// RunNow runs tasks immediately (for testing or manual trigger)
func (s *CardExpiryScheduler) RunNow() {
	s.runTasks()
}

// GetStatus returns scheduler status
func (s *CardExpiryScheduler) GetStatus() string {
	if s.ticker != nil {
		return fmt.Sprintf("Running - Next check at %s", time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"))
	}
	return "Stopped"
}

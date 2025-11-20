package scheduler

import (
"fmt"
"gym-admin/internal/service"
"time"

"github.com/robfig/cron/v3"
)

type PerformanceScheduler struct {
	cron           *cron.Cron
	coachService *service.CoachService
	perfService    *service.CoachPerformanceService
}

func NewPerformanceScheduler() *PerformanceScheduler {
	return &PerformanceScheduler{
		cron:           cron.New(cron.WithSeconds()),
		coachService: service.NewCoachService(),
		perfService:    service.NewCoachPerformanceService(),
	}
}

// Start begins the cron job
func (s *PerformanceScheduler) Start() {
	// Schedule to run at 00:00 on the 1st of every month
	s.cron.AddFunc("0 0 0 1 * *", s.updateAllCoachesPerformance)
	s.cron.Start()
	fmt.Println("Performance scheduler started. Will run on the 1st of every month.")
}

// Stop ends the cron job
func (s *PerformanceScheduler) Stop() {
	s.cron.Stop()
	fmt.Println("Performance scheduler stopped.")
}

// updateAllCoachesPerformance calculates and updates performance for all coaches for the previous month
func (s *PerformanceScheduler) updateAllCoachesPerformance() {
	fmt.Println("Running monthly coach performance update...")

	coaches, err := s.coachService.GetAllCoaches()
	if err != nil {
		fmt.Printf("Error getting all coaches: %v\n", err)
		return
	}

	// Calculate for the previous month
	now := time.Now()
	firstDayOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := firstDayOfCurrentMonth.AddDate(0, -1, 0)
	year, month := lastMonth.Year(), int(lastMonth.Month())

	for _, coach := range coaches {
		if err := s.perfService.UpdateCoachPerformance(coach.ID, year, month); err != nil {
			fmt.Printf("Error updating performance for coach %d: %v\n", coach.ID, err)
		}
	}

	fmt.Println("Monthly coach performance update completed.")
}

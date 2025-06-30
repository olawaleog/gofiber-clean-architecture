package jobs

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/robfig/cron/v3"
	"time"
)

// Create a new cron job manager
func SetupCronJobs(
	transactionService *service.TransactionService,
	truckService *service.TruckService,
	// Add other services you might need
) *cron.Cron {
	// Create a new cron with the required configuration
	// Using CRON_TZ=UTC to ensure timezone consistency
	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))

	// Add a cron job that runs every day at midnight
	// Format: second minute hour day month weekday
	_, err := c.AddFunc("0 */3 * * * *", func() {
		common.Logger.Info("Running task every 5 minutes")
		ctx := context.Background()

		//get active truck
		activeTruck := (*truckService).GetActiveTruck(ctx)
		if activeTruck.Id == 0 {
			return
		}

		// Example: Check for pending transactions older than 24 hours
		err := (*transactionService).ProcessPendingTransactions(ctx, activeTruck.Id)
		if err != nil {
			common.Logger.Error("Error in daily transaction processing: " + err.Error())
		}

		// Additional scheduled tasks can be added here
	})

	if err != nil {
		common.Logger.Error("Failed to setup cron job: " + err.Error())
	}

	// Add another example job that runs every hour
	_, err = c.AddFunc("0 0 * * * *", func() {
		common.Logger.Info("Running hourly scheduled task")
		// Your hourly task logic here
	})

	if err != nil {
		common.Logger.Error("Failed to setup hourly cron job: " + err.Error())
	}

	return c
}

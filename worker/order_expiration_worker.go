package worker

import (
	"context"
	"go-ecommerce-api/service"
	"log"
	"time"
)

type OrderExpirationWorker struct {
	orderExpirationService *service.OrderExpirationService
}

func NewOrderExpirationWorker(orderExpirationService *service.OrderExpirationService) *OrderExpirationWorker {
	return &OrderExpirationWorker{
		orderExpirationService: orderExpirationService,
	}
}

func (w *OrderExpirationWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Checking Expired Order")
			expiredOrderId, err := w.orderExpirationService.GetExpiredOrderIds(ctx)
			if err != nil {
				log.Println("Failed to retrieved expired order", err)
				continue
			}
			for _, id := range expiredOrderId {
				err := w.orderExpirationService.ExpireOrder(ctx, id)
				if err != nil {
					log.Printf("Failed to expire order %d: %v", id, err)
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

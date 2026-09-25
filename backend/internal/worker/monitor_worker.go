package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/royandi/gowatch/backend/internal/model"
	monitorchecker "github.com/royandi/gowatch/backend/internal/monitor"
	"github.com/royandi/gowatch/backend/internal/repository"
)

type MonitorWorker struct {
	monitorRepository *repository.MonitorRepository
	resultRepository  *repository.MonitorResultRepository
	pollInterval      time.Duration
	maxConcurrent     int
}

func NewMonitorWorker(
	monitorRepository *repository.MonitorRepository,
	resultRepository *repository.MonitorResultRepository,
) *MonitorWorker {
	return &MonitorWorker{
		monitorRepository: monitorRepository,

		resultRepository: resultRepository,

		pollInterval: 5 * time.Second,

		maxConcurrent: 10,
	}
}

func (w *MonitorWorker) Start(
	ctx context.Context,
) {
	log.Println(
		"📡 Background monitor worker aktif",
	)

	w.run(ctx)

	ticker :=
		time.NewTicker(
			w.pollInterval,
		)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println(
				"🛑 Background monitor worker berhenti",
			)

			return

		case <-ticker.C:
			w.run(ctx)
		}
	}
}

func (w *MonitorWorker) run(
	ctx context.Context,
) {
	if ctx.Err() != nil {
		return
	}

	monitors, err :=
		w.monitorRepository.FindDueForCheck(
			ctx,
		)

	if err != nil {
		if ctx.Err() != nil {
			return
		}

		log.Printf(
			"❌ Worker gagal mengambil monitor: %v",
			err,
		)

		return
	}

	if len(monitors) == 0 {
		return
	}

	log.Printf(
		"🔎 Worker menemukan %d monitor yang perlu dicek",
		len(monitors),
	)

	var waitGroup sync.WaitGroup

	semaphore :=
		make(
			chan struct{},
			w.maxConcurrent,
		)

	for _, currentMonitor := range monitors {
		waitGroup.Add(1)

		go func(
			currentMonitor model.Monitor,
		) {
			defer waitGroup.Done()

			select {
			case semaphore <- struct{}{}:

				defer func() {
					<-semaphore
				}()

			case <-ctx.Done():
				return
			}

			w.checkMonitor(
				ctx,
				currentMonitor,
			)
		}(
			currentMonitor,
		)
	}

	waitGroup.Wait()
}

func (w *MonitorWorker) checkMonitor(
	ctx context.Context,
	currentMonitor model.Monitor,
) {
	if ctx.Err() != nil {
		return
	}

	log.Printf(
		"🌐 Mengecek #%d %s -> %s",
		currentMonitor.ID,
		currentMonitor.Name,
		currentMonitor.URL,
	)

	result :=
		monitorchecker.CheckWebsiteContext(
			ctx,
			currentMonitor.URL,
		)

	if ctx.Err() != nil {
		return
	}

	err :=
		w.resultRepository.Create(
			ctx,
			currentMonitor.ID,
			result,
		)

	if err != nil {
		if ctx.Err() != nil {
			return
		}

		log.Printf(
			"❌ Gagal menyimpan hasil monitor #%d: %v",
			currentMonitor.ID,
			err,
		)

		return
	}

	if result.Status == "UP" {
		log.Printf(
			"✅ #%d %s UP | HTTP %d | %d ms",
			currentMonitor.ID,
			currentMonitor.Name,
			result.HTTPStatus,
			result.ResponseTime,
		)

		return
	}

	if result.Error != "" {
		log.Printf(
			"❌ #%d %s DOWN | %s | %d ms",
			currentMonitor.ID,
			currentMonitor.Name,
			result.Error,
			result.ResponseTime,
		)

		return
	}

	log.Printf(
		"❌ #%d %s DOWN | HTTP %d | %d ms",
		currentMonitor.ID,
		currentMonitor.Name,
		result.HTTPStatus,
		result.ResponseTime,
	)
}

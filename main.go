package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"main.go/chef"
	"main.go/consts"
	"main.go/dish"
	"main.go/handlers"
	"main.go/server"
	"main.go/station"
	"main.go/statistics"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stations := make(map[consts.ClassStation]chef.StationCook)

	stations[consts.BarClass] = station.NewBarStation(3)
	stations[consts.FryerClass] = station.NewFryerStation(3)
	stations[consts.GrillClass] = station.NewGrillStation(3)
	stations[consts.MolecularClass] = station.NewMolecularStation(3)

	menu := dish.NewListAvailabilityDishes()

	stat := statistics.NewStatistics()

	ch := chef.NewChef(stations, menu, stat)
	handler := handlers.NewHTTPHandlers(ctx, ch, menu)
	serv := server.NewServer(handler)

	serv.Start(ctx)

}

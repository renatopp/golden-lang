package builder

import (
	"time"

	"github.com/renatopp/golden/internal/builder/build"
	"github.com/renatopp/golden/internal/helpers/logger"
)

type Builder2 struct {
	startTime time.Time
}

func NewBuilder2() *Builder2 {
	return &Builder2{}
}

func (b *Builder2) Build(opts build.Options) error {
	logger.Debug("[builder] starting building")

	b.startTime = time.Now()
	err := build.Build(opts)
	if err != nil {
		return err
	}
	logger.Info("Finished building in %s", time.Since(b.startTime))
	return nil
}

func (b *Builder2) Run(opts build.Options) error {
	// execute the run process
	return nil
}

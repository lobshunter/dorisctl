package step

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
)

// TODO: rewrite display code

var _ Step = &StepDisplay{}

var (
	spinnerColor           = "cyan"
	spinnerCharSets        = spinner.CharSets[9] // see: https://github.com/briandowns/spinner#available-character-sets
	spinnerRefreshInterval = 200 * time.Millisecond
)

type StepDisplay struct {
	prefix string
	step   Step
}

func NewStepDisplay(prefix string, step Step) *StepDisplay {
	return &StepDisplay{
		prefix: prefix,
		step:   step,
	}
}

func (s *StepDisplay) Name() string {
	return s.prefix + "- " + s.step.Name()
}

func (s *StepDisplay) StepName() string {
	return s.Name()
}

func (s *StepDisplay) Execute(ctx context.Context) error {
	bar := spinner.New(spinnerCharSets, spinnerRefreshInterval)
	_ = bar.Color(spinnerColor)
	bar.Prefix = s.prefix + " "
	bar.Start()

	err := s.step.Execute(ctx)

	if err != nil {
		bar.FinalMSG = s.prefix + " Failed\n"
	} else {
		bar.FinalMSG = s.prefix + " Done\n"
	}
	bar.Stop()

	return err
}

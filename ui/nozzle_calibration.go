package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var nozzleCalibrationPanelInstance *nozzleCalibrationPanel

type pointCoordinates struct {
	x float64
	y float64
	z float64
}

type nozzleCalibrationPanel struct {
	CommonPanel
	zCalibrationMode bool
	activeTool       int
	cPoint           pointCoordinates
	zOffset          float64
	labZOffsetLabel  *gtk.Label
}

func NozzleCalibrationPanel(ui *UI, parent Panel) Panel {
	if nozzleCalibrationPanelInstance == nil {
		m := &nozzleCalibrationPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 3
		m.cPoint = pointCoordinates{x: 20, y: 20, z: 0}
		m.initialize()

		nozzleCalibrationPanelInstance = m
	}

	return nozzleCalibrationPanelInstance
}

func (m *nozzleCalibrationPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createIncreaseOffsetAgainButton(), 1, 0, 1, 1)
	m.Grid().Attach(m.createIncreaseOffsetHalfButton(), 2, 0, 1, 1)
	m.Grid().Attach(m.createDecreaseOffsetHalfButton(), 3, 0, 1, 1)
	m.Grid().Attach(m.createDecreaseOffsetAgainButton(), 4, 0, 1, 1)

	m.Grid().Attach(m.createIncreaseOffsetButton(), 1, 1, 1, 1)
	m.Grid().Attach(m.createAcceptButton(), 2, 1, 1, 1)
	m.Grid().Attach(m.createAbortButton(), 3, 1, 1, 1)
	m.Grid().Attach(m.createDecreaseOffsetButton(), 4, 1, 1, 1)

	m.Grid().Attach(m.createAutoZCalibrationButton(), 1, 2, 1, 1)

	m.step = MustStepButton("move-step.svg",
		Step{"5mm", 5.0}, Step{"1mm", 1.0}, Step{"0.1mm", 0.1}, Step{"0.05mm", 0.05},
	)
	m.Grid().Attach(m.step, 2, 2, 1, 1)

}

func (m *nozzleCalibrationPanel) createAutoZCalibrationButton() gtk.IWidget {
	return MustButtonImageStyle("Auto Z Calibration", "z-calibration.svg", "color3", func() {
		if m.zCalibrationMode {
			return
		}

		if m.UI.Settings != nil && m.UI.Settings.GCodes.ProbeCalibrate != "" {
			cmd := &octoprint.CommandRequest{}
			cmd.Commands = []string{
				m.UI.Settings.GCodes.ProbeCalibrate,
			}
			if err := cmd.Do(m.UI.Printer); err != nil {
				Logger.Error(err)
			}
		} else {
			cmd := &octoprint.RunZOffsetCalibrationRequest{}
			if err := cmd.Do(m.UI.Printer); err != nil {
				Logger.Error(err)
			}
		}

	})
}

func (m *nozzleCalibrationPanel) createIncreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("+ value", "z-offset-increase.svg", func() {
		distance := m.step.Value().(float64)
		m.testZ(distance)
	})
}

func (m *nozzleCalibrationPanel) createIncreaseOffsetAgainButton() gtk.IWidget {
	return MustButtonImage("+ repeat", "z-offset-increase.svg", func() {
		m.testZ("++")
	})
}

func (m *nozzleCalibrationPanel) createIncreaseOffsetHalfButton() gtk.IWidget {
	return MustButtonImage("+ 1/2", "z-offset-increase.svg", func() {
		m.testZ("+")
	})
}

func (m *nozzleCalibrationPanel) createDecreaseOffsetButton() gtk.IWidget {
	return MustButtonImage("- value", "z-offset-decrease.svg", func() {
		distance := m.step.Value().(float64) * -1.0
		m.testZ(distance)
	})
}

func (m *nozzleCalibrationPanel) createDecreaseOffsetAgainButton() gtk.IWidget {
	return MustButtonImage("- repeat", "z-offset-decrease.svg", func() {
		m.testZ("--")
	})
}

func (m *nozzleCalibrationPanel) createDecreaseOffsetHalfButton() gtk.IWidget {
	return MustButtonImage("- 1/2", "z-offset-decrease.svg", func() {
		m.testZ("-")
	})
}

func (m *nozzleCalibrationPanel) createAcceptButton() gtk.IWidget {
	return MustButtonImage("Accept", "complete.svg", func() 
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"ACCEPT",
		}

		Logger.Info("Accept auto z-offset request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *nozzleCalibrationPanel) createAbortButton() gtk.IWidget {
	return MustButtonImage("Abort", "stop.svg", func() 
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			"ABORT",
		}

		Logger.Info("Abort auto z-offset request")
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})
}

func (m *nozzleCalibrationPanel) testZ(v string) {

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		fmt.Sprintf("TESTZ Z=%f", v),
	}
	if err := cmd.Do(m.UI.Printer); err != nil {
		Logger.Error(err)
	}
}

func (m *nozzleCalibrationPanel) command(gcode string) error {
	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{gcode}
	return cmd.Do(m.UI.Printer)
}

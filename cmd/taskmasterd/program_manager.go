package main

import (
	"fmt"
	"log"
	"sort"
	"syscall"
)

type ProgramManager struct {
	Programs        Programs
	ProgramTaskChan chan ProgramTask
}

func NewProgramManager() *ProgramManager {
	var programManager ProgramManager
	programManager.Init()
	return &programManager
}

func (programManager *ProgramManager) Init() {
	programManager.ProgramTaskChan = make(chan ProgramTask)

	go func() {
		for programTask := range programManager.ProgramTaskChan {
			if programTask.Action == ProgramTaskActionStart {
				programTask.Program.Start()
			} else if programTask.Action == ProgramTaskActionStop {
				programTask.Program.Stop()
			} else if programTask.Action == ProgramTaskActionRestart {
				programTask.Program.Restart()
			} else if programTask.Action == ProgramTaskActionKill && programTask.ProcessID != "" {
				process := programTask.Program.GetProcessById(programTask.ProcessID)
				if process != nil {
					process.Cmd.Process.Signal(syscall.SIGKILL)
				}
			}
		}
	}()
}

func (programManager *ProgramManager) GetProgramByName(name string) *Program {
	for _, program := range programManager.Programs {
		if program.Config.Name == name {
			return program
		}
	}
	return nil
}

func (programManager *ProgramManager) StartPrograms() {
	for _, program := range programManager.Programs {
		programManager.ProgramTaskChan <- ProgramTask{
			Action:  ProgramTaskActionStart,
			Program: program,
		}
	}
}

func (programManager *ProgramManager) StopPrograms() {
	for _, program := range programManager.Programs {
		programManager.ProgramTaskChan <- ProgramTask{
			Action:  ProgramTaskActionStop,
			Program: program,
		}
	}
}

func (programManager *ProgramManager) StartProgramByName(name string) error {
	program := programManager.GetProgramByName(name)
	if program == nil {
		return fmt.Errorf("Program not found: \"%s\"", name)
	}
	programManager.ProgramTaskChan <- ProgramTask{
		Action:  ProgramTaskActionStart,
		Program: program,
	}
	return nil
}

func (programManager *ProgramManager) StopProgramByName(name string) error {
	program := programManager.GetProgramByName(name)
	if program == nil {
		return fmt.Errorf("Program not found: \"%s\"", name)
	}
	programManager.ProgramTaskChan <- ProgramTask{
		Action:  ProgramTaskActionStop,
		Program: program,
	}
	return nil
}

func (programManager *ProgramManager) RestartProgramByName(name string) error {
	program := programManager.GetProgramByName(name)
	if program == nil {
		return fmt.Errorf("Program not found: \"%s\"", name)
	}
	return nil
}

func (programManager *ProgramManager) GetSortedPrograms() []*Program {
	programsKeys := []string{}

	for key, _ := range programManager.Programs {
		programsKeys = append(programsKeys, key)
	}

	sort.Strings(programsKeys)

	programs := []*Program{}
	for _, key := range programsKeys {
		program := programManager.GetProgramByName(key)
		if program != nil {
			programs = append(programs, programManager.GetProgramByName(key))
		} else {
			log.Panic("GetSortedPrograms(): program is nil")
		}
	}

	return programs
}

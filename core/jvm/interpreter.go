package jvm

import "fmt"
import "math"

import "github.com/ethereum/go-ethereum/core/jvm/instructions"
import "github.com/ethereum/go-ethereum/core/jvm/instructions/base"
import "github.com/ethereum/go-ethereum/core/jvm/rtda"

func interpret(thread *rtda.Thread, logInst bool, contract *Contract, evm *EVM) (uint64, error) {
	defer catchErr(thread)
    return loop(thread, logInst, contract, evm)
}

func catchErr(thread *rtda.Thread) {
	if r := recover(); r != nil {
		logFrames(thread)
        fmt.Printf("error executing instruction: %v\n", r)
		panic(r)
	}
}

func loop(thread *rtda.Thread, logInst bool, contract *Contract, evm *EVM) (uint64, error) {
    var gas uint64
    if contract != nil {
        gas = contract.Gas
    } else {
        gas = math.MaxUint64 // de facto unlimited
    }
	reader := &base.BytecodeReader{}
	for {
		frame := thread.CurrentFrame()
		pc := frame.NextPC()
		thread.SetPC(pc)

		// decode
		reader.Reset(frame.Method().Code(), pc)
		opcode := reader.ReadUint8()
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		frame.SetNextPC(reader.PC())

		if logInst {
			logInstruction(frame, inst)
		}

		// execute
        gasConsumed := inst.Execute(frame, gas, contract, evm)
        if gasConsumed > gas {
            for !thread.IsStackEmpty() {
                thread.PopFrame()
            }
            return 0, ErrOutOfGas   
        } else {
            gas -= gasConsumed
        }
        if thread.IsStackEmpty() {
            break
        }
		if f := thread.TopFrame(); f.Method()==nil {
           // have reached the bogus frame for holding return value 
			break
		}
	}
    return gas, nil
}

func logInstruction(frame *rtda.Frame, inst base.Instruction) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	pc := frame.Thread().PC()
	fmt.Printf("%v.%v() #%2d %T %v\n", className, methodName, pc, inst, inst)
}

func logFrames(thread *rtda.Thread) {
	for !thread.IsStackEmpty() {
		frame := thread.PopFrame()
		method := frame.Method()
		className := method.Class().Name()
		lineNum := method.GetLineNumber(frame.NextPC())
		fmt.Printf(">> line:%4d pc:%4d %v.%v%v \n",
			lineNum, frame.NextPC(), className, method.Name(), method.Descriptor())
	}
}

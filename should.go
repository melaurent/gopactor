package gopactor

import 	"github.com/melaurent/gopactor/assertions"

// These assertions are mostly self-explanatory,
// but it may be helpful to go through some examples
// which can be found in the documentation for the assertions package:
// https://godoc.org/github.com/meAmidos/gopactor/assertions
var (
	ShouldReceive          		= assertions.ShouldReceive
	ShouldReceiveType			= assertions.ShouldReceiveType
	ShouldReceiveFrom      		= assertions.ShouldReceiveFrom
	ShouldReceiveSomething 		= assertions.ShouldReceiveSomething
	ShouldReceiveN         		= assertions.ShouldReceiveN

	ShouldSend         	 		= assertions.ShouldSend
	ShouldSendType				= assertions.ShouldSendType
	ShouldSendTo        		= assertions.ShouldSendTo
	ShouldSendSomething 		= assertions.ShouldSendSomething
	ShouldSendN         		= assertions.ShouldSendN

	ShouldNotSendOrReceive 		= assertions.ShouldNotSendOrReceive

	ShouldStart              	= assertions.ShouldStart
	ShouldStop               	= assertions.ShouldStop
	ShouldBeRestarting       	= assertions.ShouldBeRestarting
	ShouldObserveTermination 	= assertions.ShouldObserveTermination

	ShouldSpawn 				= assertions.ShouldSpawn
)

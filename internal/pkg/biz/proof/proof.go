package proof

import (
	"fmt"

	"github.com/go-omnibus/proof"
)

func InitProof() {
	p := proof.New()
	p.SetDivision(proof.TimeDivision)
	p.SetTimeUnit(proof.Day)
	p.SetEncoding(proof.JSONEncoder)
	//p.CloseConsoleDisplay()

	p.SetInfoFile("./logs/application.log")
	p.SetErrorFile("./logs/application_err.log")

	p.Run()
	defer proof.Sync()

	fmt.Println("proof initialized")
}

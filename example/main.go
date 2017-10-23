package main

import (
	"github.com/AlexSnet/godgdat"
)

func main() {
	dg, err := godgdat.Open("/Users/alex/Downloads/DgdatToXlsx-master/download/Almetevsk-24.0.0.dgdat")
	// dg, err := godgdat.Open("/Users/alex/Downloads/2gis/3.0/Data_Moscow.dgdat")
	// dg, err := godgdat.Open("/Users/alex/Downloads/2GISData_Odessa~mobile-167.7.3.dgdat")

	defer dg.Close()

	if err != nil {
		panic(err)
	}
}

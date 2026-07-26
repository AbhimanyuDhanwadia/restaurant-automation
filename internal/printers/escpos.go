package printers

import "bytes"

type Formatter interface{ Format(Ticket) []byte }

type ESCPosFormatter struct{}

func (ESCPosFormatter) Format(ticket Ticket) []byte {
	var output bytes.Buffer
	output.Write([]byte{0x1b, 0x40})
	if ticket.Reprint {
		output.WriteString("*** REPRINT ***\n")
	}
	output.WriteString("ORDER: " + ticket.OrderID + "\n")
	output.WriteString("DESTINATION: " + ticket.Destination + "\n")
	output.WriteString("------------------------------\n")
	for _, line := range ticket.Lines {
		output.WriteString(line.Text)
		if line.Quantity > 1 {
			output.WriteString(" x ")
			output.WriteString(string(rune('0' + line.Quantity)))
		}
		output.WriteByte('\n')
	}
	output.WriteString("\n\n")
	output.Write([]byte{0x1d, 0x56, 0x00})
	return output.Bytes()
}

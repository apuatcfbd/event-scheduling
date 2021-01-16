package customevents

import "log"

// SendEmail sends the email
func SendEmail(data interface{}) {
	log.Println("📨 Sending email with data: ", data)
}

// PayBills pays the bills
func PayBills(data interface{}) {
	log.Println("💲 Pay me a bill: ", data)
}

// domain/coffee.go
// ชั้นในสุด กฎธุรกิจและโมเดล
// ไม่ import อะไรจากภายนอก domain เลย

package domain

// Coffee = แบบพิมพ์เขียวของกาแฟ 1 เมนู
type Coffee struct {
	ID    string
	Name  string
	Price float64
	Emoji string
}
package enum

import "database/sql/driver"

type Frequency string

const (
	Hourly Frequency = "hourly"
	Daily  Frequency = "daily"
)

func (f Frequency) IsValid() bool                { return f == Hourly || f == Daily }
func (f Frequency) Value() (driver.Value, error) { return string(f), nil }
func (f *Frequency) Scan(v any) error            { *f = Frequency(v.(string)); return nil }

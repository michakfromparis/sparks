package sys

import (
	log "github.com/Sirupsen/logrus"
)

func Init() {
	initLogger()
}

func initLogger() {
	log.SetFormatter(&Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "15:04:05",
		// FieldsOrder     []string // default: fields sorted alphabetically
		// TimestampFormat string   // default: time.StampMilli = "Jan _2 15:04:05.000"
		// HideKeys        bool     // show [fieldValue] instead of [fieldKey:fieldValue]
		// NoColors        bool     // disable colors
		// NoFieldsColors  bool     // color only level, default is level + fields
		// ShowFullLevel   bool     // true to show full level [WARNING] instead [WARN]
	})
}

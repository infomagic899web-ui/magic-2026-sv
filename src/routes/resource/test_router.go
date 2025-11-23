package resources

import (
	"context"
	"encoding/json"
	"log"
	"magic-server-2026/src/db"
	"magic-server-2026/src/middlewares"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

// CSPReport represents a CSP violation document
type CSPReport struct {
	BlockedURI        string                 `bson:"blocked_uri"`
	DocumentURI       string                 `bson:"document_uri"`
	OriginalPolicy    string                 `bson:"original_policy"`
	ViolatedDirective string                 `bson:"violated_directive"`
	SourceFile        string                 `bson:"source_file,omitempty"`
	LineNumber        int                    `bson:"line_number,omitempty"`
	ColumnNumber      int                    `bson:"column_number,omitempty"`
	Disposition       string                 `bson:"disposition,omitempty"`
	Referrer          string                 `bson:"referrer,omitempty"`
	ReportTime        time.Time              `bson:"report_time"`
	Extra             map[string]interface{} `bson:"extra,omitempty"`
}

func TestRouter(app fiber.Router) {
	// Other endpoints omitted for brevity

	app.Post("/enforce-csp-report", func(c fiber.Ctx) error {
		var body map[string]map[string]interface{}
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			log.Printf("[CSP] Invalid report: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		reportData, ok := body["csp-report"]
		if !ok {
			log.Println("[CSP] Missing 'csp-report' field")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		// Build CSPReport document
		doc := CSPReport{
			BlockedURI:        utils.Stringify(reportData["blocked-uri"]),
			DocumentURI:       utils.Stringify(reportData["document-uri"]),
			OriginalPolicy:    utils.Stringify(reportData["original-policy"]),
			ViolatedDirective: utils.Stringify(reportData["violated-directive"]),
			SourceFile:        utils.Stringify(reportData["source-file"]),
			LineNumber:        utils.Intify(reportData["line-number"]),
			ColumnNumber:      utils.Intify(reportData["column-number"]),
			Disposition:       utils.Stringify(reportData["disposition"]),
			Referrer:          utils.Stringify(reportData["referrer"]),
			ReportTime:        time.Now(),
			Extra:             reportData,
		}

		collection := db.Client.Database("magic899_db").Collection("csp_reports")
		_, err := collection.InsertOne(context.Background(), doc)
		if err != nil {
			log.Printf("[CSP] MongoDB insert failed: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Printf("[CSP] Violation stored: %v -> %v", doc.DocumentURI, doc.BlockedURI)
		return c.SendStatus(fiber.StatusNoContent)
	})

	app.Post("/post", middlewares.CSRFTokenMiddleware, func(x fiber.Ctx) error {
		return x.SendStatus(fiber.StatusOK)
	})
}

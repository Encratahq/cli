package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Encratahq/cli/internal/api"
	"github.com/Encratahq/cli/internal/output"
	"github.com/spf13/cobra"
)

var phoneCmd = &cobra.Command{
	Use:   "phone [number]",
	Short: "Look up a phone number",
	Long:  "Retrieve carrier, location, and validation info for a phone number.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Validate(); err != nil {
			return err
		}

		client := api.New(cfg.BaseURL, cfg.APIKey)
		data, err := client.PhoneSearch(args[0])
		if err != nil {
			output.Error(err.Error())
			return err
		}

		if cfg.Output == "json" {
			output.JSON(data)
			return nil
		}

		var result struct {
			Phone   string `json:"phone"`
			Valid   bool   `json:"valid"`
			Format  *struct {
				International string `json:"international"`
				Local         string `json:"local"`
			} `json:"format"`
			Country *struct {
				Code   string `json:"code"`
				Name   string `json:"name"`
				Prefix string `json:"prefix"`
			} `json:"country"`
			Location string  `json:"location"`
			Type     string  `json:"type"`
			Carrier  string  `json:"carrier"`
			Credits  float64 `json:"credits"`
		}

		if err := json.Unmarshal(data, &result); err != nil {
			output.JSON(data)
			return nil
		}

		output.Header("Phone Lookup: " + args[0])

		valid := output.Err.Sprint("✗ Invalid")
		if result.Valid {
			valid = output.Success.Sprint("✓ Valid")
		}

		international := ""
		local := ""
		if result.Format != nil {
			international = result.Format.International
			local = result.Format.Local
		}

		countryName := ""
		countryCode := ""
		if result.Country != nil {
			countryName = result.Country.Name
			countryCode = result.Country.Code
		}

		output.KV(
			"Number", result.Phone,
			"Valid", valid,
			"International", international,
			"Local", local,
			"Country", fmt.Sprintf("%s (%s)", countryName, countryCode),
			"Location", result.Location,
			"Type", result.Type,
			"Carrier", result.Carrier,
		)

		fmt.Println()
		output.Dim.Printf("  Credits used: %.0f\n", result.Credits)
		return nil
	},
}

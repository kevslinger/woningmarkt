package woningmarkt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func HuurVerhogingMain() int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://opendata.cbs.nl/ODataApi/odata/70675ned/TypedDataSet", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
		return 1
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v", err)
		return 1
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
		return 1
	}
	var ds HuurVerhogingDataset
	err = json.Unmarshal(body, &ds)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v", err)
		return 1
	}

	xys, err := MakeHuurVerhogingDataset(ds)
	if err != nil {
		fmt.Printf("Error generating dataset: %v", err)
		return 1
	}

	p := plot.New()
	p.Title.Text = "Huurverhogingen in Nederland, 1959-2024"
	p.X.Label.Text = "Jaar"
	p.X.Min, p.X.Max, p.Y.Min, p.Y.Max = plotter.XYRange(xys)
	p.Y.Label.Text = "Verhoging (%)"

	err = plotutil.AddLinePoints(p, xys)
	if err != nil {
		fmt.Printf("Error adding line points: %v", err)
		return 1
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		fmt.Printf("Error saving graph: %v", err)
		return 1
	}
	return 0
}

// MakeHuurVerhogingDataset converts the Dataset we marshalled from JSON into a set of (x, y) datapoints,
// where X is the year and Y is the percent increase of the rent in that year (compared to the previous year)
func MakeHuurVerhogingDataset(ds HuurVerhogingDataset) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(ds.Value))
	for idx, val := range ds.Value {
		jaar, err := strconv.Atoi(val.Perioden[:4])
		if err != nil {
			return pts, err
		}
		pts[idx].X = float64(jaar)
		pts[idx].Y = val.HuurVerhoging
	}
	return pts, nil
}

// HuurVerhogingDataset consists of a metadata URL and the datapoints for the rent increase dataset
type HuurVerhogingDataset struct {
	Metadata string                   `json:"odata.metadata"`
	Value    []HuurVerhogingDataPoint `json:"value"`
}

// HuurVerhogingDataPoint consists of an ID per data point as well as the period (year) and % increase
// in rent for that year (compared to the previous year)
type HuurVerhogingDataPoint struct {
	ID            int     `json:"ID"`
	Perioden      string  `json:"Perioden"`
	HuurVerhoging float64 `json:"Huurverhoging_1"`
}

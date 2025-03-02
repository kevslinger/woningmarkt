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

func PrijzenBestaandeKoopwoningenMain() int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://opendata.cbs.nl/ODataApi/odata/85773ned/TypedDataSet", nil)
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
	var ds PrijzenBestaandeKoopwoningenDataset
	err = json.Unmarshal(body, &ds)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v", err)
		return 1
	}

	xys, err := MakePrijzenBestaandeKoopwoningenDataset(ds)
	if err != nil {
		fmt.Printf("Error generating dataset: %v", err)
		return 1
	}

	p := plot.New()
	p.Title.Text = "Bestaande Koopwoning prijzen in Nederland, 1995-2024"
	p.X.Label.Text = "Jaar/Maand"
	p.X.Min, p.X.Max, p.Y.Min, p.Y.Max = plotter.XYRange(xys)
	p.Y.Label.Text = "Gemiddelde Prijs (€)"

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

// MakePrijzenBestaandeKoopwoningenDataset converts the data we marshalled from JSON into an (x, y) datapoint
// for each value. The X represents the year/month (JJJJMM) of the datapoint, and the Y represents the average sell price
func MakePrijzenBestaandeKoopwoningenDataset(ds PrijzenBestaandeKoopwoningenDataset) (plotter.XYs, error) {
	pts := make(plotter.XYs, len(ds.Value))
	for idx, val := range ds.Value {
		jaar, err := strconv.Atoi(val.Perioden[:4])
		if err != nil {
			return pts, err
		}
		maand, err := strconv.Atoi(val.Perioden[len(val.Perioden)-2:])
		if err != nil {
			return pts, err
		}
		pts[idx].X = float64(jaar)*100 + float64(maand)
		pts[idx].Y = val.GemiddeldeVerkoopprijs
		fmt.Println(val.GemiddeldeVerkoopprijs)
	}
	return pts, nil
}

// PrijzenBestaandeKoopwoningenDataset consists of a metadata URL as well as the datapoints
// from the dataset
type PrijzenBestaandeKoopwoningenDataset struct {
	Metadata string                              `json:"odata.metadata"`
	Value    []PrijzenBestaandeKoopwoningenPoint `json:"value"`
}

// PrijzenBestaandeKoopwoningenPoint contains all the information from the OpenData table for the prices
// of existing houses
type PrijzenBestaandeKoopwoningenPoint struct {
	ID                             int     `json:"ID"`
	Perioden                       string  `json:"Perioden"`
	PrijsIndex                     float64 `json:"PrijsindexVerkoopprijzen_1"`
	PrijsIndexTOVVoorgaandePeriode float64 `json:"OntwikkelingTOVVoorgaandePeriode_2"`
	PrijsIndexTOVJaarEerder        float64 `json:"OntwikkelingTOVEenJaarEerder_3"`
	VerkochteWoningen              int     `json:"VerkochteWoningen_4"`
	VerkochteTOVVoordgaandePeriod  float64 `json:"OntwikkelingTOVVoorgaandePeriode_5"`
	VerkochteTOVJaarEerder         float64 `json:"OntwikkelingTOVEenJaarEerder_6"`
	GemiddeldeVerkoopprijs         float64 `json:"GemiddeldeVerkoopprijs_7"`
	TotalWaardeVerkoopPrijzen      float64 `json:"TotaleWaardeVerkoopprijzen_8"`
}

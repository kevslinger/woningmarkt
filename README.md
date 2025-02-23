# Woning Markt Data Analyse

- [Woning Markt Data Analyse](#woning-markt-data-analyse)
  - [Gameplan](#gameplan)
    - [Phase 1 \<-- (We are currently here)](#phase-1----we-are-currently-here)
    - [Phase 2](#phase-2)
    - [Phase 3](#phase-3)


The goal is to teach myself some data science techniques in go (weird, right? Why not python, julia, or R?).

I will use [gonum/plot](https://github.com/gonum/plot) to create my graphs, and use data from [CBS](https://www.cbs.nl)

The project will go as follows:

1. Familiarize myself with plotting in go using several different types of charts (line charts, scatter plots, etc.)
2. Inspect trends in the data to see if any correlations exist
3. Create a predictive model that tries to predict the future house price change year over year
4. Train and test the model, and compare to actual bank and analyst predictions


The project will consist of the following structure:

`main`

`adapters/cbs`

`plot`

`cmd/woningmarkt`

## Gameplan

We outline the approach to this project in several phases

### Phase 1 <-- (We are currently here)

In the first phase, we will simply make an HTTP GET request to the [CBS OpenData](https://www.cbs.nl/en-gb/our-services/open-data) endpoint to get a dataset.
We will then process the dataset, plot the result, and save the plot to a PNG file

We will create several "scripts" in the `cmd` package, one for each separate dataset.

### Phase 2

In the next phase, we will extract some of the code we wrote to handle the CBS OpenData data and use it as a basis of the CBS `adapter`. This package will allow us to read and process data from CBS without needing to know much about the data.

We will perform the same extraction with the `plot` package by combining plotting code into a more generic package that can be used for different types of datasets.

Once this is implemented, we will update our scripts to use these new packages.

### Phase 3

In this final phase, we will use the CBS adapter and the plot package to help us as we attempt to build a predictive model of housing prices.

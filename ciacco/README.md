# Ciacco

The `ciacco` script execute all miners from the `CIACCO_MINER_DIRECTORY`.
The output of each miner is saved inside the `CIACCO_OUTPUT_DIR`.

## Configuration file

The configuration file is `/etc/sysconfig/dante`.

### Default values

- `CIACCO_MINER_DIRECTORY`: `/usr/share/dante/miners`
- `CIACCO_OUTPUT_DIRECTORY` `/var/lib/nethserver/dante`

## Miners

A miner is a script which collects statistics from the system for the current day.

On success, each miner must always generates output in JSON format and should exit 0.
On error the miner must exit with non-zero code and can be generate custom output.

The name of the script must respect some naming conventions. The name is composed by 3 parts: `minerId-type-subtype`.
Both `minerId` and `type` are mandatory, while `subtype` is optional.

Valid naming examples:

- df-chart-pie
- hostname-label


## Miners output file

Ciacco defines the structure of the output file of a miner.
The file produced by a miner can represent different kinds of graphical widgets:

- a label
- a card that displays a title and a single numerical value (counter)
- a chart
- a table
- a ranking list

Each miner has a specific set of fields, but there are some common ones which are mandatory:

- `minerId`: miner identifier
- `type`: it describes the widget type. Valid values are:
   - `label`
   - `counter`
   - `chart`
   - `table`
   - `list`

Extra commond fields:

- `unit`: the data type (not valid for `label` type), valid types are:
  - `number`
  - `seconds`
  - `bytes`
  - `time`
- `aggregationType`: how data will aggreagate by server for the selected span of time:
   - `sum`: values will be summed
   - `average`: average calculation
   - `snapshot`: values will not be aggregated


### label

Display a simple label, labels have no trends and can't be aggregated.

Extra fields:

- `value`: a value in string format

#### Example

```json
{
  "value": "QEMU Standard PC (i440FX + PIIX, 1996)",
  "type": "label",
  "minerId": "hardware-label"
}
```

### counter

Display a counter with a trend and a line chart.

Extra fields:

- `value`: numerical value of the counter
- `trendType`: display the counter variation
  - `percentage`
  - `number`

#### Example

```json
{
  "value": 43,
  "type": "counter",
  "minerId": "mailsent-counter",
  "aggregationType": "sum",
  "unit": "number",
  "trendType": "percentage"
}
```

### chart

Display a chart without a trend. Extra fields:

- `chartType`: "pie" or "bar" or "line" or "area" or "column"
- `categories`: array of categories (for pie, bar and columns charts) or values on the x-axis (for line and area charts)
- `series`: array of values/objects associated to the categories
    - `name`: name of the serie
    - `data`: array of values
    - `unit`: see above


#### Example 1: series

```json
{
  "aggregationType": "snapshot",
  "categories": [
    "192.168.1.252",
    "192.168.1.253",
    ...
  ],
  "type": "chart",
  "title": "hosttraffic-chart-column",
  "minerId": "hosttraffic-chart-column",
  "series": [
    {
      "unit": "bytes",
      "name": "sent",
      "data": [
        132710324944,
        12560210471,
        ...
      ]
    },
    {
      "unit": "bytes",
      "name": "received",
      "data": [
        104958008699,
        7118435580,
        ...
      ]
    }
  ],
  "chartType": "column"
}

```

### Example 2: pie

```json
{
  "type": "chart",
  "chartType": "pie",
  "minerId": "mailfilter-chart-pie",
  "aggregationType": "snapshot",
  "unit": "number",
  "series": [
    {
      "name": "mails",
      "data": [
        2,
        143
      ]
    }
  ],
  "categories": [
    "virus",
    "spam"
  ]
}
```

### table

Display a simple table.

Extra fields:

- `columnsHeader`: column headers
- `rowHeader`: row headers
- `rows`: array of array values for the table

#### Example

```json
{
  "rows": [
    [
      132710712601,
      104958634973
    ],
    [
      12560714747,
      7118848245
    ],
    ...
  ],
  "rowHeader": [
    "host0",
    "host1",
    ...
  ],
  "minerId": "hosttraffic-table",
  "aggregationType": "snapshot",
  "unit": "bytes",
  "title": "hosttraffic-table",
  "type": "table",
  "columnHeader": [
    "sent",
    "received"
  ]
}

```

### list

Display an ordered list of items.
The server will aggregate data from all days and calculate the "top X items" or "bottom X items".

Extra fields:

- `data`: a list of objects with a label and a counter

#### Example

```json
{
  "aggregationType": "sum",
  "unit": "bytes",
  "data": [
    {
      "count": 1121576493,
      "name": "host1.nethesis.it"
    },
    {
      "count": 487075598,
      "name": "host2.nethesis.it"
    },
    ...
  ],
  "type": "list",
  "title": "squidsources-list",
  "minderId": "squidsources-list"
}

```


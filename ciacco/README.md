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
- `title`: a i18n string representing a description of the information associated to the label
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
    "type": "label",
    "title": "hostname-label",
    "minerId": "miner-total-emails-received",
    "value": "mail.nethserver.org"
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
    "type": "counter",
    "title": "received_mails",
    "minerId": "miner-total-emails-received",
    "value":  42,
    "unit": "number",
    "aggregationType": "snapshot",
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


#### Example 1: single series

```json
{
    "type": "chart",
    "chartType": "pie",
    "title": "traffic_by_protocol",
    "minerId": "miner-traffic-by-protocol", 
    "categories": [ "TCP", "UDP", "ICMP", "Other" ],
    "series": [
        {
            "name": "traffic",
            "data": [ 34, 42, 45, 38 ]
        }
    ],
    "aggregationType": "sum"
}
```

#### Example 2: multiple series

```json
{
    "type": "chart",
    "chartType": "bar",
    "title": "Calls sent and received by hour",
    "minerId": "miner-calls-sent-and-received-by-hour", 
    "categories": [ "9:00-10:00", "10:00-11:00", "11:00-12:00", "12:00-13:00" ],
    "series": [
        {
            "name": "Calls sent",
            "data": [ 34, 42, 45, 38 ]
        },
        {
            "name": "Calls received",
            "data": [ 24, 33, 35, 28 ]
        }
    ],
    "aggregationType": "average"
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
    "type": "table",
    "title": "host-traffic",
    "minerId": "miner-host-traffic",
    "unit": "bytes",
    "columnHeader": [ "Total", "Sent", "Received" ],
    "rowHeader": "true | false // if true first item of columnHeader is blank: [ '', 'Sent', 'Received' ]",
    "rows":  [ 
        [ 720, 400, 320 ],
        [ 550, 300, 250 ]
    ],
    "aggregationType": "snapshot"
}
```

### list

Display an ordered list of items.
The server will aggregate data from all days and calculate the "top X items" or "bottom X items".

#### Example

```json
{
    "type": "list",
    "title": "blockedcategories",
    "minerId": "blockedcategories",
    "unit": "number",
    "data":  [
        {
            "count": 865,
            "name": "88.60.zz.xx"
        },
        {
            "count": 272,
            "name": "49.248.yyy.xxx"
        },
    ],
}
```


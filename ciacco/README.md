# Ciacco

The `ciacco` script execute all miners from the `CIACCO_MINER_DIRECTORY`.
The output of each miner is saved inside the `CIACCO_OUTPUT_DIR`.

## Configuration file

The configuration file is `/etc/sysconfig/dante`.

### Default values

- `CIACCO_MINER_DIRECTORY`: `/usr/share/dante/miners`
- `CIACCO_OUTPUT_DIRECTORY` `/var/lib/nethserver/dante`

## Output file

Ciacco defines the structure of the output file of a miner.
The file produced by a miner can represent different kinds of graphical widgets:

- a card that displays a title and a single numerical value (counter)
- a chart
- a table

### Label

This is the structure of the output JSON of a miner that feeds a label widget:

- `type`: "label"
- `minerId`: miner identifier
- `title`: a i18n string representing a description of the information associated to the label
- `value`: a value in string format
- `snapshot`:  can be `true` or `false`. If set to `true`, values will not be aggregated over a span of time. It's fixed to `true` for the `label` type

#### Example

```json
{
    "type": "label",
    "title": "hostname-label",
    "minerId": "miner-total-emails-received",
    "value": "mail.nethserver.org",
    "snapshot": true
}
```

### Counter

This is the structure of the output JSON of a miner that feeds a counter widget:

- `type`: "counter"
- `minerId`: miner identifier
- `title`: a i18n string representing a description of the information associated to the counter
- `value`: numerical value of the counter
- `snapshot`: can be `true` or `false`. If set to `true`, values will not be aggregated over a span of time

#### Example

```json
{
    "type": "counter",
    "title": "received_mails",
    "minerId": "miner-total-emails-received",
    "value":  42,
    "snapshot": true
}
```

### Chart

This is the structure of the output JSON of a miner that feeds a chart widget:

- `type`: "chart"
- `chartType`: "pie" or "bar" or "line" or "area" or "column"
- `title`: the title of the chart
- `minerId`: miner identifier
- `unit`: unit of measurement. Currently only "bytes" is supported
- `categories`: array of categories (for pie, bar and columns charts) or values on the x-axis (for line and area charts)
- `series`: array of values/objects associated to the categories
    - `i18n`: boolean value that specifies if category and series labels should be translated or not (e.g. numbers, IP addresses, ...)
    - `name`: name of the series, it will be translated if `i18n` is set to true
    - `data`: array of values
- `snapshot`: can be `true` or `false`. If set to `true`, values will not be aggregated over a span of time


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
            "i18n": true,
            "data": [ 34, 42, 45, 38 ]
        }
    ],
    "snapshot": false
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
            "i18n": true,
            "data": [ 34, 42, 45, 38 ]
        },
        {
            "name": "Calls received",
            "i18n": true,
            "data": [ 24, 33, 35, 28 ]
        }
    ],
    "snapshot": false
}
```

### Table

This is the the structure of output JSON of a miner that feeds a table widget:

- `type`: "table"
- `title`: the title of the table
- `minerId`: miner identifier
- `unit`: unit of measurement. Currently only "bytes" is supported
- `columnsHeader`: column headers
- `rowHeader`: row headers
- `rows`: array of array values for the table
- `snapshot`: can be `true` or `false`. If set to `true`, values will not be aggregated over a span of time

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
    "snapshot": true
}
```

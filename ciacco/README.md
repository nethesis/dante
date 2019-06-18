# Ciacco

The `ciacco` script execute all miners from the `CIACCO_MINER_DIRECTORY`.
The output of each miner is saved inside the `CIACCO_OUTPUT_DIR`.

## Configuration file

The configuration file is `/etc/sysconfig/dante`.

## Default values

- CIACCO_MINER_DIRECTORY=/usr/share/dante/miners
- CIACCO_OUTPUT_DIRECTORY=/var/lib/nethserver/dante

## Output file

Ciacco defines the structure of the output file of a miner.
The file produced by a miner can represent different kinds of graphical widgets:

- a card that displays a title and a single numerical value (counter)
- a chart
- a table

### Counter

This is the structure of the output JSON of a miner that feeds a counter widget:

- type: "counter"
- minerID: miner identifier
- title: a description of the information associated to the counter
- value: numerical value of the counter
- tags: list of keywords related to the information displayed by the widget
- position: integer number that defines the position of the widget on the dashboard. Small numbers are displayed top-left.

#### Example

```
{
    type: "counter",
    title: "Total e-mails received",
    minerID: "miner-total-emails-received",
    value:  42,
    tags: [ "email", "e-mail", "received", "total", "mail" ],
    position: 10
}
```

### Chart

This is the structure of the output JSON of a miner that feeds a chart widget:

- type: "chart"
- chartType: "pie" or "bar" or "line" or "area" or "column"
- title: the title of the chart
- minerID: miner identifier
- position: integer number that defines the position of the widget on the dashboard. Small numbers are displayed top-left.
- tags: list of keywords related to the information displayed by the widget
- i18n: boolean value that specifies if category and series labels should be translated or not (e.g. numbers, IP addresses, ...)
- unit (? dataType): unit of measurement. Currently only "bytes" is supported
- categories: array of categories (for pie, bar and columns charts) or values on the x-axis (for line and area charts)
- series: array of values/objects associated to the categories


#### Example 1: single series

```
{
    type: "chart",
    chartType: "pie",
    title: "Traffic by protocol",
    minerID: "miner-traffic-by-protocol", 
    position: 20,
    tags: [ "traffic", "protocol", "tcp", "udp", "icmp", "network" ],
    i18n: true,
    categories: [ "TCP", "UDP", "ICMP", "Other" ],
    series: [ 1204, 767, 32, 184 ]
}
```

#### Example 2: multiple series

```
{
    type: "chart",
    chartType: "bar",
    title: "Calls sent and received by hour",
    minerID: "miner-calls-sent-and-received-by-hour", 
    position: 30,
    tags: [ "calls", "call", "sent", "received", "hour", "phone" ],
    i18n: true,
    categories: [ "9:00-10:00", "10:00-11:00", "11:00-12:00", "12:00-13:00" ],
    series: [
        {
            name: 'Calls sent',
            i18n: true,
            data: [ 34, 42, 45, 38 ]
        },
        {
            name: 'Calls received',
            i18n: true,
            data: [ 24, 33, 35, 28 ]
        }
    ]
```

### Table

This is the the structure of output JSON of a miner that feeds a table widget:

- type: "table"
- title: the title of the table
- minerID: miner identifier
- position: integer number that defines the position of the widget on the dashboard. Small numbers are displayed top-left.
- tags: list of keywords related to the information displayed by the widget
- i18nColumns: boolean value that specifies if columns headers should be translated or not (e.g. numbers, IP addresses, ...).
- i18nRows: boolean value that specifies if row headers should be translated or not.
- unit (? dataType): unit of measurement. Currently only "bytes" is supported
- columns: column headers
- rows: array of values for the table

#### Example

```
{
    type: "table",
    title: "Host traffic",
    minerID: "miner-host-traffic",
    position: 40,
    tags: [ "host", "traffic", "total", "received", "sent" ],
    i18nColumns: true,
    i18nRows: false,
    unit: "bytes",
    colums: [ "Total", "Sent", "Received" ],
    rows:  [ 
        {
            name: "192.168.5.252",
            i18n: false,
            data: [ 720, 400, 320 ]
        },
        {
            name: "192.168.5.211",
            i18n: false,
            data: [ 550, 300, 250 ]
        }
    ]
}
```

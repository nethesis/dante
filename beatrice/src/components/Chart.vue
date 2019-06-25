<template>
  <div :id="chartId">
    <apexchart
      :ref="chartId"
      :type="type"
      :width="width"
      :height="height-30"
      :options="options"
      :series="series"
    />
  </div>
</template>

<script>
import Filters from "../filters";

export default {
  name: "Chart",
  props: {
    chartId: String,
    type: String,
    series: Array,
    categories: Array,
    width: Number,
    height: Number,
    theme: Boolean,
    palette: String,
    sparkline: Boolean,
    title: String,
    labels: Array,
    unit: String
  },
  watch: {
    theme: function(theme) {
      this.options = this.initOptions();
    },
    palette: function() {
      this.options = this.initOptions();
    },
    title: function() {
      this.options = this.initOptions();
    },
    sparkline: function() {
      this.options = this.initOptions();
    }
  },
  data() {
    return {
      options: this.initOptions()
    };
  },
  methods: {
    initOptions() {
      var context = this;
      return {
        chart: {
          toolbar: {
            show: true,
            tools: this.sparkline
              ? false
              : {
                  download: true,
                  selection: false,
                  zoom: false,
                  zoomin: false,
                  zoomout: false,
                  pan: false,
                  reset: false
                }
          },
          id: this.id,
          background: this.theme ? "#fffff" : "#1d1e1e",
          sparkline: {
            enabled: this.sparkline
          }
        },
        dataLabels: {
          enabled: false
        },
        labels: this.labels ? this.labels : [],
        title: {
          text: this.sparkline ? "" : this.title,
          floating: false,
          align: "left",
          style: {
            fontSize: "14px",
            color: this.theme ? "black" : "white"
          }
        },
        markers: {
          size: this.sparkline ? 0 : 4
        },
        xaxis: {
          categories: this.categories
            ? this.categories.map(function(c) {
                return Filters.formatter(c, "date");
              })
            : []
        },
        yaxis: {
          labels: {
            formatter: function(value, timestamp, index) {
              return Filters.formatter(value, context.unit);
            }
          }
        },
        legend: {
          position: "top",
          onItemClick: {
            toggleDataSeries: false
          }
        },
        theme: {
          mode: this.theme ? "light" : "dark",
          palette: this.palette || "palette1"
        },
        tooltip: {
          fillSeriesColor: true
        }
      };
    }
  }
};
</script>

<style>
.apexcharts-menu.open {
  width: 130px !important;
}

.dark .apexcharts-menu > .apexcharts-menu-item:hover {
  background: #88898a;
}

.apexcharts-tooltip {
  box-shadow: 0px 0px 3px 0px #e0e1e2 !important;
}
</style>

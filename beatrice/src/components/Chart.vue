<template>
  <div :id="chartId">
    <apexchart
      :ref="chartId"
      :type="type"
      :width="width"
      :height="height"
      :options="options"
      :series="series"
    />
  </div>
</template>

<script>
export default {
  name: "Chart",
  props: {
    chartId: String,
    type: String,
    series: Array,
    width: Number,
    height: Number,
    theme: Boolean,
    palette: String,
    sparkline: Boolean,
    title: String
  },
  watch: {
    theme: function() {
      this.options = this.initOptions();
    },
    palette: function() {
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
          title: {
            text: "Month"
          }
        },
        yaxis: {
          title: {
            text: "Temperature"
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
        }
      };
    }
  }
};
</script>

<style scoped>
</style>

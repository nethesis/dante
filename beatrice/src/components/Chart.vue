<template>
  <div :id="chartId">
    <apexchart :type="type" :width="width" :height="height" :options="options" :series="series"/>
  </div>
</template>

<script>
export default {
  name: "Chart",
  props: {
    chartId: String,
    type: String,
    width: Number,
    height: Number,
    series: Array,
    theme: Boolean
  },
  watch: {
    theme: function() {
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
            tools: {
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
          background: this.theme ? "#fffff" : "#1c1d1d"
        },
        //colors: ["#77B6EA", "#545454"],
        dataLabels: {
          enabled: false
        },
        title: {
          text: "Average High & Low Temperature",
          floating: false,
          align: "left",
          style: {
            fontSize: "14px",
            color: this.theme ? "black" : "white"
          }
        },
        markers: {
          size: 4
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
          horizontalAlign: "right",
          onItemClick: {
            toggleDataSeries: false
          }
        },
        theme: {
          mode: this.theme ? "light" : "dark",
          palette: "palette1"
        }
      };
    }
  }
};
</script>

<style scoped>
</style>

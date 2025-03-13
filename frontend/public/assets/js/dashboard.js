    /* sessions overview */
    var options = {
        series: [{
            name: "This Year",
            type: "area",
            data: [34, 22, 65, 35, 86, 58, 70, 34, 56, 78, 35, 55],
        }, {
             name: "This Year",
            type: "line",
            data: [12, 40, 35, 55, 25, 45, 34, 55, 37, 56, 21, 35],
        }, {
            name: 'Total Revenue',
            type: "bar",
            data: [12, 20, 30, 20, 40, 45, 35, 30, 10, 10, 25, 35]
        }],
        chart: {
            height: 338,
            type: 'line',
            toolbar: {
                show: false
            },
            stacked: "true",
            dropShadow: {
                enabled: true,
                enabledOnSeries: undefined,
                top: 7,
                left: 0,
                blur: 1,
                color: ["var(--primary-color)", "rgba(0, 0, 0, 0.5)" , "rgba(215, 124, 247,0.6)"],
                opacity: 0.1,
            },
        },
        grid: {
            borderColor: '#f1f1f1',
            strokeDashArray: 3
        },
        plotOptions: {
            bar: {
              horizontal: false,
              columnWidth: "18%",
              endingShape: "rounded",
              borderRadius: 3,
            },
          },
        dataLabels: {
            enabled: false
        },
        stroke: {
            width: [1.5, 4, 0],
            curve: ["smooth","straight","smooth"],
            dashArray: [0, 3, 0],
        },
        legend: {
            show: true,
            position: 'bottom',
            horizontalAlign: 'center',
            height: 40,
            offsetX: 0,
            offsetY: 10,
            markers: {
              strokeWidth: 0,
              size: 7,
          },
        },
        yaxis: {
            title: {
                style: {
                    color: '#adb5be',
                    fontSize: '14px',
                    fontFamily: 'poppins, sans-serif',
                    fontWeight: 500,
                    cssClass: 'apexcharts-yaxis-label',
                },
            },
        },
        xaxis: {
            type: "month",
            categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
            axisBorder: {
                show: true,
                color: 'rgba(119, 119, 142, 0.05)',
                offsetX: 0,
                offsetY: 0,
            },
            axisTicks: {
                show: true,
                borderType: 'solid',
                color: 'rgba(119, 119, 142, 0.05)',
                offsetX: 0,
                offsetY: 0
            },
            labels: {
                rotate: -90,
                style: {
                    colors: "#8c9097",
                    fontSize: '11px',
                    fontWeight: 500,
                    cssClass: 'apexcharts-xaxis-label',
                },
            }
        },
        fill: {
            type: 'gradient',
            gradient: {
                shadeIntensity: 1,
                opacityFrom: 0.4,
                opacityTo: 0.1,
                stops: [0, 90, 100],
                colorStops: [
                    [
                        {
                            offset: 0,
                            color: "var(--primary08)",
                            opacity: 1
                        },
                        {
                            offset: 75,
                            color: "var(--primary08)",
                            opacity: 1
                        },
                        {
                            offset: 100,
                            color: "var(--primary08)",
                            opacity: 1
                        }
                    ],
                    [
                        {
                            offset: 0,
                            color: "rgba(179, 254, 74, 1)",
                            opacity: 1
                        },
                        {
                            offset: 75,
                            color: "rgba(179, 254, 74,1)",
                            opacity: 1
                        },
                        {
                            offset: 100,
                            color: "rgba(179, 254, 74,1)",
                            opacity: 1
                        }
                    ],
                    [
                        {
                            offset: 0,
                            color: "rgba(230, 231, 233, 0.4)",
                            opacity: 1
                        },
                        {
                            offset: 75,
                            color: "rgba(230, 231, 233, 0.4)",
                            opacity: 1
                        },
                        {
                            offset: 100,
                            color: "rgba(230, 231, 233, 0.4)",
                            opacity: 1
                        }
                    ],
                ]
            }
        },
        colors: ["var(--primary09)","rgb(179, 254, 74)", "rgba(230, 231, 233, 1)"],
        tooltip: {
            shared: true,
            intersect: false,
            theme: "dark",
        },
    }
    var chart = new ApexCharts(document.querySelector("#sales-overview"), options);
    chart.render();
    /* sessions overview */
   /* sessions by device */
   var options01 = {
    series: [1754, 1234, 878],
    labels: ["Sales", "Income", "Profit"],
    chart: {
        height: 287,
        type: 'donut',
    },
    dataLabels: {
        enabled: false,
    },

    legend: {
        show: false,
    },
    stroke: {
        show: true,
        curve: 'smooth',
        lineCap: 'round',
        colors: "#fff",
        width: 0,
        dashArray: 0,
    },
    plotOptions: {
        pie: {
            expandOnClick: false,
            donut: {
                size: '80%',
                background: 'transparent',
                labels: {
                    show: false,
                    name: {
                        show: true,
                        fontSize: '20px',
                        color: '#495057',
                        offsetY: -4
                    },
                    value: {
                        show: false,
                        fontSize: '18px',
                        color: undefined,
                        offsetY: 8,
                        formatter: function (val) {
                            return val + "%"
                        }
                    },
                    total: {
                        show: false,
                        showAlways: true,
                        label: 'Total',
                        fontSize: '22px',
                        fontWeight: 600,
                        color: '#495057',
                    }

                }
            }
        }
    },
    colors: ["var(--primary-color)", "rgba(179, 254, 74, 1)", "var(--muted2)"],
};
var chart01 = new ApexCharts(document.querySelector("#sales-donut"), options01);
chart01.render();
/* sessions by device */

var swiper3 = new Swiper(".related", {
    loop: true,
    slidesPerView: 1,
    spaceBetween: 30,
    autoplay: {
      delay: 1500,
      disableOnInteraction: false,
    },
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
    breakpoints: {
      768: {
        slidesPerView: 2,
      },
      1024: {
        slidesPerView: 2,
      },
      1400: {
        slidesPerView: 3,
      },
    },
  });
  (function () {
    "use strict";
    // Select all elements with the specified class
    var elements = document.querySelectorAll(".rating-stars2");
  
    // Loop through each element and create an instance of raterJs
    elements.forEach(function (element) {
      var starRating = raterJs({
        starSize: 24,
        max: 5,
        rating: 4,
        element: element,
        disableText: "Custom disable text!",
        ratingText: "My custom rating text {rating}",
        showToolTip: true,
        rateCallback: function ratingCallback(rating, done) {
          starRating.setRating(rating);
          starRating.disable();
          done();
        },
      });
    });
    /* Custom messages */
  })();

  
 /* Earnings */
 var options = {
    series: [
      {
        name: "Earnings",
        data: [
          {
            x: "25",
            y: 21
          },
          {
            x: "30",
            y: 12
          },
          {
            x: "35",
            y: 22
          },
          {
            x: "40",
            y: 12
          },
          {
            x: "45",
            y: 32
          }
        ]
      },
      {
        name: "Profit",
        data: [
          {
            x: "25",
            y: 11
          },
          {
            x: "30",
            y: 19
          },
          {
            x: "35",
            y: 12
          },
          {
            x: "40",
            y: 18
          },
          {
            x: "45",
            y: 2
          }
        ]
      }
    ],
    chart: {
      height: 263,
      type: "line",
      dropShadow: {
        enabled: true,
        color: "#000",
        top: 18,
        left: 7,
        blur: 10,
        opacity: 0.06
      },
      toolbar: {
        show: false
      }
    },
    stroke: {
      width: 6,
      curve: "smooth",
  
    },
    colors: ["var(--primary-color)", "rgb(179, 254, 74)"],
    dataLabels: {
      enabled: false,
    },
    grid: {
      borderColor: "#edeef1",
      strokeDashArray: 2,
    },
    legend: {
      show: true,
      position: "top",
      horizontalAlign: 'right', 
      markers: {
        width: 0,        // Marker width
        height: 0,       // Marker height
        strokeWidth: 0,
        strokeColor: "#fff",
        fillColors: undefined,
        radius: 10,       // Radius for rounded markers
        customHTML: undefined,
        onClick: undefined,
        offsetX: -2,
        offsetY: -2,
    },
    },
    xaxis: {
      axisBorder: {
        show: false,
        offsetX: 0,
        offsetY: 0,
      },
    },
    yaxis: {
      axisBorder: {
        show: false,
        offsetX: 0,
        offsetY: 0,
      },
    }
  };
  
  var chart = new ApexCharts(document.querySelector("#earnings"), options);
  chart.render();
  /* Earnings */
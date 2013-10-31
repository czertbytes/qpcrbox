'use strict';

qpcrApp.service('contentParser', function() {
    this.parse = function(type, qpcrData) {
        var detectors = {};

        if (type == 'ab7300') {
            //  AB7300 SDS 1.4
            if (qpcrData.indexOf('Applied Biosystems 7300 Real-Time PCR System') > 0 &&
                qpcrData.indexOf('SDS v1.4') > 0) {

                var lines = qpcrData.split('\n'), section = 1, lineValues = [];
                angular.forEach(lines, function(line) {
                    if (line.length == 0) {
                        section = section + 1;
                    } else {
                        if (section == 11) {
                            lineValues = line.split(',');

                            //  suggest mock name
                            detectors[lineValues[3]] = (lineValues[3].toLowerCase() == 'mock');
                        }
                    }
                });
            }
        }

        return detectors;
    };
});


qpcrApp.service('chartGenerator', function() {
    return {
        create: function(elm, chartData) {
            var margin = {top: 20, right: 20, bottom: 30, left: 40},
                width = 900 - margin.left - margin.right,
                height = 300 - margin.top - margin.bottom;

            var svg = d3.select(elm).append("svg:svg")
                .attr("width", width + margin.left + margin.right)
                .attr("height", height + margin.top + margin.bottom)
                .on('mousemove', function(d) {
                    var coord = d3.mouse(elm);
                    infobox.style("left", (coord[0] + elm.offsetLeft + 10) + "px");
                    infobox.style("top", (coord[1] + elm.offsetTop - 55) + "px");
                })
                .append("svg:g")
                .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

            var x = d3.scale.ordinal().domain(chartData.map(function(d) { return d.name; })).rangeRoundBands([0, width], .1);
            var y = d3.scale.linear().domain([0, d3.max(chartData, function(d) { return d.value + d.error })]).range([height, 0]);

            var xAxis = d3.svg.axis().scale(x).orient("bottom");
            svg.append("svg:g").attr("class", "x axis").attr("transform", "translate(0," + height + ")").call(xAxis);

            var yAxis = d3.svg.axis().scale(y).orient("left");
            svg.append("svg:g").attr("class", "y axis").call(yAxis).append("svg:text").attr("transform", "rotate(-90)").attr("y", 6).attr("dy", ".71em").style("text-anchor", "end").text("RQ");

            var bar = svg.selectAll("bar").data(chartData);
            bar.enter().append("svg:rect").attr("class", "bar")
                .on("mouseover", function(d) { barMouseOver(d); })
                .on("mouseout", function(d) { barMouseOut(); })
                .attr("x", function(d) { return x(d.name); })
                .attr("y", function(d) { return y(d.value); })
                .attr("rx", 2)
                .attr("ry", 2)
                .attr("width", x.rangeBand())
                .attr("height", function(d) { return height - y(d.value); });

            var errorLine = svg.selectAll("err").data(chartData);
            //  center dashed-line
            errorLine.enter().append("line").attr("class", "err-center")
                .on("mouseover", function(d) { barMouseOver(d); })
                .on("mouseout", function(d) { barMouseOut(); })
                .attr("x1", function(d) { return x(d.name) + 34})
                .attr("x2", function(d) { return x(d.name) + 34})
                .attr("y1", function(d) { return y(d.value + (d.error / 2)); })
                .attr("y2", function(d) { return y(d.value - (d.error / 2)); });

            //  top whisker
            errorLine.enter().append("line").attr("class", "err-whisker")
                .on("mouseover", function(d) { barMouseOver(d); })
                .on("mouseout", function(d) { barMouseOut(); })
                .attr("x1", function(d) { return x(d.name) + 28})
                .attr("x2", function(d) { return x(d.name) + 40})
                .attr("y1", function(d) { return y(d.value + (d.error / 2)); })
                .attr("y2", function(d) { return y(d.value + (d.error / 2)); });

            //  bottom whisker
            errorLine.enter().append("line").attr("class", "err-whisker")
                .on("mouseover", function(d) { barMouseOver(d); })
                .on("mouseout", function(d) { barMouseOut(); })
                .attr("x1", function(d) { return x(d.name) + 28})
                .attr("x2", function(d) { return x(d.name) + 40})
                .attr("y1", function(d) { return y(d.value - (d.error / 2)); })
                .attr("y2", function(d) { return y(d.value - (d.error / 2)); });
        }
    }
});


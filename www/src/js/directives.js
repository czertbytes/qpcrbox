'use strict';

qpcrApp.directive('fileUpload', function($parse) {
    return function(scope, element, attrs) {
        element.bind('drop', function(e) {
            e.preventDefault();
            e.stopPropagation();
            var file = e.originalEvent.dataTransfer.files[0], reader = new FileReader;
            reader.onload = function(event) {
                scope.$apply(function() {
                    var qpcrData = event.target.result.replace(/\r\n/g, '\n');
                    $parse(attrs.ngModel).assign(scope, qpcrData);
                });
            };

            reader.readAsText(file);
            return false;
        });
    };
});

qpcrApp.directive('withChart', function(chartGenerator) {
    return function(scope, element, attrs) {
        var chartData = [], detectorData = scope.expData.Detectors[scope.detectorName];

        Object.keys(detectorData).map(function(tg) {
            if (tg != '$$hashKey') {
                chartData.push({name: tg, value: detectorData[tg].RQ, error: detectorData[tg].RQErr});
            }
        });

        chartGenerator.create(element[0].parentNode, chartData);
    };
});

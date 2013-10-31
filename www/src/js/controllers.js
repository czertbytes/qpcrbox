'use strict';

qpcrApp.controller('MainCtrl', function($scope, $location) {
    $scope.setRoute = function(route) {
        $location.path(route);
    };
});

qpcrApp.controller('AB7300Ctrl', function($scope, $http, contentParser, $compile) {
    $scope.qpcrData = '';
    $scope.mock = '';
    $scope.genes = [];
    $scope.showMockGeneSelection = false;
    $scope.expData = '';
    $scope.detectors = [];
    $scope.rateLimitExceed = false;
    $scope.rateLimitReset = '';

    $scope.$watch('qpcrData', function() {
        $scope.update();
    }, true);

    $scope.$watch('expData', function() {
        if ($scope.expData) {
            $scope.detectors = Object.keys($scope.expData.Detectors);
        }
    }, true);

    $scope.update = function() {
        $scope.updateGenesSelection();
        $scope.compute();
    };

    $scope.compute = function() {
        if ($scope.qpcrData && $scope.mock) {

            //  post experiment data
            $http.post('http://api.qpcrbox.com/v1/qpcr/ab7300?mock=' + $scope.mock, $scope.qpcrData)
                .success(function(data, status, headers, config) {
                    console.log("response code: " + status);
                    if (status == 201) {
                        //  Get experiment result data
                        $http.get('http://api.qpcrbox.com/v1/experiment/' + data.ExperimentId, {headers: {'Accept': 'application/json'}})
                            .success(function(data, status, headers, config) {
                                if (status == 200) {
                                    $scope.expData = data;
                                }
                            })
                            .error(function(data, status, headers, config) {
                                if (status == 429) {
                                    $scope.rateLimitExceed = true;
                                    $scope.rateLimitReset = data.RetryAfter;
                                    console.log('limit exceeded, try later!');
                                }

                                console.log("get exp by id failed: " + status);
                            });
                    }
                })
                .error(function(data, status, headers, config) {
                    if (status == 429) {
                        $scope.rateLimitExceed = true;
                        $scope.rateLimitReset = data.RetryAfter;
                        console.log('limit exceeded, try later!');
                    }

                    console.log("post exp data failed: " + status);
                });
        }
    };

    $scope.updateGenesSelection = function() {
        if ($scope.qpcrData) {
            console.log("updating mock gene selection!");
            $scope.genes = [];
            $scope.mock = '';
            $scope.expData = '';

            var qpcrDataGenes = contentParser.parse('ab7300', $scope.qpcrData);
            var qpcrDataGeneNames = Object.keys(qpcrDataGenes);
            angular.forEach(qpcrDataGeneNames, function(geneName) {
                $scope.genes.push({id: geneName, name: geneName});

                //  set suggested mock name
                if (qpcrDataGenes[geneName] == true) {
                    $scope.mock = geneName;
                }
            });

            $scope.showMockGeneSelection = true;
        }
    };
});

'use strict';

var qpcrApp = angular.module('qpcrApp', [], function($routeProvider) {
    $routeProvider.when('/qpcr/ab7300', { templateUrl: '/partials/qpcr/ab7300.html', controller: qpcrApp.AB7300Ctrl });
    $routeProvider.otherwise({ redirectTo: '/qpcr/ab7300' });

    /*$locationProvider.html5Mode(true).hashPrefix('!');

    $routeProvider.when('/qpcr/ab7300', { templateUrl: '/partials/qpcr/ab7300.html', controller: qpcrApp.AB7300Ctrl });
    $routeProvider.otherwise({ redirectTo: '/qpcr/ab7300' });*/
});

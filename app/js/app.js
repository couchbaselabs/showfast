/*global angular, Timeline, RunList: true*/


angular.module('sf', ['ngRoute', 'nvd3ChartDirectives'])
	.config([
	    '$routeProvider', function($routeProvider) {
	        $routeProvider.
		        when('/timeline', {templateUrl: 'partials/timeline.html', controller: Timeline}).
		        when('/runs/:metric/:build', {templateUrl: 'partials/runs.html', controller: RunList}).
		        otherwise({redirectTo: '/timeline'});
        }
    ]);

angular.module('sf-admin', [])
    .controller("AdminList", AdminList);

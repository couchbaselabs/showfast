/*global angular, MetricList, RunList: true*/


angular.module('showfast', ['ngRoute', 'nvd3ChartDirectives']).
	config(['$routeProvider', function($routeProvider) {
	$routeProvider.
		when('/timeline', {templateUrl: 'partials/timeline.html', controller: MetricList}).
		when('/runs/:metric/:build', {templateUrl: 'partials/runs.html', controller: RunList}).
		otherwise({redirectTo: '/timeline'});
}]);

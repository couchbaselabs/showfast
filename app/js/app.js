/*global angular, MetricList, RunList: true*/


angular.module('sf', ['ngRoute', 'nvd3ChartDirectives'])
	.config([
	    '$routeProvider', function($routeProvider) {
	        $routeProvider.
		        when('/timeline', {templateUrl: 'partials/timeline.html', controller: MetricList}).
		        when('/runs/:metric/:build', {templateUrl: 'partials/runs.html', controller: RunList}).
		        otherwise({redirectTo: '/timeline'});
        }
    ]);

angular.module('sf-admin', [])
    .controller("AdminList", AdminList);

angular.module('sf-feed', [])
    .controller("Feed", Feed);

angular.module('sf-release', [])
    .controller("ReleaseList", ReleaseList);

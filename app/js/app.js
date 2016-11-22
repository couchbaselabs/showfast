angular.module('sf', ['ngRoute', 'nvd3ChartDirectives'])
	.config([
	    '$routeProvider', function($routeProvider) {
	        $routeProvider.
		        when('/timeline/:os/:component/:category', {templateUrl: '/static/timeline.html', controller: MenuRouter}).
		        when('/runs/:metric/:build', {templateUrl: '/static/runs.html', controller: RunList}).
		        otherwise({redirectTo: 'timeline/Linux/kv/max_ops'});
        }
    ])
	.controller('MainDashboard', MainDashboard);

angular.module('sf-admin', ['tableSort'])
    .controller("AdminList", AdminList);

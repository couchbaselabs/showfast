/*jshint jquery: true, browser: true*/
/*global jlinq, d3: true*/


function MetricList($scope, $http) {
	'use strict';

	$http.get('/all_metrics').success(function(data) {
		if (data.length) {
			$scope.metrics = data;
		} else {
			$scope.error = true;
			return;
		}

		$http.get('/all_clusters').success(function(data) {
			for (var i = 0, l = $scope.metrics.length; i < l; i++ ) {
				$scope.metrics[i].cluster = jlinq.from(data)
					.starts('Name', $scope.metrics[i].cluster)
					.ends('Name', $scope.metrics[i].cluster)
					.select()[0];
			}
		});

		var format = d3.format(',');
		$scope.valueFormatFunction = function(){
			return function(d) {
				return format(d);
			};
		};

		$http.get('/all_timelines').success(function(data) {
			for (var i = 0, l = $scope.metrics.length; i < l; i++ ) {
				var id = $scope.metrics[i].id;
				$scope.metrics[i].chartData = [{"values": data[id]}];
			}
		});

		$scope.setSelectedCategory = function (value) {
			if ($scope.selectedCategory === value) {
				$scope.selectedCategory = undefined;
			} else {
				$scope.selectedCategory = value;
			}
		};

		$scope.selectedCategory = "all";

		$scope.categories = [{
			"id": "all", "title": "All"
		}, {
			"id": "comp", "title": "Compaction"
		}, {
			"id": "index", "title": "Indexing"
		}, {
			"id": "reb", "title": "Rebalance"
		}, {
			"id": "kv", "title": "KV"
		}, {
			"id": "query", "title": "Query"
		}, {
			"id": "xdcr", "title": "XDCR"
		}];

		$scope.byCategory = function(entry){
			var category = $scope.selectedCategory;

			if (category === "all") {
				return true;
			} else if (entry.id.substring(0, category.length) === category) {
				return true;
			} else {
				return false;
			}
		};
	});
}

function RunList($scope, $routeParams, $http) {
	$http({method: 'GET', url: '/all_runs', params: $routeParams})
	.success(function(data) {
		$scope.runs = data;
	});
}

function AdminList($scope, $http) {
	$scope.deleteBenchmark = function(entry) {
		$http({method: 'POST', url: '/delete', data: {id: entry}})
		.success(function(data) {
			location.reload();	
		});
	};

	$http.get('/all_benchmarks').success(function(data) {
		$scope.benchmarks = data;
	});
}

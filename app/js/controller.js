function MainDashboard($scope, $http, $routeParams) {
	'use strict';

	DefineMenu($scope, $http);

	GetData($scope, $http);

	var format = d3.format(',');
	$scope.valueFormatFunction = function() {
		return function(d) {
			return format(d);
		};
	};

	$scope.$on('elementClick.directive', function(event, data) {
		var build = data.data[0],
			metric = event.targetScope.id,
			a = $("#run_"  + metric);
		a.attr("href", "/#/runs/" + metric + "/" + build);
		a[0].click();
	});
}

function GetData($scope, $http) {
	$http.get('/api/v1/metrics').success(function(metrics) {
		$http.get('/api/v1/clusters').success(function(data) {
			for (var i = 0, l = metrics.length; i < l; i++ ) {
				metrics[i].cluster = jlinq.from(data)
					.starts('name', metrics[i].cluster)
					.ends('name', metrics[i].cluster)
					.select()[0];
			}
		});

		$http.get('/api/v1/timelines').success(function(data) {
			$scope.metrics = [];
			var j = 0;
			for (var i = 0, l = metrics.length; i < l; i++ ) {
				var id = metrics[i].id;
				if (id in data) {
					$scope.metrics[j] = metrics[i];
					$scope.metrics[j].chartData = [{"key": id, "values": data[id]}];
					$scope.metrics[j].link = id.replace(".", "_");
					j++;
				}
			}
		});
	});
}

function MenuRouter($scope, $routeParams, $location) {
	$scope.activeOS = $routeParams.os;
	$scope.activeComponent = $routeParams.component;
	$scope.activeCategory = $routeParams.category;

	$scope.setActiveOS = function(os) {
		$location.path("/timeline/" + os + "/" + $scope.activeComponent + "/" + $scope.activeCategory);
	};

	$scope.setActiveComponent = function(component) {
		$location.path("/timeline/" + $scope.activeOS + "/" + component + "/" + $scope.components[component].categories[0].id);
	};

	$scope.setActiveCategory = function(category) {
		$location.path("/timeline/" + $scope.activeOS + "/" + $scope.activeComponent + "/" + category);
	};

	DefineFilters($scope);
}

function DefineMenu($scope, $http) {
	$http.get('/static/menu.json').success(function(menu) {
		$scope.components = menu.components;
		$scope.oses =  menu.oses;
	});
}

function DefineFilters($scope) {
	$scope.byOS = function(metric) {
		switch($scope.activeOS) {
			case "Linux":
				return metric.cluster.os.indexOf("Windows") === -1;
			case "Windows":
				return metric.cluster.os.indexOf("Windows") === 0;
		}
	};

	$scope.byComponentAndCategory = function(metric) {
		switch($scope.activeComponent) {
			case "kv":
			case "reb":
			case "index":
			case "query":
			case "n1ql":
			case "secondary":
			case "xdcr":
			case "ycsb":
			case "fts":
				if (metric.component === $scope.activeComponent) {
					return metric.category === $scope.activeCategory;
				}
				break;
			case "tools":
				if (metric.component === $scope.activeComponent) {
					return metric.id.indexOf($scope.activeCategory) !== -1;
				}
				break;
			default:
				return false;
		}
	};
}

function RunList($scope, $routeParams, $http) {
	$http({method: 'GET', url: '/api/v1/runs/' + $routeParams.metric + '/' + $routeParams.build})
	.success(function(data) {
		$scope.runs = data;
	});
}

function AdminList($scope, $http) {
	$scope.deleteBenchmark = function(benchmark) {
		$http({method: 'DELETE', url: '/api/v1/benchmarks/' + benchmark.id})
		.success(function(data) {
			var index = $scope.benchmarks.indexOf(benchmark);
			$scope.benchmarks.splice(index, 1);
		});
	};

	$scope.reverseHidden = function(id) {
		$http({method: 'PATCH', url: '/api/v1/benchmarks/' + id});
	};

	$http.get('/static/menu.json').success(function(menu) {
		$scope.components = {};
		for (component in menu.components) {
			$scope.components[component] = menu.components[component].title;
		}

		$scope.activeComponent = Object.keys($scope.components)[0];
	});

	$http.get('/api/v1/benchmarks').success(function(data) {
		$scope.benchmarks = data;

		$scope.setActiveComponent = function(component) {
			$scope.activeComponent = component;
		};

		$scope.byComponent = function(benchmark) {
			return benchmark.metric.indexOf($scope.activeComponent) === 0;
		};
	});
}

angular
	.module('showfast', ['tableSort'])
    .controller("AdminList", AdminList);

function GetBenchmarks($scope, $http) {
	$http.get('/api/v1/benchmarks/' + $scope.activeComponent + '/' + $scope.activeCategory).success(function(data) {
		$scope.benchmarks = data;
	});
}

function InitMenu($scope, $http) {
	$http.get('/static/menu.json').success(function(menu) {
		$scope.components = menu.components;
		$scope.activeComponent = Object.keys($scope.components)[0];
		$scope.activeCategory = $scope.components[$scope.activeComponent].categories[0].id;

		GetBenchmarks($scope, $http);
	});
}

function AdminList($scope, $http) {
	InitMenu($scope, $http);

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

	$scope.setActiveComponent = function(component) {
		$scope.activeComponent = component;
		$scope.activeCategory = $scope.components[component].categories[0].id;

		GetBenchmarks($scope, $http);
	};

	$scope.setActiveCategory = function(category) {
		$scope.activeCategory = category;

		GetBenchmarks($scope, $http)
	};
}

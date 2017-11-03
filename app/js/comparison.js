angular
    .module('showfast', [])
    .controller('ComparisonCtrl', ComparisonCtrl);

function ComparisonCtrl($scope, $http) {
    $( "#view" ).hide();

    $http.get('/static/menu.json').success(function(menu) {
        $scope.components = menu.components;
    });

    $http.get('/api/v1/builds').success(function(builds) {
        $scope.builds = builds;
        $scope.lhb = '4.6.3-4136';
        $scope.rhb = '5.0.0-3519';
    });

    $scope.$watch('lhb', function() {
        Compare($scope, $http);
    });

    $scope.$watch('rhb', function() {
        Compare($scope, $http);
    });

    $scope.calcWidthClass = calcWidthClass;

    $scope.calcOffsetClass = calcOffsetClass;

    $scope.fmtDiff = fmtDiff;
}

function fmtDiff(results) {
    var diff = 100 * (results[1].value / results[0].value - 1);
    if (diff > 0) {
        return "+" + diff.toFixed(1);
    }
    return diff.toFixed(1);
}

function calcWidthClass(results) {
    var diff = (results[1].value - results[0].value) / results[0].value;

    if ( diff > 0.5 ) {
        return "col-md-4 col-md-offset-4 bar-pos";
    } else if ( diff > 0.25 ) {
        return "col-md-3 col-md-offset-4 bar-pos";
    }    else if ( diff > 0.1 ) {
        return "col-md-2 col-md-offset-4 bar-pos";
    } else if ( diff > 0 ) {
        return "col-md-1 col-md-offset-4 bar-pos";
    } else if ( diff < -0.5 ) {
        return "col-md-4 col-md-offset-0 bar-neg";
    } else if ( diff < -0.25 ) {
        return "col-md-3 col-md-offset-1 bar-neg";
    } else if ( diff < -0.1 ) {
        return "col-md-2 col-md-offset-2 bar-neg";
    } else {
        return "col-md-1 col-md-offset-3 bar-neg";
    }
}

function calcOffsetClass(results) {
    var diff = (results[1].value - results[0].value) / results[0].value;

    if ( diff > 0.5 ) {
        return "col-md-1 col-md-offset-0";
    } else if ( diff > 0.25 ) {
        return "col-md-1 col-md-offset-1";
    } else if ( diff > 0.1 ) {
        return "col-md-1 col-md-offset-2";
    } else if ( diff > 0 ) {
        return "col-md-1 col-md-offset-3";
    } else {
        return "col-md-1 col-md-offset-4";
    }
}

function Compare($scope, $http) {
    if ( $scope.lhb === undefined || $scope.rhb === undefined ) {
        return;
    }

    $( "#view" ).hide();
    $http.get('/api/v1/comparison/' + $scope.lhb + '/' + $scope.rhb).success(function(comparisons) {
        $scope.comparisons = comparisons;
        $( "#view" ).show();
    });
}

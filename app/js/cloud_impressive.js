angular
    .module('showfast', [])
    .controller("ImpressiveCtrl", ImpressiveCtrl);


function GetImpressiveTests($scope, $http) {
    if ( $scope.lhb === undefined || $scope.rhb === undefined ) {
        return;
    }

    var url = '/api/v1/impressive/' + $scope.activeComponent + '/' + $scope.lhb + '/' +  $scope.rhb  + '/' + $scope.activeMode;
    $http.get(url).success(function(data) {
        $scope.tests = data;
        var summaryMissing = 0;
        var summaryTotal = data.length;
        data.forEach(function (arrayItem) {
            if (arrayItem.b2 === '') {
                summaryMissing += 1
            }
        });

        $scope.summaryTotal = summaryTotal;
        $scope.summaryMissing = summaryMissing;
        $scope.summaryCompleted = summaryTotal - summaryMissing;

    });

}

function InitMenu($scope, $http) {
    $http.get('/static/cloud_menu.json').success(function(menu) {
        $scope.components = menu.components;
        $scope.activeComponent = Object.keys($scope.components)[0];
        $scope.activeCategory = $scope.components[$scope.activeComponent].categories[0].id;
        GetImpressiveTests($scope, $http);
    });
}

function ImpressiveCtrl($scope, $http) {
    InitMenu($scope, $http);

    $scope.activeMode = 'Active';

    $http.get('/api/v1/builds').success(function(builds) {
        $scope.builds = builds;
        $scope.lhb = '6.0.0-1693';
        $scope.rhb = '6.5.0-1441';
    });

    $scope.$watch('lhb', function() {
        GetImpressiveTests($scope, $http);
    });

    $scope.$watch('rhb', function() {
        GetImpressiveTests($scope, $http);
    });

    $scope.setActiveComponent = function(component) {
        $scope.activeComponent = component;
        GetImpressiveTests($scope, $http);

    };


    $scope.setActiveMode = function(mode) {
        $scope.activeMode = mode;
        GetImpressiveTests($scope, $http);
    };


    $scope.calcWidthClass = calcWidthClass;

    $scope.fmtDiff = fmtDiff;

    $scope.GetRowStyle = GetRowStyle;

}


function GetRowStyle(isactive, iswindows){
    if (isactive) {
        if (iswindows) {
            return "t-windowstest";
        }
        return "t-activetest";
    }
    return "t-nonactivetest";
}

function fmtDiff(v1, v2){
    var diff = Math.abs(100 - ((100 * v2) / v1));
    return diff.toFixed(1);
}

function calcWidthClass(v1, v2){
    var diff = Math.abs(100 - ((100 * v2) / v1));
    if ( diff > 50 ) {
        return "col-md-10 col-md-offset-0 bar-bad";
    } else if ( diff > 25 ) {
        return "col-md-7 col-md-offset-0 bar-bad";
    }    else if ( diff > 10) {
        return "col-md-5 col-md-offset-0 bar-bad";
    } else if ( diff > 0 ) {
        return "col-md-3 col-md-offset-0 bar-good";
    }
    return "col-md-3 col-md-offset-0 bar-good"
}
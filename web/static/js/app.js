var app = angular.module("sparrowUI", ["ngRoute"]);
app.config(function($routeProvider) {  
    $routeProvider
        .when("/", {
            templateUrl: "main.html" ,
            controller: 'mainController'
        })
        .when("/login", {
            templateUrl: "login.html" ,
            controller: 'loginController'
        })
        .when("/db", {
            templateUrl: "database.html",
            controller: 'dbController'
        })
        .when("/db/upload", {
            templateUrl: "database.html",
            controller: 'dbController'
        });
});

app.factory('sparrow', function() {
    var self = {};
    self.token = '';
    self.currentDb = null;

    self.connect = function(_host) {
        self.client = new SparrowDb({
            host: _host
        });
    };

    self.getClient = function() {
        return self.client;
    };

    return self;
});

app.controller('loginController', function($scope, $location, sparrow, $rootScope) {
    $scope.loginData = { host: '127.0.0.1:8081', username: 'sparrow', password: 'sparrow' };
    $scope.error = '';

    $scope.doLogin = function() {
        console.log($scope.loginData)
        sparrow.connect($scope.loginData.host);

        sparrow.getClient().login($scope.loginData.username, $scope.loginData.password)
            .success(function(r) {
                sparrow.token = r.token;
                $rootScope.$apply(function() {
                    $location.path("/");
                });
            }).error(function(xhr) {
                $scope.$apply(function() {
                    $scope.error = 'Invalid user and/or password';
                });
            });
    };
});

app.controller('mainController', function($scope, $location, sparrow, $rootScope) {
    sparrow.getClient().showDatabases().success(function(r) {
        $scope.$apply(function() {
            $scope.databases = r.content._all;
        });
    });

    $scope.dbInfo = function(db) {
        sparrow.currentDb = db;
        $location.path("/db");
    }

    $scope.dbDrop = function(db) {
        if (confirm('Drop ' + db + ' ?') == true) {
            sparrow.getClient().dropDatabase(_currentDb)
                .success(function(r) {
                    alert('Database dropped')
                }).error(function(xhr) {
                    if (xhr.status == 404) {
                        $rootScope.$apply(function() {
                            $location.path("/");
                        });
                        return;
                    }
                    alert('Could not drop database');
                });
        }
    }
});


app.controller('dbController', function($scope, $location, sparrow, $rootScope) {
    $scope.currentDb = sparrow.currentDb;

    var updateInfo = function() {
        sparrow.getClient().infoDatabase(sparrow.currentDb)
            .success(function(r) {
                $scope.$apply(function() {
                    $scope.info = r.content;
                });
            }).error(function(xhr) {});
    }
    updateInfo();

    $scope.refresh = function() {
        updateInfo();
    }
});
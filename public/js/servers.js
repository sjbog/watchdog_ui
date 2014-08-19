angular.module('ServersApp', [])
	.controller('ServersController', ['$scope', '$http', function($scope, $http) {
		$scope.servers = {};
		if	( typeof servers != "undefined" )	{
			$scope.servers	= servers;
		}
		$scope.error	= null;
		if	( typeof error != "undefined" )	{
			$scope.error	= error;
		}

		$scope.errorHandlerFn	= function(data, status, headers, config) {
			//console.log ( data );
			if	( data	&& typeof data [ "error" ] )	{
				$scope.error	= "[ "+ status +" ] " + data [ "error" ];
			}

			setTimeout ( $scope.ServersAction, 60000 );
		};

		$scope.ServerAction = function ( serverId, action ) {
			$http({method: 'GET', url: '/api/v1/servers/' + serverId, params: { action: action }}).
				success(function(data, status, headers, config) {
					//console.log ( data );
					if	( data )	{
						if	( typeof data [ "error" ] )	{
							$scope.error	= data [ "error" ];
						}	else if	( typeof data [ "Label" ] != "undefined" ) {
							$scope.servers [data.Label] = data;
						}
					}
					setTimeout ( $scope.ServersAction, 30000 );
				} ).error ( $scope.errorHandlerFn );
		};
		$scope.ServersAction = function ( action ) {
			$rootScope	= $scope;
			$http({method: 'GET', url: '/api/v1/servers', params: { action: action }}).
				success(function(data, status, headers, config) {
					console.log ( data );
					if	( data )	{
						if	( typeof data [ "error" ] != "undefined" )	{
							$rootScope.error	= data [ "error" ];
						}	else	{
							$rootScope.servers	= data;
							$rootScope.error	= null;
						}
					}
				}).error ( $scope.errorHandlerFn );
		};


	}
	])
;
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
		$scope.updateHandle	= null;

		$scope.errorHandlerFn	= function(data, status, headers, config) {
			//console.log ( data );
			if	( data	&& typeof data [ "error" ] )	{
				$scope.error	= "[ "+ status +" ] " + data [ "error" ] + " " + data [ "error_msg" ];
			}

			$scope.updateHandle	= setTimeout ( $scope.ServersAction, 60000 );
		};

		$scope.ServerAction = function ( serverId, action ) {
			$http({method: 'GET', url: '/api/v1/servers/' + serverId, params: { action : action } }).
				success(function(data, status, headers, config) {
					//console.log ( data );
					if	( data )	{
						if	( typeof data [ "error" ] )	{
							$scope.error	= data [ "error" ];
						}	else if	( typeof data [ "label" ] != "undefined" ) {
							$scope.servers [data [ "label" ]] = data;
						}
					}
				} ).error ( $scope.errorHandlerFn );
		};

		$scope.ServersAction = function ( action ) {
			clearTimeout ( $scope.updateHandle );
			$scope.updateHandle	= null;

			$http({method: 'GET', url: '/api/v1/servers', params: { action: action }}).
				success(function(data, status, headers, config) {
					//console.log ( data );
					if	( data )	{
						if	( typeof data [ "error" ] != "undefined" )	{
							$scope.error	= data [ "error" ];
						}	else	{
							$scope.servers	= data;
							$scope.error	= null;
						}
					}

					if	( ! $scope.updateHandle ) {
						$scope.updateHandle	= setTimeout ( $scope.ServersAction, 30000 );
					}
				}).error ( $scope.errorHandlerFn );
		};


		$scope.ServersPost = function ( serverId, data_to_send ) {
			if	( data_to_send == null )	{
				data_to_send	= {};
			}

			$http({method: 'POST', url: '/api/v1/servers/' + serverId, data: data_to_send }).
				success(function(data, status, headers, config) {
					//console.log ( data );
					$scope.$serverToEditResult	= data;
					if	( ! data [ "error" ] ) {
						$scope.ServersAction ();
						$scope.$serverToEditId = data_to_send ["label"];
					}
				} )
				//TODO: error handling
				.error ( $scope.errorHandlerFn )
			;
		};

		$scope.$serverToEdit	= null;
		$scope.$serverToEditId	= null;
		$scope.$serverToEditResult	= null;
		$scope.$modal_elem	= $( '#editServerModal' ).on ( 'hidden.bs.modal', function () {
			$scope.CancelServerEdit ();
		});

		$scope.OpenServerEdit	= function ( serverId )	{
			$scope.$serverToEdit	= $scope.servers [ serverId ] || {};
			$scope.$serverToEdit	= JSON.parse ( JSON.stringify ( $scope.$serverToEdit ) );
			if	( $scope.$serverToEdit )	{
				$scope.$serverToEditId	= serverId;
				$scope.$modal_elem.modal('show');
			}
		};
		$scope.CancelServerEdit	= function () {
			$scope.$serverToEdit = null;
			$scope.$serverToEditId	= null;
		};
		$scope.SaveServerEdit	= function ()	{
			if	( $scope.$serverToEdit	&& $scope.$serverToEditId ) {

				$scope.ServersPost ( $scope.$serverToEditId, $scope.$serverToEdit );
				//$scope.servers [ $scope.$serverToEdit.label ] = $scope.$serverToEdit;
				$scope.$serverToEditResult	= { status : "Saving.. please wait" };
			}
			//$scope.$modal_elem.modal('hide');
		};

		$scope.ServersAction();
	}
	])
;
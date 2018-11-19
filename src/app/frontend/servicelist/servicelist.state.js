/**
 * Configures states for the service view
 *
 * @param {!ui.router.$stateProvider} $stateProvider
 * @ngInject
 */
import ServiceListController from "./servicelist.controller";

export default function stateConfig($stateProvider) {
    $stateProvider.state('servicelist', {
        url: '/servicelist',
        templateUrl: 'servicelist/servicelist.html',
        controller: ServiceListController,
        controllerAs: 'ctrl',
    });
}
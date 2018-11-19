import {MainController} from "./main.controller";

/**
 * Configures states for the zero state view
 *
 * @param {!ui.router.$stateProvider} $stateProvider
 * @ngInject
 */
export default function stateConfig($stateProvider) {
    $stateProvider.state('main', {
        url: '/',
        templateUrl: 'main/main.html',
        controller: MainController,
        controllerAs: 'ctrl',
    });
}
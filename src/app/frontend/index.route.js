/**
 * Global route configuration for the application.
 *
 * @param {!ui.router.$urlRouterProvider} $urlRouterProvider
 * @ngInject
 */
export default function routerConfig($urlRouterProvider) {
    // When no state is matched by an URL, redirect to default one.
    $urlRouterProvider.otherwise('/');
}
import stateConfig from "./main.state";

/**
 * Angular module for main page view
 *
 * This view is active on the first launch of the application
 */
export default angular.module(
    'k8sconsole.main',
    [
        'ui.router',
    ])
    .config(stateConfig);
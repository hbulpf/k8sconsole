import stateConfig from "./servicelist.state";

/**
 * Angular module for the service list view.
 *
 * The view shows services running in the cluster and allows to manager them
 */
export default angular.module(
    'k8sconsole.serviceList',
    [
        'ui.router',
    ])
    .config(stateConfig);
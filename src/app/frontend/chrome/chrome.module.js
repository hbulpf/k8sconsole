import chromeDirective from "./chrome.directive";

export default angular.module(
    'k8sconsole.chrome',
    [
        'ngMaterial',
        'ui.router',
    ])
    .directive('chrome', chromeDirective);
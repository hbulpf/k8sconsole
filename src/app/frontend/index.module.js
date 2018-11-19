import chromeModule from './chrome/chrome.module';
import indexConfig from './index.config';
import mainModule from './main/main.module';
import routerConfig from './index.route';
import serviceListModule from './servicelist/servicelist.module';

export default angular.module(
    'k8sconsole',
    [
        'ngAnimate',
        'ngAria',
        'ngMaterial',
        'ngMessages',
        'ngResource',
        'ngSanitize',
        'ui.router',
        mainModule.name,
        chromeModule.name,
        serviceListModule.name,
    ])
    .config(indexConfig)
    .config(routerConfig);

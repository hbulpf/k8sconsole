/**
 * @param {!md.$mdThemingProvider} $mdThemingProvider
 * @ngInject
 */
export default function config($mdThemingProvider) {
    // Create a color palette for k8sconsole
    let k8sconsoleColorPaletteName = 'k8sconsoleColorPalette';
    let k8sconsoleColorPalette = $mdThemingProvider.extendPalette('blue', {
        '500': '326de6',
    });

    // Use the k8sconsole color palette as default one.
    $mdThemingProvider.definePalette(k8sconsoleColorPaletteName, k8sconsoleColorPalette);
    $mdThemingProvider.theme("default").primaryPalette(k8sconsoleColorPaletteName);
}

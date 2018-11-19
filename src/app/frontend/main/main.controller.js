export class MainController {
    /**
     * @param {!angular.$timeout} $timeout
     * @ngInject
     */
    constructor($timeout) {
        /** @export */
        this.testValue = 9;

        $timeout(() => {
            this.testValue = 8;
        }, 4000);
    }
}

import {MainController} from 'main/main.controller';

describe('Main controller', () => {
    let vm;

    beforeEach(inject(($timeout) => {
        vm = new MainController($timeout);
    }));

    it('should do something', () => {
        expect(vm.testValue).toEqual(9);
    });
});

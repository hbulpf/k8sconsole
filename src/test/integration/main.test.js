import {MainPage} from './main.po';


describe('The main view', function () {
    let page;

    beforeEach(function () {
        browser.get('/index.html');
        page = new MainPage();
    });

    it('should do something', function() {
        expect(page.header.getText()).toContain('page');
    });
});
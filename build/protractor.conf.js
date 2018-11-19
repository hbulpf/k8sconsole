require('babel-core/register');
let conf = require('./config').default;
let path = require('path');


exports.config = {
    capabilities: {
        'browserName': 'chrome',
    },

    baseUrl: 'http://localhost:3000',

    specs: [path.join(conf.paths.integrationTest, '**/*.js')],

    framework: 'jasmine',
};

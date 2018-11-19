import fs from 'fs';
import path from 'path';
import wiredep from 'wiredep';

import config from './config';

/**
 * Returns an array of files required by Karma to run the tests.
 * @returns {!Array<string>}
 */
function getFileList() {
    // All app dependencies are required for tests. Include them.
    let wiredepOptions = {
        bowerJson: JSON.parse(fs.readFileSync(path.join(config.paths.base, 'package.json'))),
        directory: config.paths.nodeModules,
        devDependencies: false,
        customDependencies: ['angular-mocks', 'google-closure-library'],
        onError: (msg) => {
            console.log(msg);
        },
    };

    return wiredep(wiredepOptions).js
        .concat([
            path.join(config.paths.frontendTest, '**/*.json'),
            path.join(config.paths.frontendTest, '**/*.js'),
            path.join(config.paths.frontendSrc, '**/*.html'),
        ]);
}

/**
 * Exported default function which sets Karma configuration. Required by the framework.
 *
 * @param {!Object} conf
 */
export default function (conf) {
    let configuration = {
        basePath: '.',

        files: getFileList(),

        LogLevel: 'INFO',

        browserConsoleLogOptions: {terminal: true, level: ''},

        frameworks: ['jasmine-jquery', 'jasmine', 'browserify', 'closure'],

        browserNoActivityTimeout: 5 * 60 * 1000,  // 5 minutes.

        reporters: ['dots', 'coverage'],

        coverageReporter: {
            dir: config.paths.coverage,
            reporters: [
                {type: 'html', subdir: 'html'},
                {type: 'lcovonly', subdir: 'lcov'},
            ],
        },

        preprocessors: {}, // This field is filled with values later.

        // karma-browserify plugin config.
        browserify: {
            // Add source maps to outpus bundles.
            debug: true,
            // Make 'import ...' statements relative to the following paths.
            paths: [config.paths.frontendSrc, config.paths.frontendTest],
            transform: [
                // Transform ES6 code into ES5 so that browsers can digest it.
                ['babelify'],
            ],
        },

        // karma-ng-html2js-preprocessor plugin config.
        ngHtml2JsPreprocessor: {
            stripPrefix:  `${config.paths.frontendSrc}/`,
            moduleName: 'ng',
        },
    };

    // for remote linux server
    configuration.browsers = ['PhantomJS'];

    // Convert all JS code written ES6 with modules to ES5 bundles that browsers can digest.
    configuration.preprocessors[path.join(config.paths.frontendTest, '**/*.js')] =
        ['browserify', 'closure', 'closure-iit'];
    configuration.preprocessors[path.join(
        config.paths.nodeModules, 'google-closure-library/closure/goog/deps.js')] = ['closure-deps'];

    // Convert HTML templates into JS files that serve code through $templateCache.
    configuration.preprocessors[path.join(config.paths.frontendSrc, '**/*.html')] = ['ng-html2js'];

    conf.set(configuration);
}